package client

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPathBuilder(t *testing.T) {
	qs := url.Values{}
	qs.Set("key", "val")
	guid := "57b816f7-bbac-49d7-a3b6-2f342d676997"

	require.Equal(t, "/v3/apps/57b816f7-bbac-49d7-a3b6-2f342d676997/env",
		path("/v3/apps/%s/env", guid))

	require.Equal(t, "/v3/apps/57b816f7-bbac-49d7-a3b6-2f342d676997/env?key=val",
		path("/v3/apps/%s/env?%s", guid, qs))

	require.Equal(t, "/v3/apps/57b816f7-bbac-49d7-a3b6-2f342d676997?key=val",
		path("/v3/apps/%s?%s", guid, qs))

	require.Equal(t, "/v3/apps/57b816f7-bbac-49d7-a3b6-2f342d676997/env",
		path("/v3/apps/%s/env?%s", guid, url.Values{}))
}
