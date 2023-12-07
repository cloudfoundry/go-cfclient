package config

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	internalhttp "github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/ios"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/jwt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

const (
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeAuthorizationCode = "authorization_code"

	DefaultRequestTimeout = 30 * time.Second
	DefaultUserAgent      = "Go-CF-Client/3.0"
	DefaultClientID       = "cf"
)

// Config is used to configure the creation of a client
type Config struct {
	mu sync.Mutex

	apiEndpointURL   string
	loginEndpointURL string
	uaaEndpointURL   string
	sshOAuthClient   string

	username          string
	password          string
	clientID          string
	clientSecret      string
	grantType         string
	origin            string
	scopes            []string
	oAuthToken        *oauth2.Token
	httpClient        *http.Client
	skipTLSValidation bool
	requestTimeout    time.Duration
	userAgent         string
}

// New creates a new Config with specified API root URL and options.
func New(apiRootURL string, options ...Option) (*Config, error) {
	u, err := url.Parse(apiRootURL)
	if err != nil {
		return nil, fmt.Errorf("expected an http(s) CF API root URI, but got %s: %w", apiRootURL, err)
	}
	cfg := &Config{
		apiEndpointURL: strings.TrimRight(u.String(), "/"),
		userAgent:      DefaultUserAgent,
		requestTimeout: DefaultRequestTimeout,
		clientID:       DefaultClientID,
	}
	err = initConfig(cfg, options...)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewFromCFHome creates a go-cfclient config from the CF CLI config.
//
// This will use the currently configured CF_HOME env var if it exists, otherwise attempts to use the
// default CF_HOME directory.
//
// If CF_USERNAME and CF_PASSWORD env vars are set then those credentials will be used to get an oauth2 token. If
// those env vars are not set then the stored oauth2 token is used.
func NewFromCFHome(options ...Option) (*Config, error) {
	dir, err := findCFHomeDir()
	if err != nil {
		return nil, err
	}
	return NewFromCFHomeDir(dir, options...)
}

// NewFromCFHomeDir creates a go-cfclient config from the CF CLI config using the specified directory.
//
// This will attempt to read the CF CLI config from the specified directory only.
//
// If CF_USERNAME and CF_PASSWORD env vars are set then those credentials will be used to get an oauth2 token. If
// those env vars are not set then the stored oauth2 token is used.
func NewFromCFHomeDir(cfHomeDir string, options ...Option) (*Config, error) {
	cfg, err := createConfigFromCFCLIConfig(cfHomeDir)
	if err != nil {
		return nil, err
	}
	err = initConfig(cfg, options...)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// initConfig fully populates and then validates the provided base config
func initConfig(cfg *Config, options ...Option) error {
	// Apply any user provided config overrides
	err := applyOptions(cfg, options...)
	if err != nil {
		return err
	}

	// Ensure an HTTP client is available and then query the CF API for UAA/Login endpoints
	configureHTTPClient(cfg)
	err = tokenServiceURLDiscovery(context.Background(), cfg)
	if err != nil {
		return err
	}

	// Ensure the config object is valid for use by the client
	return cfg.Validate()
}

// configureHTTPClient configures the base http client
func configureHTTPClient(c *Config) {
	// Only configure the client if it has not been configured before
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		}
	}
	if transport := getHTTPTransport(c.httpClient); transport != nil {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = c.skipTLSValidation
	}
	c.httpClient.CheckRedirect = internalhttp.CheckRedirect
	c.httpClient.Timeout = c.requestTimeout
}

func tokenServiceURLDiscovery(ctx context.Context, c *Config) error {
	// Return immediately if URLs have already been configured
	if strings.TrimSpace(c.loginEndpointURL) != "" && strings.TrimSpace(c.uaaEndpointURL) != "" {
		return nil
	}

	// Query the CF API root for the service locator records
	root, err := c.globalAPIRoot(ctx)
	if err != nil {
		return fmt.Errorf("error while discovering token service URL: %w", err)
	}
	c.loginEndpointURL = root.Links.Login.Href
	c.uaaEndpointURL = root.Links.Uaa.Href
	c.sshOAuthClient = root.Links.AppSSH.Meta.OauthClient
	return nil
}

