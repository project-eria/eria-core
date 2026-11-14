package eriaconsumer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A neighbour refusing for want of a token (401) is retried until it answers;
// a 404 is final at once.
func TestBackoffMethodRetriesRefusals(t *testing.T) {
	saved := backoffSchedule
	backoffSchedule = []time.Duration{time.Millisecond, time.Millisecond, time.Millisecond}
	defer func() { backoffSchedule = saved }()

	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"id":"urn:eria:test:1","title":"Test"}`))
	}))
	defer srv.Close()
	td, err := BackoffMethod(srv.URL)
	require.NoError(t, err)
	assert.Equal(t, "urn:eria:test:1", td.ID)
	assert.Equal(t, 3, calls)

	calls = 0
	gone := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer gone.Close()
	_, err = BackoffMethod(gone.URL)
	assert.Error(t, err)
	assert.Equal(t, 1, calls, "a 404 is not retried")

	calls = 0
	refused := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusForbidden)
	}))
	defer refused.Close()
	_, err = BackoffMethod(refused.URL)
	assert.Error(t, err, "a lasting refusal fails after the schedule")
	assert.Equal(t, 3, calls)
}
