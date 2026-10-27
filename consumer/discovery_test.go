package eriaconsumer

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/project-eria/go-wot/protocolDnssd"
	"github.com/stretchr/testify/assert"
)

func fakeBrowse(entries ...protocolDnssd.Entry) func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error) {
	return func(ctx context.Context, _ protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error) {
		out := make(chan protocolDnssd.Entry, len(entries))
		for _, e := range entries {
			out <- e
		}
		go func() {
			<-ctx.Done()
			close(out)
		}()
		return out, nil
	}
}

func entry(instance, host, rawURL string) protocolDnssd.Entry {
	u, _ := url.Parse(rawURL)
	return protocolDnssd.Entry{
		Instance: instance,
		Host:     host,
		TDURL:    u,
	}
}

func Test_DiscoverMdnsMatchInstance(t *testing.T) {
	defer func(orig func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error)) { browseFunc = orig }(browseFunc)
	browseFunc = fakeBrowse(
		entry("other:thing:1", "10-0-0-9.local.", "http://10.0.0.9:80/"),
		entry("eria:test:1", "192-168-1-42.local.", "http://192.168.1.42:8888/"),
	)
	url, err := DiscoverMdns(context.Background(), "eria:test:1", 2*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "http://192.168.1.42:8888/", url)
}

func Test_DiscoverMdnsMatchHostLabel(t *testing.T) {
	defer func(orig func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error)) { browseFunc = orig }(browseFunc)
	browseFunc = fakeBrowse(
		entry("eria:esp:1wire:aabbcc", "eria-aabbcc.local.", "http://192.168.1.77:80/"),
	)
	url, err := DiscoverMdns(context.Background(), "eria-aabbcc", 2*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "http://192.168.1.77:80/", url)
}

func Test_DiscoverMdnsIgnoresNonMatching(t *testing.T) {
	defer func(orig func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error)) { browseFunc = orig }(browseFunc)
	browseFunc = fakeBrowse(
		entry("other:thing:1", "10-0-0-9.local.", "http://10.0.0.9:80/"),
	)
	_, err := DiscoverMdns(context.Background(), "eria:test:1", 300*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func Test_DiscoverMdnsTimeout(t *testing.T) {
	defer func(orig func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error)) { browseFunc = orig }(browseFunc)
	browseFunc = fakeBrowse() // nothing on the network
	start := time.Now()
	_, err := DiscoverMdns(context.Background(), "eria:test:1", 300*time.Millisecond)
	assert.Error(t, err)
	assert.Less(t, time.Since(start), 2*time.Second)
}

func Test_DiscoverMdnsBrowsesItsNetwork(t *testing.T) {
	defer func(orig func(context.Context, protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error)) { browseFunc = orig }(browseFunc)
	defer SetNetwork("")
	SetNetwork("dummy")
	var browsed string
	fake := fakeBrowse(entry("eria:test:1", "192-168-1-42.local.", "http://192.168.1.42:8888/"))
	browseFunc = func(ctx context.Context, opts protocolDnssd.BrowseOptions) (<-chan protocolDnssd.Entry, error) {
		browsed = opts.Network
		return fake(ctx, opts)
	}
	_, err := DiscoverMdns(context.Background(), "eria:test:1", 2*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "dummy", browsed)
}
