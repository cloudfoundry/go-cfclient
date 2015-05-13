package cfclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"net/url"
)

//Client used to communicate with Cloud Foundry
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

	ctx := oauth2.NoContext

	if config.SkipSslValidation == false {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, defConfig.HttpClient)
	} else {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: tr})

	}

	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.LoginAddress + "/oauth/auth",
			TokenURL: config.LoginAddress + "/oauth/token",
		},
	}

	token, err := authConfig.PasswordCredentialsToken(ctx, config.Username, config.Password)

	if err != nil {
		log.Printf("Error getting token %v\n", err)
	}

	config.HttpClient = authConfig.Client(ctx, token)
	config.TokenSource = authConfig.TokenSource(ctx, token)

	client := &Client{
		config: *config,
	}
	return client
}

// newRequest is used to create a new request
func (c *Client) newRequest(method, path string) *request {
	r := &request{
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
	token, err := c.config.TokenSource.Token()
	if err != nil {
		log.Printf("Error getting token %v\n", err)
		return ""
	}

	return "bearer " + token.AccessToken
}