// userAuth authenticates using user credentials or refreshes an existing token.
func (c *Config) userAuth(ctx context.Context, shouldCreateToken bool) error {
	authConfig := &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Scopes:       c.scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.loginEndpointURL + "/oauth/auth",
			TokenURL: c.uaaEndpointURL + "/oauth/token",
		},
	}
	// If shouldCreateToken is true, generate a new token
	if shouldCreateToken {
		if c.origin != "" {
			// Add login hint to the token URL
			authConfig.Endpoint.TokenURL = addLoginHintToURL(authConfig.Endpoint.TokenURL, c.origin)
		}
		// Authenticate with user credentials
		return c.obtainUserToken(ctx, authConfig)
	}
	// Refresh an existing token
	return c.refreshToken(ctx, authConfig)
}

// obtainUserToken obtains a user authentication token.
func (c *Config) obtainUserToken(ctx context.Context, authConfig *oauth2.Config) error {
	var err error
	if c.oAuthToken, err = authConfig.PasswordCredentialsToken(ctx, c.username, c.password); err != nil {
		return fmt.Errorf("an error occurred while obtaining the user authentication token: %w", err)
	}
	return nil
}

// refreshToken refreshes an existing token.
func (c *Config) refreshToken(ctx context.Context, authConfig *oauth2.Config) error {
	var err error
	if c.oAuthToken, err = authConfig.TokenSource(ctx, c.oAuthToken).Token(); err != nil {
		var oauthErr *oauth2.RetrieveError
		if errors.As(err, &oauthErr) && oauthErr.Response.StatusCode == http.StatusUnauthorized {
			if err = c.generateToken(ctx, true); err == nil {
				return nil
			}
		}
		return fmt.Errorf("an error occurred while attempting to refresh the token: %w", err)
	}
	return nil
}

// clientAuth authenticates using client credentials.
func (c *Config) clientAuth(ctx context.Context) error {
	authConfig := &clientcredentials.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Scopes:       c.scopes,
		TokenURL:     c.uaaEndpointURL + "/oauth/token",
	}
	var err error
	// Authenticate with client credentials
	if c.oAuthToken, err = authConfig.Token(ctx); err != nil {
		return fmt.Errorf("an error occurred while obtaining the client authentication token: %w", err)
	}
	return nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Ensure at least one of clientID, username, or token is provided
	if c.clientID == "" && c.username == "" && c.oAuthToken == nil {
		return errors.New("either client credentials, user credentials, or tokens are required")
	}

	// If a non-default clientID is provided, check for clientSecret
	if c.clientID != DefaultClientID && c.clientSecret == "" {
		return errors.New("client secret is required when using client credentials")
	}

	// If username is provided, check for password
	if c.username != "" && c.password == "" {
		return errors.New("password is required when using user credentials")
	}

	return nil
}

// token ensures the token is valid, refreshing it if necessary.
func (c *Config) token(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.oAuthToken.Valid() {
		return nil
	}
	// If the current token is not valid, refresh it.
	shouldCreateToken := c.oAuthToken == nil || strings.TrimSpace(c.oAuthToken.RefreshToken) == ""
	// If the token is missing or has no refresh token, re-authenticate.
	return c.generateToken(ctx, shouldCreateToken)
}

// generateToken generates a new token or refreshes an existing one.
func (c *Config) generateToken(ctx context.Context, shouldCreateToken bool) error {
	oauthCtx := context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
	if !shouldCreateToken {
		return c.userAuth(oauthCtx, shouldCreateToken)
	}
	switch c.grantType {
	case GrantTypeClientCredentials:
		return c.clientAuth(oauthCtx)
	case GrantTypeAuthorizationCode:
		return c.userAuth(oauthCtx, shouldCreateToken)
	default:
		return fmt.Errorf("unsupported grant type: `%s`", c.grantType)
	}
}

func (c *Config) ToURL(urlPath string) string {
	return path.Join(c.apiEndpointURL, urlPath)
}

func (c *Config) ToAuthenticateURL(urlPath string) string {
	return path.Join(c.uaaEndpointURL, urlPath)
}

func (c *Config) SSHOAuthClient(ctx context.Context) (string, error) {
	if c.sshOAuthClient == "" {
		r, err := c.globalAPIRoot(ctx)
		if err != nil {
			return "", err
		}
		c.sshOAuthClient = r.Links.AppSSH.Meta.OauthClient
	}
	return c.sshOAuthClient, nil
}

