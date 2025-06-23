package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/internal/jwt"
)

// Option is a functional option for configuring the client.
type Option func(*Config) error

// ClientCredentials is a functional option to set client credentials.
func ClientCredentials(clientID, clientSecret string) Option {
	return func(c *Config) error {
		// don't override the default client with empty
		if clientID = strings.TrimSpace(clientID); clientID != "" {
			c.clientID = clientID
		}

		// client/secret grant type takes precedence over nothing & token
		// but a secret must be set to be a real client
		if clientSecret = strings.TrimSpace(clientSecret); clientSecret != "" {
			c.clientSecret = clientSecret
		}
		return nil
	}
}

// ClientAssertion is a functional option to set client assertion.
func ClientAssertion(assertion string) Option {
	return func(c *Config) error {
		// if set, must be a valid JWT token
		// alternative to client secret
		// usually used with ClientCredentials
		// can be combined with JWTBearerAssertion
		// refer RFC 7523 for more details
		if assertion = strings.TrimSpace(assertion); assertion == "" {
			return errors.New("assertion must be valid JWT Token")
		} else {
			c.clientAssertion = assertion
		}

		return nil
	}
}

// UserPassword is a functional option to set user credentials.
func UserPassword(username, password string) Option {
	return func(c *Config) error {
		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)
		if username == "" || password == "" {
			return errors.New("username and password are required when using using user credentials")
		}
		c.username = username
		c.password = password
		return nil
	}
}

// JWTBearerAssertion is a functional option to set JWT Bearer credentials.
func JWTBearerAssertion(assertion string) Option {
	return func(c *Config) error {
		// if set, must be a valid JWT token
		// can be used alone
		// or with ClientCredentials
		// or with ClientCredentials + ClientAssertion
		// refer RFC 7523 for more details
		if assertion = strings.TrimSpace(assertion); assertion == "" {
			return errors.New("assertion is required for JWT Bearer grant type")
		} else {
			c.assertion = assertion
		}

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

// SSHOAuthClient configures a clientID used to request an SSH code.
func SSHOAuthClient(clientID string) Option {
	return func(c *Config) error {
		c.sshOAuthClient = clientID
		return nil
	}
}
