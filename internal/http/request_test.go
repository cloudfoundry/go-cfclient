package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Run("Test IsIgnoredRedirect with default request", func(t *testing.T) {
		r, _ := http.NewRequestWithContext(context.Background(), "GET", "/v3/apps", nil)
		require.False(t, IsIgnoredRedirect(r))
	})

	t.Run("Test CheckRedirect with default request", func(t *testing.T) {
		r, _ := http.NewRequestWithContext(context.Background(), "GET", "/v3/apps", nil)
		require.Nil(t, CheckRedirect(r, nil))
	})

	t.Run("Test CheckRedirect with max redirects exceeded", func(t *testing.T) {
		r, _ := http.NewRequestWithContext(context.Background(), "GET", "/v3/apps", nil)
		err := CheckRedirect(r, make([]*http.Request, 11))
		require.Error(t, err)
		require.Equal(t, ErrMaxRedirects, err.Error())
	})

	t.Run("Test IgnoreRedirect and IsIgnoredRedirect", func(t *testing.T) {
		r, _ := http.NewRequestWithContext(context.Background(), "GET", "/v3/apps", nil)
		r = IgnoreRedirect(r)
		require.True(t, IsIgnoredRedirect(r))

		err := CheckRedirect(r, nil)
		require.ErrorIs(t, err, http.ErrUseLastResponse)
	})
}
