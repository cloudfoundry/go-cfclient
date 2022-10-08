package client

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPathBuilder(t *testing.T) {
	qs := url.Values{}
	qs.Set("key", "val")
	require.Equal(t, "/v3/apps/GUID/env", joinPath("/v3/apps", "GUID", "env"))
	require.Equal(t, "/v3/apps/GUID/env?key=val", joinPathAndQS(qs, "/v3/apps", "GUID", "env"))
}
