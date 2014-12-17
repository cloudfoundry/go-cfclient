package cfclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Client used to communicate with Clod Foundry
type Client struct {
	config Config
}

//Config is used to configure the creation of a client
type Config struct {
	ApiAddress   string
	LoginAddress string
	Username     string
	Password     string
	HttpClient   *http.Client
	Token        string
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

type AuthResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

//DefaultConfig configuration for client
func DefaultConfig() *Config {
	return &Config{
		ApiAddress:   "https://api.10.244.0.34.xip.io",
		LoginAddress: "https://login.10.244.0.34.xip.io",
		Username:     "admin",
		Password:     "admin",
		Token:        "",
		HttpClient:   http.DefaultClient,
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
		config.HttpClient = defConfig.HttpClient
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

	req.Header.Set("Authorization", "bearer "+c.config.Token)
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
	if c.config.Token == "" {
		c.config.Token = c.getToken()
	}
	return c.config.Token
}

func (c *Client) getToken() string {
	var authResp AuthResp
	data := url.Values{
		"grant_type": {"password"},
		"scope":      {""},
		"username":   {c.config.Username},
		"password":   {c.config.Password},
	}
	req, err := http.NewRequest("POST", c.config.LoginAddress+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error posting token %v\n", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")))

	resp, err := c.config.HttpClient.Do(req)
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response %v\n", err)
	}
	err = json.Unmarshal(resBody, &authResp)
	if err != nil {
		log.Printf("Error unmarshalling %v\n", err)
	}
	return "bearer " + authResp.AccessToken
}
