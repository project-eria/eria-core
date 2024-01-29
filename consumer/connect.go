package eriaconsumer

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/project-eria/go-wot/thing"
	zlog "github.com/rs/zerolog/log"
)

type ConnectMethod func(string) (*thing.Thing, error)

func BackoffMethod(url string) (*thing.Thing, error) {
	// A backoff schedule for when and how often to retry failed HTTP
	// requests. The first element is the time to wait after the
	// first failure, the second the time to wait after the second
	// failure, etc. After reaching the last element, retries stop
	// and the request is considered failed.
	var backoffSchedule = []time.Duration{
		1 * time.Second,
		5 * time.Second,
		5 * time.Second,
	}
	var err error
	var resp *http.Response
	for _, backoff := range backoffSchedule {
		resp, err = http.Get(url)
		if err == nil {
			break
		}

		zlog.Error().Str("url", url).Err(err).Msgf("[core:BackoffMethod] Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, errors.New("can't open thing url")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		zlog.Error().Str("status", resp.Status).Str("url", url).Msg("[core:BackoffMethod] incorrect response")
		return nil, errors.New("incorrect response")
	}

	var td thing.Thing
	if err := json.NewDecoder(resp.Body).Decode(&td); err != nil {
		zlog.Error().Str("url", url).Err(err).Msg("[core:BackoffMethod]")
		return nil, errors.New("can't decode json")
	}

	return &td, nil
}
