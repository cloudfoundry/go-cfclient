package cfclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

//Client used to communicate with Clod Foundry
type Client struct {
	config Config
}

//Config is used to configure the creation of a client
type Config struct {
	ApiAddress        string
	LoginAddress      string
	Username          string
	Password          string
	SkipSslValidation bool
	HttpClient        *http.Client
	Token             string
	TokenSource       oauth2.TokenSource
}

// request is used to help build up a request
type request struct {
	config *Config
	method string
	url    string
	params url.Values
	body   io.Reader
	obj    interface{}
}

//DefaultConfig configuration for client
func DefaultConfig() *Config {
	return &Config{
		ApiAddress:        "https://api.10.244.0.34.xip.io",
		LoginAddress:      "https://login.10.244.0.34.xip.io",
		Username:          "admin",
		Password:          "admin",
		Token:             "",
		SkipSslValidation: false,
		HttpClient:        http.DefaultClient,
	}
}

// NewClient returns a new client
func NewClient(config *Config) *Client {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.ApiAddress) == 0 {
		config.ApiAddress = defConfig.ApiAddress
	}

	if len(config.LoginAddress) == 0 {
		config.LoginAddress = defConfig.LoginAddress
	}

	if len(config.Username) == 0 {
		config.Username = defConfig.Username
	}

	if len(config.Password) == 0 {
		config.Password = defConfig.Password
	}

	if len(config.Token) == 0 {
		config.Token = defConfig.Token
	}

	if config.HttpClient == nil {
		if config.SkipSslValidation == false {
			config.HttpClient = defConfig.HttpClient
		} else {

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			config.HttpClient = &http.Client{Transport: tr}
		}
	}

	client := &Client{
		config: *config,
	}
	return client
}

// newRequest is used to create a new request
func (c *Client) newRequest(method, path string) *request {
	r := &request{
		config: &c.config,
		method: method,
		url:    c.config.ApiAddress + path,
		params: make(map[string][]string),
	}
	return r
}

// doRequest runs a request with our client
func (c *Client) doRequest(r *request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}

	c.GetToken()

	req.Header.Set("Authorization", c.GetToken())
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.body = b
		}
	}

	// Create the HTTP request
	return http.NewRequest(r.method, r.url, r.body)
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) GetToken() string {
	if c.config.Token != "" {
		return "bearer" + c.config.Token
	}

	if c.config.TokenSource == nil {
		tokenSource, err := c.getTokenSource()
		if err != nil {
			log.Printf("Error getting token %v\n", err)
		}

		c.config.TokenSource = tokenSource
	}

	token, err := c.config.TokenSource.Token()
	if err != nil {
		log.Printf("Error getting token %v\n", err)
		return ""
	}

	return "bearer " + token.AccessToken
}

func (c *Client) getTokenSource() (oauth2.TokenSource, error) {
	authConf := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.config.LoginAddress + "/oauth/auth",
			TokenURL: c.config.LoginAddress + "/oauth/token",
		},
	}

	token, err := authConf.PasswordCredentialsToken(
		oauth2.NoContext, c.config.Username, c.config.Password)

	if err != nil {
		log.Printf("Error getting token %v\n", err)
		return nil, err
	}

	tokenSource := authConf.TokenSource(oauth2.NoContext, token)
	return tokenSource, nil
}