func (c *Config) executeHTTPRequest(req *http.Request, includeAuthHeader bool) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("request is empty or invalid")
	}

	req.Header.Set("User-Agent", c.userAgent)
	if includeAuthHeader {
		// Get a new access token if the current one is invalid.
		if err := c.token(req.Context()); err != nil {
			return nil, fmt.Errorf("unable to get new access token: %w", err)
		}
		// Set the OAuth header on the request.
		c.oAuthToken.SetAuthHeader(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed during HTTP request send: %w", err)
	}

	// If the response status code is 401 Unauthorized, refresh the access token and retry the request.
	if includeAuthHeader && resp.StatusCode == http.StatusUnauthorized {
		ios.Close(resp.Body)
		// Return from the function to retry the request.
		return c.reAuthenticateAndRetry(req)
	}

	if internalhttp.IsStatusSuccess(resp.StatusCode) {
		return resp, err
	}
	return nil, internalhttp.DecodeError(resp)
}

// ExecuteAuthRequest executes an HTTP request with authentication.
func (c *Config) ExecuteAuthRequest(req *http.Request) (*http.Response, error) {
	return c.executeHTTPRequest(req, true)
}

func (c *Config) ExecuteRequest(req *http.Request) (*http.Response, error) {
	return c.executeHTTPRequest(req, false)
}

func (c *Config) reAuthenticateAndRetry(req *http.Request) (*http.Response, error) {
	// Lock the mutex to prevent concurrent access to the OAuth token.
	c.mu.Lock()
	// Check if another goroutine already refreshed the token.
	if !c.oAuthToken.Valid() {
		// Set the OAuth token to nil so that a new one will be obtained on the next request.
		c.oAuthToken = nil
	}
	c.mu.Unlock()
	// Retry the request with the new token.
	return c.ExecuteAuthRequest(req)
}

// globalAPIRoot queries the CF API service discovery root endpoint
func (c *Config) globalAPIRoot(ctx context.Context) (*resource.Root, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ToURL("/"), nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred while generating the request for the global API root: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed during HTTP request send: %w", err)
	}
	if !internalhttp.IsStatusSuccess(resp.StatusCode) {
		return nil, internalhttp.DecodeError(resp)
	}
	defer ios.Close(resp.Body)

	var root resource.Root
	if err := internalhttp.DecodeBody(resp, &root); err != nil {
		return nil, fmt.Errorf("failed to decode API root response: %w", err)
	}
	return &root, nil
}

// createConfigFromCFCLIConfig reads the CF Home configuration from the specified directory.
func createConfigFromCFCLIConfig(cfHomeDir string) (*Config, error) {
	cf, err := loadCFCLIConfig(cfHomeDir)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		apiEndpointURL:    cf.Target,
		loginEndpointURL:  cf.AuthorizationEndpoint,
		uaaEndpointURL:    cf.UaaEndpoint,
		clientID:          cf.UAAOAuthClient,
		clientSecret:      cf.UAAOAuthClientSecret,
		skipTLSValidation: cf.SSLDisabled,
		sshOAuthClient:    cf.SSHOAuthClient,
		userAgent:         DefaultUserAgent,
		requestTimeout:    DefaultRequestTimeout,
	}

	// if the username and password are specified via env vars use password based auth
	if os.Getenv("CF_USERNAME") != "" && os.Getenv("CF_PASSWORD") != "" {
		cfg.username = os.Getenv("CF_USERNAME")
		cfg.password = os.Getenv("CF_PASSWORD")
		cfg.grantType = GrantTypeAuthorizationCode
	} else {
		oAuthToken, err := jwt.ToOAuth2Token(cf.AccessToken, cf.RefreshToken)
		if err != nil {
			return nil, err
		}
		cfg.oAuthToken = oAuthToken
		cfg.grantType = GrantTypeRefreshToken
	}

	return cfg, nil
}

func applyOptions(cfg *Config, options ...Option) error {
	for _, option := range options {
		if err := option(cfg); err != nil {
			return err
		}
	}
	return nil
}

// findCFHomeDir finds the CF Home directory.
func findCFHomeDir() (string, error) {
	cfHomeDir := os.Getenv("CF_HOME")
	if cfHomeDir != "" {
		return cfHomeDir, nil
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user's home directory: %w", err)
	}
	return userHomeDir, nil
}

func getHTTPTransport(client *http.Client) *http.Transport {
	switch t := client.Transport.(type) {
	case *http.Transport:
		return t
	case *oauth2.Transport:
		if httpTransport, ok := t.Base.(*http.Transport); ok {
			return httpTransport
		}
	}
	return nil
}

func addLoginHintToURL(tokenURL, origin string) string {
	u, err := url.Parse(tokenURL)
	if err != nil {
		// Handle the error, or return the original URL
		return tokenURL
	}

	q := u.Query()
	q.Add("login_hint", fmt.Sprintf(`{"origin":"%s"}`, origin))
	u.RawQuery = q.Encode()

	return u.String()
}
