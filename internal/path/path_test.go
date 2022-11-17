package path

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPathFormat(t *testing.T) {
	qs := url.Values{}
	qs.Set("key", "val")
	guid := "GUID"

	require.Equal(t, "/v3/apps/GUID/env",
		Format("/v3/apps/%s/env", guid))

	require.Equal(t, "/v3/apps/GUID/env?key=val",
		Format("/v3/apps/%s/env?%s", guid, qs))

	require.Equal(t, "/v3/apps/GUID?key=val",
		Format("/v3/apps/%s?%s", guid, qs))

	require.Equal(t, "/v3/apps/GUID/env",
		Format("/v3/apps/%s/env?%s", guid, url.Values{}))
}

func TestPathJoin(t *testing.T) {
	type pathTest struct {
		parts    []string
		expected string
	}
	tests := []pathTest{
		{
			parts:    []string{"/v3/apps/", "GUID/env"},
			expected: "/v3/apps/GUID/env",
		},
		{
			parts:    []string{"/v3/apps", "GUID"},
			expected: "/v3/apps/GUID",
		},
		{
			parts:    []string{"/v3/apps/", "/GUID/env"},
			expected: "/v3/apps/GUID/env",
		},
		{
			parts:    []string{"/v3/apps/", "/GUID?key=val"},
			expected: "/v3/apps/GUID?key=val",
		},
		{
			parts:    []string{"https://api.example.org/v3/apps/", "/GUID/env"},
			expected: "https://api.example.org/v3/apps/GUID/env",
		},
		{
			parts:    []string{"https://api.example.org/v3/apps/", ""},
			expected: "https://api.example.org/v3/apps",
		},
	}
	for _, tt := range tests {
		require.Equal(t, tt.expected, Join(tt.parts...))
	}
}
