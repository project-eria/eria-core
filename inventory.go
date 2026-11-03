package eria

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/auth"
	"github.com/project-eria/go-wot/inventory"
	zlog "github.com/rs/zerolog/log"
)

// The master's inventory (AUTH-DESIGN §2.4, F037): an app with a service
// token reads it, finds there the Things declared by `mdns:`, and sends a
// token only to the URLs it lists: the one auth issued for that very Thing
// (step 4e), worthless at any other.

// inventoryRefresh is how often the inventory is read again; discoverPoll,
// while a Thing is awaited.
var inventoryRefresh, discoverPoll = 30 * time.Second, 5 * time.Second

type inventoryCache struct {
	master   string
	token    func() string             // the app's token for the master
	forThing func(thing string) string // its token for a Thing
	client   *http.Client

	mu   sync.RWMutex
	list []inventory.Entry
}

var _inventory *inventoryCache // nil: no service token, mDNS as before

// startInventory reads the inventory once, then again in the background.
func startInventory(master string, tokens *auth.ServiceTokens) {
	c := &inventoryCache{master: master, token: tokens.Token, forThing: tokens.For, client: &http.Client{Timeout: 10 * time.Second}}
	c.read()
	go func() {
		for range time.Tick(inventoryRefresh) {
			c.read()
		}
	}()
	_inventory = c
	eriaconsumer.Token = c.tokenFor
}

func (c *inventoryCache) read() {
	token := c.token()
	if token == "" {
		return
	}
	list, err := inventory.Read(context.Background(), c.client, c.master, token)
	if err != nil {
		zlog.Warn().Err(err).Msg("[core:inventory] Can't read the master's inventory, the last one kept")
		return
	}
	c.mu.Lock()
	c.list = list
	c.mu.Unlock()
}

// tokenFor is the service token for the Thing at u, "" when u is not a
// Thing of the inventory.
func (c *inventoryCache) tokenFor(u string) string {
	e, ok := c.entryFor(u)
	if !ok {
		zlog.Debug().Str("url", u).Msg("[core:inventory] Not in the master's inventory: no token sent")
		return ""
	}
	return c.forThing(e.ID)
}

// entryFor finds the Thing u belongs to: same origin as its TD, the longest
// TD path u starts with (an app serves its Things under /<ref>, an ESP32
// under /).
func (c *inventoryCache) entryFor(u string) (inventory.Entry, bool) {
	origins := originsOf(inventory.Origin(u))
	path := strings.TrimSuffix(pathOf(u), "/")
	c.mu.RLock()
	defer c.mu.RUnlock()
	var best inventory.Entry
	length := -1
	for _, e := range c.list {
		p := strings.TrimSuffix(pathOf(e.URL), "/")
		if slices.Contains(origins, inventory.Origin(e.URL)) && (path == p || strings.HasPrefix(path, p+"/")) && len(p) > length {
			best, length = e, len(p)
		}
	}
	return best, length >= 0
}

// originsOf is origin, plus its addresses for a host name of the config: the
// master lists the addresses mDNS announces.
//
// ponytail: one DNS lookup per request to a named Thing; cache it if the
// traffic ever weighs. Trusts the network's DNS like the mDNS before it.
func originsOf(origin string) []string {
	if origin == "" {
		return nil
	}
	origins := []string{origin}
	scheme, hostport, _ := strings.Cut(origin, "://")
	host, port, err := net.SplitHostPort(hostport)
	if err != nil || net.ParseIP(host) != nil {
		return origins
	}
	ips, _ := net.LookupIP(host)
	for _, ip := range ips {
		origins = append(origins, scheme+"://"+net.JoinHostPort(ip.String(), port))
	}
	return origins
}

func pathOf(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return p.Path
}

// find matches the DNS-SD instance (TD id without "urn:"), then the mDNS host label.
func (c *inventoryCache) find(instance string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, e := range c.list {
		if e.ID == instance {
			return e.URL
		}
	}
	for _, e := range c.list {
		if e.Host == instance {
			return e.URL
		}
	}
	return ""
}

// discover finds a remote Thing declared by `mdns:`: in the inventory when
// the app reads it, else by browsing mDNS. timeout <= 0 waits without deadline.
func discover(instance string, timeout time.Duration) (string, error) {
	if _inventory == nil {
		return eriaconsumer.DiscoverMdns(context.Background(), instance, timeout)
	}
	deadline := time.Now().Add(timeout)
	for {
		if url := _inventory.find(instance); url != "" {
			return url, nil
		}
		if timeout > 0 && time.Now().After(deadline) {
			return "", fmt.Errorf("%q not in the master's inventory", instance)
		}
		time.Sleep(discoverPoll)
		_inventory.read()
	}
}
