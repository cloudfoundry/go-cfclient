package client_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientWithInvalidConfig(t *testing.T) {
	_, err := client.New(nil)
	require.Error(t, err)
	require.EqualError(t, err, "config is nil")

	cfg := &config.Config{}
	_, err = client.New(cfg)
	require.Error(t, err)
	require.Equal(t, config.ErrConfigInvalid, err)
}
