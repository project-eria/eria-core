package eria

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolveNetwork(t *testing.T) {
	os.Unsetenv("ERIA_NETWORK")
	network, err := resolveNetwork("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "dummy", network, "configured network without ERIA_NETWORK")

	t.Setenv("ERIA_NETWORK", "banc")
	network, err = resolveNetwork("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "banc", network, "ERIA_NETWORK wins over the config")

	t.Setenv("ERIA_NETWORK", "")
	network, err = resolveNetwork("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "", network, "an empty ERIA_NETWORK imposes the default network")

	t.Setenv("ERIA_NETWORK", "Bad_Name")
	_, err = resolveNetwork("dummy")
	assert.Error(t, err, "invalid network rejected")
}
