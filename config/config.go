package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Config is used to configure the creation of a client
type Config struct {
	APIEndpointURL   string
	LoginEndpointURL string
	UAAEndpointURL   string

	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	UserAgent    string
	Origin       string
	Token        string

	BaseHTTPClient *http.Client

	skipSSLValidation bool
}

type cfHomeConfig struct {
	AccessToken           string
	RefreshToken          string
	Target                string
	AuthorizationEndpoint string
	OrganizationFields    struct {
		Name string
	}
	SpaceFields struct {
		Name string
	}
	SSLDisabled bool
}

func NewUserPassword(apiRootURL, username, password string) (*Config, error) {
	if username == "" {
		return nil, errors.New("expected an non-empty CF API username")
	}
	if password == "" {
		return nil, errors.New("expected an non-empty CF API password")
	}

	c, err := newDefault(apiRootURL)
	if err != nil {
		return nil, err
	}
	c.Username = username
	c.Password = password

	return c, nil
}

func NewClientSecret(apiRoot, clientID, clientSecret string) (*Config, error) {
	if clientID == "" {
		return nil, errors.New("expected an non-empty CF API clientID")
	}
	if clientSecret == "" {
		return nil, errors.New("expected an non-empty CF API clientSecret")
	}

	c, err := newDefault(apiRoot)
	if err != nil {
		return nil, err
	}
	c.ClientID = clientID
	c.ClientSecret = clientSecret

	return c, nil
}

func NewToken(apiRoot, token string) (*Config, error) {
	if token == "" {
		return nil, errors.New("expected an non-empty CF API token")
	}

	c, err := newDefault(apiRoot)
	if err != nil {
		return nil, err
	}
	c.Token = token

	return c, nil
}

func NewFromCFHome() (*Config, error) {
	dir, err := findCFHomeDir()
	if err != nil {
		return nil, err
	}
	return NewFromCFHomeDir(dir)
}

func NewFromCFHomeDir(cfHomeDir string) (*Config, error) {
	cfHomeConfig, err := loadCFHomeConfig(cfHomeDir)
	if err != nil {
		return nil, err
	}

	cfg, err := newDefault(cfHomeConfig.Target)
	if err != nil {
		return nil, err
	}
	cfg.Token = cfHomeConfig.AccessToken
	cfg.skipSSLValidation = cfHomeConfig.SSLDisabled

	return cfg, nil
}

func (c *Config) HTTPClient(httpClient *http.Client) {
	c.BaseHTTPClient = httpClient
	c.setHTTPClientSSLConfig()
}

func (c *Config) SkipSSLValidation(skip bool) {
	c.skipSSLValidation = skip
	c.setHTTPClientSSLConfig()
}

func (c *Config) setHTTPClientSSLConfig() {
	var tp *http.Transport
	switch t := c.BaseHTTPClient.Transport.(type) {
	case *http.Transport:
		tp = t
	case *oauth2.Transport:
		if bt, ok := t.Base.(*http.Transport); ok {
			tp = bt
		}
	}

	if tp != nil {
		if tp.TLSClientConfig == nil {
			tp.TLSClientConfig = &tls.Config{}
		}
		tp.TLSClientConfig.InsecureSkipVerify = c.skipSSLValidation
	}
}

func newDefault(apiRootURL string) (*Config, error) {
	u, err := url.ParseRequestURI(apiRootURL)
	if err != nil {
		return nil, fmt.Errorf("expected an http(s) CF API root URI, but got %s: %w", apiRootURL, err)
	}
	c := &Config{
		APIEndpointURL:    strings.TrimRight(u.String(), "/"),
		UserAgent:         "Go-CF-client/2.0",
		BaseHTTPClient:    http.DefaultClient,
		skipSSLValidation: false,
	}
	c.BaseHTTPClient.Transport = shallowDefaultTransport()
	c.setHTTPClientSSLConfig()
	return c, nil
}

func shallowDefaultTransport() *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	return &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
	}
}

func loadCFHomeConfig(cfHomeDir string) (*cfHomeConfig, error) {
	cfConfigDir := filepath.Join(cfHomeDir, ".cf")
	cfJSON, err := os.ReadFile(filepath.Join(cfConfigDir, "config.json"))
	if err != nil {
		return nil, err
	}

	var cfg cfHomeConfig
	err = json.Unmarshal(cfJSON, &cfg)
	if err == nil {
		if len(cfg.AccessToken) > len("bearer ") {
			cfg.AccessToken = cfg.AccessToken[len("bearer "):]
		}
	}

	return &cfg, nil
}

func findCFHomeDir() (string, error) {
	cfHomeDir := os.Getenv("CF_HOME")
	if cfHomeDir != "" {
		return cfHomeDir, nil
	}
	return os.UserHomeDir()
}
