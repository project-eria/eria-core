package eria

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/auth"
	"github.com/project-eria/go-wot/inventory"
	zlog "github.com/rs/zerolog/log"
)

// The master's inventory (AUTH-DESIGN §2.4, F037): an app with a service
// token reads it, finds there the Things declared by `mdns:`, and sends its
// token only to the URLs it lists.

// inventoryRefresh is how often the inventory is read again; discoverPoll,
// while a Thing is awaited.
var inventoryRefresh, discoverPoll = 30 * time.Second, 5 * time.Second

type inventoryCache struct {
	master string
	token  func() string // the app's service token
	client *http.Client

	mu      sync.RWMutex
	list    []inventory.Entry
	origins map[string]bool
}

var _inventory *inventoryCache // nil: no service token, mDNS as before

// startInventory reads the inventory once, then again in the background.
func startInventory(master string, tokens *auth.ServiceTokens) {
	c := &inventoryCache{master: master, token: tokens.Token, client: &http.Client{Timeout: 10 * time.Second}}
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
	origins := map[string]bool{}
	for _, e := range list {
		origins[inventory.Origin(e.URL)] = true
	}
	c.mu.Lock()
	c.list, c.origins = list, origins
	c.mu.Unlock()
}

// tokenFor is the service token for url, "" when url is not a Thing of the
// inventory.
func (c *inventoryCache) tokenFor(url string) string {
	if !c.listed(inventory.Origin(url)) {
		zlog.Debug().Str("url", url).Msg("[core:inventory] Not in the master's inventory: no token sent")
		return ""
	}
	return c.token()
}

// listed tells whether origin is a Thing of the inventory. The master lists
// the addresses mDNS announces: a host name of the config is resolved first.
//
// ponytail: one DNS lookup per request to a named Thing; cache it if the
// traffic ever weighs. Trusts the network's DNS like the mDNS before it.
func (c *inventoryCache) listed(origin string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if origin == "" || c.origins[origin] {
		return origin != ""
	}
	scheme, hostport, _ := strings.Cut(origin, "://")
	host, port, err := net.SplitHostPort(hostport)
	if err != nil || net.ParseIP(host) != nil {
		return false
	}
	ips, _ := net.LookupIP(host)
	for _, ip := range ips {
		if c.origins[scheme+"://"+net.JoinHostPort(ip.String(), port)] {
			return true
		}
	}
	return false
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
