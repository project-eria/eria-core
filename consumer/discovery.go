package eriaconsumer

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/project-eria/go-wot/protocolDnssd"
	zlog "github.com/rs/zerolog/log"
)

// browseFunc is injectable for tests (default: the real mDNS browse).
var browseFunc = protocolDnssd.Browse

// DiscoverMdns browses `_wot._tcp` until the requested Thing is found, and
// returns the URL of its Thing Description.
//
// `instance` is matched against the DNS-SD instance name (the TD id without
// its "urn:" prefix, e.g. "eria:test:1"), then against the mDNS host label
// (the hostname without ".local.", e.g. "eria-aabbcc" for an ESP32 device).
//
// Private (LAN) addresses are accepted: ERIA targets a trusted local network.
// timeout <= 0 browses without deadline (used for background connections and
// re-discovery after a loss).
func DiscoverMdns(ctx context.Context, instance string, timeout time.Duration) (string, error) {
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	entries, err := browseFunc(ctx, protocolDnssd.BrowseOptions{
		PreferIPv4:   true,
		AllowPrivate: true,
	})
	if err != nil {
		return "", err
	}

	waiting := time.NewTicker(15 * time.Second)
	defer waiting.Stop()
	// Some announcers advertise unusable addresses first (e.g. a link-local
	// IPv6 of the loopback interface): keep such a match as fallback and
	// keep browsing for an entry with a routable address.
	fallback := ""
	for {
		select {
		case entry, ok := <-entries:
			if !ok {
				if fallback != "" {
					return fallback, nil
				}
				return "", fmt.Errorf("mDNS instance %q not found", instance)
			}
			if !matchInstance(entry, instance) {
				zlog.Debug().Str("mdns", instance).Str("instance", entry.Instance).Msg("[core:consumer:DiscoverMdns] ignoring non-matching entry")
				continue
			}
			if usableURL(entry) {
				zlog.Info().Str("mdns", instance).Str("url", entry.TDURL.String()).Msg("[core:consumer:DiscoverMdns] Thing found")
				return entry.TDURL.String(), nil
			}
			zlog.Debug().Str("mdns", instance).Str("url", entry.TDURL.String()).Msg("[core:consumer:DiscoverMdns] match with non-routable address, keeping as fallback")
			fallback = entry.TDURL.String()
		case <-waiting.C:
			zlog.Warn().Str("mdns", instance).Msg("[core:consumer:DiscoverMdns] still browsing...")
		case <-ctx.Done():
			if fallback != "" {
				return fallback, nil
			}
			return "", fmt.Errorf("mDNS instance %q not found: %w", instance, ctx.Err())
		}
	}
}

// usableURL reports whether the TD URL points to an address a consumer can
// actually reach (link-local and loopback addresses need a zone/same-host).
func usableURL(entry protocolDnssd.Entry) bool {
	if entry.TDURL == nil {
		return false
	}
	ip := net.ParseIP(strings.Trim(entry.TDURL.Hostname(), "[]"))
	if ip == nil {
		return true // hostname, let the connection decide
	}
	return !(ip.IsLinkLocalUnicast() || ip.IsLoopback() || ip.IsUnspecified())
}

func matchInstance(entry protocolDnssd.Entry, instance string) bool {
	if entry.Instance == instance {
		return true
	}
	host := strings.TrimSuffix(entry.Host, ".")
	if label, _, found := strings.Cut(host, "."); found {
		return label == instance
	}
	return host == instance
}
