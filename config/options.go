package config

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/jwt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Option is a functional option for configuring the client.
type Option func(*Config) error

// ClientCredentials is a functional option to set client credentials.
func ClientCredentials(clientId, clientSecret string) Option {
	return func(c *Config) error {
		if clientId = strings.TrimSpace(clientId); clientId == "" {
			return errors.New("expected a non-empty CF API clientID")
		}
		if clientSecret = strings.TrimSpace(clientSecret); clientSecret == "" {
			return errors.New("expected a non-empty CF API clientSecret")
		}
		c.clientID = clientId
		c.clientSecret = clientSecret
		c.grantType = GrantTypeClientCredentials
		return nil
	}
}

// UserPassword is a functional option to set user credentials.
func UserPassword(username, password string) Option {
	return func(c *Config) error {
		if username = strings.TrimSpace(username); username == "" {
			return errors.New("expected a non-empty CF API username")
		}
		if password = strings.TrimSpace(password); password == "" {
			return errors.New("expected a non-empty CF API password")
		}
		c.username = username
		c.password = password
		c.grantType = GrantTypeAuthorizationCode
		return nil
	}
}

// Scopes is a functional option to set scopes.
func Scopes(scopes ...string) Option {
	return func(c *Config) error {
		c.scopes = scopes
		return nil
	}
}

// UserAgent is a functional option to set user agent.
func UserAgent(userAgent string) Option {
	return func(c *Config) error {
		if userAgent = strings.TrimSpace(userAgent); userAgent == "" {
			c.userAgent = DefaultUserAgent
		} else {
			c.userAgent = userAgent
		}
		return nil
	}
}

// Origin is a functional option to set the origin.
func Origin(origin string) Option {
	return func(c *Config) error {
		c.origin = origin
		return nil
	}
}

// AuthTokenURL is a functional option to set the authorize and token url.
func AuthTokenURL(loginURL, tokenURL string) Option {
	return func(c *Config) error {
		l, err := url.Parse(loginURL)
		if err != nil {
			return fmt.Errorf("expected an http(s) CF login URI, but got %s: %w", loginURL, err)
		}
		c.loginEndpointURL = strings.TrimRight(l.String(), "/")

		t, err := url.Parse(tokenURL)
		if err != nil {
			return fmt.Errorf("expected an http(s) CF token URI, but got %s: %w", tokenURL, err)
		}
		c.uaaEndpointURL = strings.TrimRight(t.String(), "/")
		return nil
	}
}

// HttpClient is a functional option to set the HTTP client.
func HttpClient(client *http.Client) Option {
	return func(c *Config) error {
		c.httpClient = client
		return nil
	}
}

// RequestTimeout is a functional option to set the request timeout.
func RequestTimeout(timeout time.Duration) Option {
	return func(c *Config) error {
		if timeout <= 0 {
			c.requestTimeout = DefaultRequestTimeout
		} else {
			c.requestTimeout = timeout
		}
		return nil
	}
}

// SkipTLSValidation is a functional option to skip TLS validation.
func SkipTLSValidation() Option {
	return func(c *Config) error {
		c.skipTLSValidation = true
		return nil
	}
}

// Token is a functional option to set the access and refresh tokens.
func Token(accessToken, refreshToken string) Option {
	return func(c *Config) error {
		oAuthToken, err := jwt.ToOAuth2Token(accessToken, refreshToken)
		if err != nil {
			return fmt.Errorf("invalid CF API token: %w", err)
		}
		c.oAuthToken = oAuthToken
		c.grantType = GrantTypeRefreshToken
		return nil
	}
}
