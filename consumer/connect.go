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

// backoffSchedule is when and how often a TD read is retried: the wait
// after the first failure, after the second, etc. After the last one the
// read has failed. About two minutes: a neighbour refusing for want of a
// service token (auth restarting, or the neighbour in minimal mode) answers
// again by then.
var backoffSchedule = []time.Duration{
	1 * time.Second,
	5 * time.Second,
	10 * time.Second,
	15 * time.Second,
	30 * time.Second,
	30 * time.Second,
	30 * time.Second,
}

// transient says a status worth retrying: no token yet or refused for now
// (401, 403), auth's rate limit (429), a neighbour failing (5xx). A 403 can
// also be a lasting refusal: each retry logs it, so it stays visible.
func transient(status int) bool {
	return status == http.StatusUnauthorized || status == http.StatusForbidden ||
		status == http.StatusTooManyRequests || status >= 500
}

func BackoffMethod(url string) (*thing.Thing, error) {
	var resp *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		if t := token(url); t != "" { // asked again each time: it may have come
			req.Header.Set("Authorization", "Bearer "+t)
		}
		resp, err = http.DefaultClient.Do(req)
		switch {
		case err != nil:
			zlog.Warn().Str("url", url).Err(err).Msgf("[core:BackoffMethod] Retrying in %v", backoff)
		case transient(resp.StatusCode):
			zlog.Warn().Str("url", url).Str("status", resp.Status).Msgf("[core:BackoffMethod] Retrying in %v", backoff)
			resp.Body.Close()
		default:
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
		time.Sleep(backoff)
	}
	if err != nil {
		return nil, errors.New("can't open thing url")
	}
	zlog.Error().Str("status", resp.Status).Str("url", url).Msg("[core:BackoffMethod] incorrect response")
	return nil, errors.New("incorrect response")
}
