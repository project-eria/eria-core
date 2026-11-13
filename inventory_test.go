package eria

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/project-eria/go-wot/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventory(t *testing.T) {
	coordinator := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/inventory" || r.Header.Get("Authorization") != "Bearer svc" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode([]inventory.Entry{
			{ID: "eria:dummy:etage:bureau-volet", Host: "etage", URL: "http://etage/bureau-volet"},
			{ID: "eria:dummy:etage:bureau", Host: "etage", URL: "http://etage/bureau"},
			{ID: "eria:esp:1", Host: "eria-aabbcc", URL: "http://10.0.0.7/"},
			{ID: "eria:notices", URL: "http://127.0.0.1:8104/"},
		})
	}))
	defer coordinator.Close()

	token := ""
	c := &inventoryCache{coordinator: coordinator.URL + "/", token: func() string { return token },
		forThing: func(id string) string { return "for " + id }, client: coordinator.Client()}
	c.read()
	assert.Empty(t, c.find("eria:esp:1"), "no token, no inventory")
	token = "svc"
	c.read()

	// Step 4e: the token of the very Thing the URL belongs to.
	assert.Equal(t, "for eria:dummy:etage:bureau-volet", c.tokenFor("http://etage/bureau-volet"))
	assert.Equal(t, "for eria:dummy:etage:bureau-volet", c.tokenFor("ws://etage:80/bureau-volet/properties/open"), "its WebSocket")
	assert.Equal(t, "for eria:dummy:etage:bureau", c.tokenFor("http://etage/bureau/properties/on"), "a Thing whose path starts another's")
	assert.Empty(t, c.tokenFor("http://etage/salon"), "no Thing at that path")
	assert.Equal(t, "for eria:esp:1", c.tokenFor("http://10.0.0.7/properties/temp"))
	assert.Empty(t, c.tokenFor("http://10.9.9.9/bureau-volet"), "not in the inventory: no token")
	assert.Empty(t, c.tokenFor("http://10.0.0.7:8080/"), "another port of the same host")
	assert.Equal(t, "for eria:notices", c.tokenFor("http://localhost:8104/x"), "a name of the config, resolved to a listed address")
	assert.Empty(t, c.tokenFor("http://localhost:8105/x"))

	assert.Equal(t, "http://etage/bureau-volet", c.find("eria:dummy:etage:bureau-volet"))
	assert.Equal(t, "http://10.0.0.7/", c.find("eria-aabbcc"), "by mDNS host label")

	_inventory, discoverPoll = c, time.Millisecond
	defer func() { _inventory, discoverPoll = nil, 5*time.Second }()
	url, err := discover("eria-aabbcc", time.Second)
	require.NoError(t, err)
	assert.Equal(t, "http://10.0.0.7/", url)
	_, err = discover("unknown", time.Millisecond)
	assert.Error(t, err)
}
