package config

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/internal/jwt"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"

	"github.com/stretchr/testify/require"
)

const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk"
const refreshToken = "secret-refresh-token"
const clientAssertion = "client-assertion-token"

func TestInvalidConfig(t *testing.T) {
	c := &Config{}
	err := c.Validate()
	require.Error(t, err)
	require.Equal(t, err, ErrConfigInvalid)
}

func TestUsernamePassword(t *testing.T) {
	t.Run("with empty username", func(t *testing.T) {
		_, err := New("https://api.example.com",
			UserPassword("", "test"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.EqualError(t, err, "username and password are required when using using user credentials")
	})

	t.Run("with empty password", func(t *testing.T) {
		_, err := New("https://api.example.com", UserPassword("user", ""))
		require.Error(t, err)
		require.EqualError(t, err, "username and password are required when using using user credentials")
	})

	t.Run("with username and password", func(t *testing.T) {
		// user/pass hits the token endpoint
		uaaURL := testutil.SetupFakeUAAServer(300)
		c, err := New("https://api.example.com",
			UserPassword("username", "password"),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, uaaURL, c.loginEndpointURL)
		require.Equal(t, uaaURL, c.uaaEndpointURL)
		require.Equal(t, "username", c.username)
		require.Equal(t, "password", c.password)
		require.Equal(t, "cf", c.clientID)
		require.Equal(t, GrantTypePassword, c.grantType)
	})

	t.Run("with username and password with non-default client", func(t *testing.T) {
		// user/pass hits the token endpoint
		uaaURL := testutil.SetupFakeUAAServer(300)
		c, err := New("https://api.example.com",
			UserPassword("username", "password"),
			ClientCredentials("clientID", ""),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "username", c.username)
		require.Equal(t, "password", c.password)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, GrantTypePassword, c.grantType)
	})
}

func TestClientCredentials(t *testing.T) {
	t.Run("with invalid URL", func(t *testing.T) {
		_, err := New(":", ClientCredentials("clientID", "clientSecret"))
		require.ErrorContains(t, err, "expected an http(s) CF API root URI, but got")
	})

	t.Run("with clientID and empty client secret", func(t *testing.T) {
		_, err := New("https://api.example.com",
			ClientCredentials("clientID", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.EqualError(t, err, "CF API credentials were not provided")
	})

	t.Run("with empty clientID", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("", "clientSecret"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "cf", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("with clientID and client secret", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", "clientSecret"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("with clientID and client secret and access token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", "clientSecret"),
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, accessToken, c.oAuthToken.AccessToken)
		require.Equal(t, refreshToken, c.oAuthToken.RefreshToken)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("with clientID and client assertion", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", ""),
			ClientAssertion(clientAssertion),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, clientAssertion, c.clientAssertion)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})
}

func TestToken(t *testing.T) {
	t.Run("with empty token", func(t *testing.T) {
		_, err := New("https://api.example.com",
			Token("", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid CF API token")
	})

	t.Run("with access token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, accessToken, c.oAuthToken.AccessToken)
		require.Equal(t, "", c.oAuthToken.RefreshToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with refresh token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token("", refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with access token and refresh token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with access token and custom clientID", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, ""),
			ClientCredentials("myapp", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with custom http.Client without a transport and skip TLS verification", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"), // skip service discovery
			SkipTLSValidation(),
			HttpClient(&http.Client{}))
		require.NoError(t, err)
		require.NotNil(t, c.HTTPClient())
		require.NotNil(t, c.HTTPClient().Transport)
		require.NotNil(t, c.HTTPClient().Transport.(*http.Transport).TLSClientConfig)
		require.True(t, c.HTTPClient().Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
	})
}

func TestNewConfigFromCFHomeDir(t *testing.T) {
	cfHomeDir := writeTestCFCLIConfig(t)

	t.Run("without overrides", func(t *testing.T) {
		cfg, err := NewFromCFHomeDir(cfHomeDir)
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, "https://login.sys.example.com", cfg.loginEndpointURL)
		require.Equal(t, "https://uaa.sys.example.com", cfg.uaaEndpointURL)
		require.Equal(t, accessToken, cfg.oAuthToken.AccessToken)
		require.Equal(t, refreshToken, cfg.oAuthToken.RefreshToken)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeRefreshToken, cfg.grantType)
	})

	t.Run("with with CF_USERNAME and CF_PASSWORD set", func(t *testing.T) {
		uaaURL := testutil.SetupFakeUAAServer(300)
		require.NoError(t, os.Setenv("CF_USERNAME", "admin"))
		require.NoError(t, os.Setenv("CF_PASSWORD", "pass"))
		defer func() {
			_ = os.Unsetenv("CF_USERNAME")
			_ = os.Unsetenv("CF_PASSWORD")
		}()
		cfg, err := NewFromCFHomeDir(cfHomeDir, AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, uaaURL, cfg.loginEndpointURL)
		require.Equal(t, uaaURL, cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypePassword, cfg.grantType)
	})

	t.Run("with override options", func(t *testing.T) {
		uaaURL := testutil.SetupFakeUAAServer(300)
		cfg, err := NewFromCFHomeDir(cfHomeDir,
			UserPassword("admin", "pass"),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, uaaURL, cfg.loginEndpointURL)
		require.Equal(t, uaaURL, cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypePassword, cfg.grantType)
	})
}

func TestJwBearerAssertion(t *testing.T) {
	// A syntactically valid JWT (three base64url segments). The fake UAA
	// server in testutil does not verify the assertion's signature or
	// claims, so this is sufficient for exercising the exchange flow.
	const validAssertion = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ0ZXN0In0.c2ln"

	t.Run("eager exchange populates oAuth state", func(t *testing.T) {
		uaaURL := testutil.SetupFakeUAAServer(300)
		cfg, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, validAssertion, cfg.assertion)
		require.Equal(t, GrantTypeJwtBearer, cfg.grantType)
	})

	t.Run("hands off to refresh-token flow when UAA issues a refresh_token", func(t *testing.T) {
		// SetupFakeUAAServer always returns a refresh_token. The token
		// source returned by CreateOAuth2TokenSource should therefore be
		// the standard oauth2 reuse-source over the refresh-token grant,
		// not a *jwt.JWTAssertionTokenSource.
		uaaURL := testutil.SetupFakeUAAServer(300)
		cfg, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)

		src, err := cfg.CreateOAuth2TokenSource(context.Background())
		require.NoError(t, err)
		_, isAssertionSrc := src.(*jwt.JWTAssertionTokenSource)
		require.False(t, isAssertionSrc, "after eager exchange the long-lived source must not be the assertion source")

		token, err := src.Token()
		require.NoError(t, err)
		require.NotEmpty(t, token.AccessToken)
		require.Equal(t, "barfoo", token.RefreshToken)
	})

	t.Run("caches access token when UAA does not issue a refresh_token", func(t *testing.T) {
		var hits []string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits = append(hits, r.Method+" "+r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"access_token":"a","token_type":"bearer","expires_in":600}`)
		}))
		defer ts.Close()

		cfg, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			AuthTokenURL(ts.URL, ts.URL))
		require.NoError(t, err)
		require.Equal(t, []string{"POST /oauth/token"}, hits, "eager exchange should hit UAA exactly once at config-init time")

		src, err := cfg.CreateOAuth2TokenSource(context.Background())
		require.NoError(t, err)
		// Cached token is reused; UAA must not be called again.
		token, err := src.Token()
		require.NoError(t, err)
		require.Equal(t, "a", token.AccessToken)
		token, err = src.Token()
		require.NoError(t, err)
		require.Equal(t, "a", token.AccessToken)
		require.Equal(t, []string{"POST /oauth/token", "POST /oauth/token"}, hits,
			"the second CreateOAuth2TokenSource performs a fresh exchange; "+
				"subsequent Token() calls on the cached source must not re-hit UAA")
	})

	t.Run("returns diagnostic when cached access token expires without a refresh_token", func(t *testing.T) {
		// expires_in:1 puts the token within oauth2's 10-second expiry buffer,
		// so ReuseTokenSource treats it as already-expired and falls through
		// to the diagnostic source on the very next Token() call.
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"access_token":"a","token_type":"bearer","expires_in":1}`)
		}))
		defer ts.Close()

		cfg, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			AuthTokenURL(ts.URL, ts.URL))
		require.NoError(t, err)

		src, err := cfg.CreateOAuth2TokenSource(context.Background())
		require.NoError(t, err)
		_, err = src.Token()
		require.Error(t, err)
		require.Contains(t, err.Error(), "no refresh token was issued by UAA")
	})

	t.Run("invalid assertion fails at config.New time", func(t *testing.T) {
		// UAA rejects the exchange. config.New must surface the error
		// rather than deferring it to the first authenticated API call.
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":"invalid_grant","error_description":"assertion expired"}`)
		}))
		defer ts.Close()

		_, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			AuthTokenURL(ts.URL, ts.URL))
		require.Error(t, err)
		require.Contains(t, err.Error(), "jwt-bearer assertion exchange failed")
	})

	t.Run("rejects ClientAssertion combined with JWTBearerAssertion", func(t *testing.T) {
		// The refresh-token handoff cannot replay a ClientAssertion: refreshes
		// would silently fall back to client_id/client_secret auth and fail at
		// runtime if UAA requires client_assertion. Reject up front.
		_, err := New("https://api.example.com",
			JWTBearerAssertion(validAssertion),
			ClientCredentials("clientID", ""),
			ClientAssertion(clientAssertion),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "client assertion is not supported with the JWT Bearer assertion grant")
	})
}
