package eria

import (
	"context"
	"crypto/ed25519"
	"os"
	"sync"
	"time"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/auth"
	zlog "github.com/rs/zerolog/log"
)

// The app key and service token (AUTH-DESIGN §1.3, step 3c): with a key
// (ERIA_APP_KEY, set by the agent) and the guard config of an activated
// house, the app proves its key to auth and sends its service token to the
// Things it consumes.
//
// ponytail: the token goes to every Thing consumed, found by URL or mDNS,
// as eria-ui-server does; only to the master's inventory with step 3c-2
// (F037).

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
		a := eriaConfig.Guard
		key := appKey()
		if key == nil || a.Channel == "" || a.Keyset == "" {
			return
		}
		ks, err := loadKeyset(a.Keyset, a.Roots)
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:startServiceTokens] Can't load the auth keyset")
		}
		tokens, err := auth.NewServiceTokens(a.Channel, key, ks)
		if err != nil {
			zlog.Fatal().Err(err).Msg("[core:startServiceTokens] Can't prepare the service token")
		}
		eriaconsumer.Token = tokens.Token
		ready := make(chan struct{})
		go tokens.Run(context.Background(), ready)
		select {
		case <-ready:
		case <-time.After(30 * time.Second):
		}
	})
}

func loadKeyset(keysetFile, rootsFile string) (auth.Keyset, error) {
	data, err := os.ReadFile(keysetFile)
	if err != nil {
		return auth.Keyset{}, err
	}
	rootsPEM, err := os.ReadFile(rootsFile)
	if err != nil {
		return auth.Keyset{}, err
	}
	deploy, err := auth.DeploymentKeys(rootsPEM)
	if err != nil {
		return auth.Keyset{}, err
	}
	ks, _, err := auth.OpenKeyset(data, deploy)
	return ks, err
}
