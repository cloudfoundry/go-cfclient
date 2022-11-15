package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"time"
)

// OAuthSessionManager creates and manages OAuth http client instances
type OAuthSessionManager struct {
	config *config.Config

	tokenSource         oauth2.TokenSource
	tokenSourceDeadline *time.Time

	authenticatedHTTPClient *http.Client
}

// NewOAuthSessionManager creates a new OAuth session manager
func NewOAuthSessionManager(config *config.Config) *OAuthSessionManager {
	return &OAuthSessionManager{
		config: config,
	}
}

// Client returns an authenticated OAuth http client
func (m *OAuthSessionManager) Client() (*http.Client, error) {
	if m.shouldRenewToken() {
		err := m.refreshAuthenticatedHTTPClient()
		if err != nil {
			return nil, err
		}
	}
	return m.authenticatedHTTPClient, nil
}

// Token returns the OAuth token in "bearer: <token>" format
func (m *OAuthSessionManager) Token() (string, error) {
	if m.shouldRenewToken() {
		err := m.refreshAuthenticatedHTTPClient()
		if err != nil {
			return "", err
		}
	}

	token, err := m.tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error getting bearer token: %w", err)
	}
	return "bearer " + token.AccessToken, nil
}

// refreshAuthenticatedHTTPClient creates a new authenticated OAuth http client
func (m *OAuthSessionManager) refreshAuthenticatedHTTPClient() error {
	if m.config.LoginEndpointURL == "" || m.config.UAAEndpointURL == "" {
		return errors.New("login and UAA endpoints must not be empty")
	}

	loginEndpoint := path.Join(m.config.LoginEndpointURL, "/oauth/auth")
	uaaEndpoint := path.Join(m.config.UAAEndpointURL, "/oauth/token")

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, m.config.BaseHTTPClient)

	switch {
	case m.config.Token != "":
		m.userTokenAuth(ctx, loginEndpoint, uaaEndpoint)
	case m.config.ClientID != "":
		m.clientAuth(ctx, loginEndpoint)
	default:
		err := m.userAuth(ctx, loginEndpoint, uaaEndpoint)
		if err != nil {
			return err
		}
	}

	return nil
}

// userAuth initializes a http client using standard username and password
func (m *OAuthSessionManager) userAuth(ctx context.Context, loginEndpoint, uaaEndpoint string) error {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  loginEndpoint,
			TokenURL: uaaEndpoint,
		},
	}
	if m.config.Origin != "" {
		type LoginHint struct {
			Origin string `json:"origin"`
		}
		loginHint := LoginHint{m.config.Origin}
		origin, err := json.Marshal(loginHint)
		if err != nil {
			return fmt.Errorf("error creating login_hint for user auth: %w", err)
		}
		val := url.Values{}
		val.Set("login_hint", string(origin))
		authConfig.Endpoint.TokenURL = path.Format("%s?%s", authConfig.Endpoint.TokenURL, val)
	}

	token, err := authConfig.PasswordCredentialsToken(ctx, m.config.Username, m.config.Password)
	if err != nil {
		return fmt.Errorf("error getting token for user auth: %w", err)
	}

	m.tokenSourceDeadline = &token.Expiry
	m.tokenSource = authConfig.TokenSource(ctx, token)
	m.authenticatedHTTPClient = oauth2.NewClient(ctx, m.tokenSource)

	return nil
}

// clientAuth initializes a http client using OAuth client id and secret
func (m *OAuthSessionManager) clientAuth(ctx context.Context, uaaEndpoint string) {
	authConfig := &clientcredentials.Config{
		ClientID:     m.config.ClientID,
		ClientSecret: m.config.ClientSecret,
		TokenURL:     uaaEndpoint,
	}

	m.tokenSource = authConfig.TokenSource(ctx)
	m.authenticatedHTTPClient = authConfig.Client(ctx)
}

// userTokenAuth initializes client credentials from existing bearer token.
func (m *OAuthSessionManager) userTokenAuth(ctx context.Context, loginEndpoint, uaaEndpoint string) {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  loginEndpoint,
			TokenURL: uaaEndpoint,
		},
	}

	// Token is expected to have no "bearer" prefix
	token := &oauth2.Token{
		AccessToken: m.config.Token,
		TokenType:   "Bearer"}

	m.tokenSource = authConfig.TokenSource(ctx, token)
	m.tokenSourceDeadline = &token.Expiry
	m.authenticatedHTTPClient = oauth2.NewClient(ctx, m.tokenSource)
}

func (m *OAuthSessionManager) shouldRenewToken() bool {
	if m.tokenSourceDeadline != nil {
		expiresAt := time.Now()
		return m.tokenSourceDeadline.Before(expiresAt)
	}
	return true
}
