package eria

import (
	"context"
	"crypto/ed25519"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/project-eria/go-wot/auth"
	zlog "github.com/rs/zerolog/log"
)

// The app key and service token (AUTH-DESIGN §1.3, step 3c): with a key
// (ERIA_APP_KEY, set by the agent) and the guard config of an activated
// house, the app proves its key to auth and sends its service token to the
// Things it consumes.
//
// The token goes only to the Things of the master's inventory (F037,
// inventory.go): without `guard.master`, nowhere.

var (
	_appKeyOnce sync.Once
	_appKey     ed25519.PrivateKey
	_tokensOnce sync.Once
)

// appKey is the key the agent gave the app, nil without one.
func appKey() ed25519.PrivateKey {
	_appKeyOnce.Do(func() {
		var err error
		if _appKey, err = auth.LoadAppKey(); err != nil {
			zlog.Fatal().Err(err).Msg("[core:appKey] Can't load the app key")
		}
	})
	return _appKey
}

// startServiceTokens gets the first token, waiting at most for one try,
// then renews it in the background; the consumer sends it.
func startServiceTokens() {
	_tokensOnce.Do(func() {
		cfg, master := guardConfig()
		key := appKey()
		if key == nil || cfg.Channel == "" || DeployPath {
			return
		}
		ks, _, err := cfg.OpenKeyset()
		if errors.Is(err, os.ErrNotExist) {
			return // not activated
		}
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:startServiceTokens] Can't load the auth keyset")
		}
		tokens, err := auth.NewServiceTokens(cfg.Channel, key, ks)
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:startServiceTokens] Can't prepare the service token")
		}
		ready := make(chan struct{})
		go tokens.Run(context.Background(), ready)
		select {
		case <-ready:
		case <-time.After(30 * time.Second):
		}
		if master == "" {
			zlog.Warn().Msg("[core:startServiceTokens] No guard.master: the service token is sent nowhere")
			return
		}
		startInventory(master, tokens)
	})
}
