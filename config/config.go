package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

// Option is a functional option for configuring the client.
type Option func(*Config) error

// cfHomeConfig represents the CF Home configuration.
type cfHomeConfig struct {
	AccessToken           string
	RefreshToken          string
	Target                string
	AuthorizationEndpoint string
	UaaEndpoint           string
	UAAOAuthClient        string
	UAAOAuthClientSecret  string
	UAAGrantType          string
	SSHOAuthClient        string
	SSLDisabled           bool
}

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

	isValidated bool
}

// configureHTTPClient configures the base http client
func (c *Config) configureHTTPClient() {
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

func (c *Config) tokenServiceURLDiscovery() error {
	// Return immediately if URLs have already been discovered
	if strings.TrimSpace(c.loginEndpointURL) != "" && strings.TrimSpace(c.uaaEndpointURL) != "" {
		return nil
	}

	// Perform the URL discovery
	root, err := c.Root(context.Background())
	if err != nil {
		return fmt.Errorf("error while discovering token service URL: %w", err)
	}
	// Update the URLs based on the discovery results
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
	if c.isValidated {
		return nil
	}

	// Trim spaces only once for efficiency
	c.clientID = strings.TrimSpace(c.clientID)
	c.username = strings.TrimSpace(c.username)
	c.clientSecret = strings.TrimSpace(c.clientSecret)
	c.password = strings.TrimSpace(c.password)

	// Ensure at least one of clientID, username, or token is provided
	if c.clientID == "" && c.username == "" && c.oAuthToken == nil {
		return errors.New("either client credentials, user credentials, or tokens are required")
	}

	if c.oAuthToken != nil {
		c.grantType = internalhttp.GrantTypeNone
	}

	// If clientID is provided, check for clientSecret
	if c.clientID != "" {
		if c.clientSecret == "" {
			if c.clientID != internalhttp.DefaultClientId {
				return errors.New("client secret is required when using client credentials")
			}
		} else {
			c.grantType = internalhttp.GrantTypeClientCredentials
		}
	} else {
		// Set a default clientID if not provided
		c.clientID = internalhttp.DefaultClientId
	}

	// If username is provided, check for password
	if c.username != "" {
		if c.password == "" {
			return errors.New("password is required when using user credentials")
		}
		c.grantType = internalhttp.GrantTypePassword
	}

	c.configureHTTPClient()

	if err := c.tokenServiceURLDiscovery(); err != nil {
		return err
	}

	c.isValidated = true // Mark as validated

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
	case internalhttp.GrantTypeClientCredentials:
		return c.clientAuth(oauthCtx)
	case internalhttp.GrantTypePassword:
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
		r, err := c.Root(ctx)
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

// Root queries the global API root
func (c *Config) Root(ctx context.Context) (*resource.Root, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ToURL("/"), nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred while generating the request for the global API root: %w", err)
	}
	resp, err := c.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred while attempting to query the global API root: %w", err)
	}

	defer ios.Close(resp.Body)

	var root resource.Root
	if err := internalhttp.DecodeBody(resp, &root); err != nil {
		return nil, fmt.Errorf("failed to decode API root response: %w", err)
	}
	return &root, nil
}

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

// Useragent is a functional option to set user agent.
func Useragent(userAgent string) Option {
	return func(c *Config) error {
		if userAgent = strings.TrimSpace(userAgent); userAgent == "" {
			c.userAgent = internalhttp.DefaultUserAgent
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
			c.requestTimeout = internalhttp.DefaultRequestTimeout
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
		return nil
	}
}

// New creates a new Config with specified API root URL and options.
func New(apiRootURL string, options ...Option) (*Config, error) {
	u, err := url.Parse(apiRootURL)
	if err != nil {
		return nil, fmt.Errorf("expected an http(s) CF API root URI, but got %s: %w", apiRootURL, err)
	}
	return validateConfig(&Config{
		apiEndpointURL: strings.TrimRight(u.String(), "/"),
		userAgent:      internalhttp.DefaultUserAgent,
		requestTimeout: internalhttp.DefaultRequestTimeout,
	}, options...)
}

// NewFromCFHome is similar to NewToken but reads the access token from the CF_HOME config, which must
// exist and have a valid access token.
//
// This will use the currently configured CF_HOME env var if it exists, otherwise attempts to use the
// default CF_HOME directory.
func NewFromCFHome(options ...Option) (*Config, error) {
	dir, err := findCFHomeDir()
	if err != nil {
		return nil, err
	}
	return NewFromCFHomeDir(dir, options...)
}

// NewFromCFHomeDir is similar to NewToken but reads the access token from the config in the specified directory
// which must exist and have a valid access token.
func NewFromCFHomeDir(cfHomeDir string, options ...Option) (*Config, error) {
	config, err := loadConfigFromCFHome(cfHomeDir)
	if err != nil {
		return nil, err
	}
	return validateConfig(config, options...)
}

// loadConfigFromCFHome reads the CF Home configuration from the specified directory.
func loadConfigFromCFHome(cfHomeDir string) (*Config, error) {
	configFile := filepath.Join(filepath.Join(cfHomeDir, ".cf"), "config.json")
	cfJSON, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", configFile, err)
	}
	var cfgHome cfHomeConfig
	if err = json.Unmarshal(cfJSON, &cfgHome); err != nil {
		return nil, fmt.Errorf("error while unmarshalling CF home config: %w", err)
	}
	cfg := &Config{
		apiEndpointURL:    cfgHome.Target,
		loginEndpointURL:  cfgHome.AuthorizationEndpoint,
		uaaEndpointURL:    cfgHome.UaaEndpoint,
		clientID:          cfgHome.UAAOAuthClient,
		clientSecret:      cfgHome.UAAOAuthClientSecret,
		skipTLSValidation: cfgHome.SSLDisabled,
		grantType:         cfgHome.UAAGrantType,
		sshOAuthClient:    cfgHome.SSHOAuthClient,
		username:          os.Getenv("CF_USERNAME"),
		password:          os.Getenv("CF_PASSWORD"),
		userAgent:         internalhttp.DefaultUserAgent,
		requestTimeout:    internalhttp.DefaultRequestTimeout,
	}
	if oAuthToken, err := jwt.ToOAuth2Token(cfgHome.AccessToken, cfgHome.RefreshToken); err == nil {
		cfg.oAuthToken = oAuthToken
	}
	return cfg, nil
}

func validateConfig(cfg *Config, options ...Option) (*Config, error) {
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, err
		}
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating configuration: %w", err)
	}
	return cfg, nil
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
