package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"io"
	"net/http"
	"net/url"
	"strings"

	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client used to communicate with Cloud Foundry
type Client struct {
	Config   Config
	Endpoint Endpoint

	common commonClient // Reuse a single struct instead of allocating one for each commonClient on the heap.

	Organizations             *OrgClient
	Applications              *AppClient
	AppFeatures               *AppFeatureClient
	AppUsageEvents            *AppUsageClient
	AuditEvents               *AuditEventClient
	Buildpacks                *BuildpackClient
	Builds                    *BuildClient
	Deployments               *DeploymentClient
	Domains                   *DomainClient
	Droplets                  *DropletClient
	EnvVarGroups              *EnvVarGroupClient
	FeatureFlags              *FeatureFlagClient
	IsolationSegments         *IsolationSegmentClient
	Jobs                      *JobClient
	Manifests                 *ManifestClient
	Packages                  *PackageClient
	Roles                     *RoleClient
	Routes                    *RouteClient
	SecurityGroups            *SecurityGroupClient
	ServiceCredentialBindings *ServiceCredentialBindingClient
	ServiceInstances          *ServiceInstanceClient
	Spaces                    *SpaceClient
	Stacks                    *StackClient
	Users                     *UserClient
}

type commonClient struct {
	client *Client
}

type Endpoint struct {
	DopplerEndpoint   string `json:"doppler_logging_endpoint"`
	LoggingEndpoint   string `json:"logging_endpoint"`
	AuthEndpoint      string `json:"authorization_endpoint"`
	TokenEndpoint     string `json:"token_endpoint"`
	AppSSHEndpoint    string `json:"app_ssh_endpoint"`
	AppSSHOauthClient string `json:"app_ssh_oauth_client"`
}

type LoginHint struct {
	Origin string `json:"origin"`
}

// Request is used to help build up a request
type Request struct {
	method string
	url    string
	params url.Values
	body   io.Reader
	obj    interface{}
}

var ErrPreventRedirect = errors.New("prevent-redirect")

func DefaultEndpoint() *Endpoint {
	return &Endpoint{
		DopplerEndpoint: "wss://doppler.10.244.0.34.xip.io:443",
		LoggingEndpoint: "wss://loggregator.10.244.0.34.xip.io:443",
		TokenEndpoint:   "https://uaa.10.244.0.34.xip.io",
		AuthEndpoint:    "https://login.10.244.0.34.xip.io",
	}
}

// New returns a new CF client
func New(config *Config) (*Client, error) {
	client := &Client{
		Config: *config,
	}
	err := client.refreshEndpoint()
	if err != nil {
		return nil, err
	}
	client.common.client = client
	client.Organizations = (*OrgClient)(&client.common)
	client.Applications = (*AppClient)(&client.common)
	client.AppFeatures = (*AppFeatureClient)(&client.common)
	client.AppUsageEvents = (*AppUsageClient)(&client.common)
	client.AuditEvents = (*AuditEventClient)(&client.common)
	client.Buildpacks = (*BuildpackClient)(&client.common)
	client.Builds = (*BuildClient)(&client.common)
	client.Deployments = (*DeploymentClient)(&client.common)
	client.Domains = (*DomainClient)(&client.common)
	client.Droplets = (*DropletClient)(&client.common)
	client.EnvVarGroups = (*EnvVarGroupClient)(&client.common)
	client.FeatureFlags = (*FeatureFlagClient)(&client.common)
	client.IsolationSegments = (*IsolationSegmentClient)(&client.common)
	client.Jobs = (*JobClient)(&client.common)
	client.Manifests = (*ManifestClient)(&client.common)
	client.Packages = (*PackageClient)(&client.common)
	client.Roles = (*RoleClient)(&client.common)
	client.Routes = (*RouteClient)(&client.common)
	client.SecurityGroups = (*SecurityGroupClient)(&client.common)
	client.ServiceCredentialBindings = (*ServiceCredentialBindingClient)(&client.common)
	client.ServiceInstances = (*ServiceInstanceClient)(&client.common)
	client.Spaces = (*SpaceClient)(&client.common)
	client.Stacks = (*StackClient)(&client.common)
	client.Users = (*UserClient)(&client.common)
	return client, nil
}

func getUserAuth(ctx context.Context, config Config, endpoint *Endpoint) (Config, error) {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  endpoint.AuthEndpoint + "/oauth/auth",
			TokenURL: endpoint.TokenEndpoint + "/oauth/token",
		},
	}
	if config.Origin != "" {
		loginHint := LoginHint{config.Origin}
		origin, err := json.Marshal(loginHint)
		if err != nil {
			return config, fmt.Errorf("error creating login_hint: %w", err)
		}
		val := url.Values{}
		val.Set("login_hint", string(origin))
		authConfig.Endpoint.TokenURL = fmt.Sprintf("%s?%s", authConfig.Endpoint.TokenURL, val.Encode())
	}

	token, err := authConfig.PasswordCredentialsToken(ctx, config.Username, config.Password)
	if err != nil {
		return config, fmt.Errorf("error getting token: %w", err)
	}

	config.tokenSourceDeadline = &token.Expiry
	config.tokenSource = authConfig.TokenSource(ctx, token)
	config.httpClient = oauth2.NewClient(ctx, config.tokenSource)

	return config, err
}

func getClientAuth(ctx context.Context, config Config, endpoint *Endpoint) Config {
	authConfig := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     endpoint.TokenEndpoint + "/oauth/token",
	}

	config.tokenSource = authConfig.TokenSource(ctx)
	config.httpClient = authConfig.Client(ctx)
	return config
}

// getUserTokenAuth initializes client credentials from existing bearer token.
func getUserTokenAuth(ctx context.Context, config Config, endpoint *Endpoint) Config {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  endpoint.AuthEndpoint + "/oauth/auth",
			TokenURL: endpoint.TokenEndpoint + "/oauth/token",
		},
	}

	// Token is expected to have no "bearer" prefix
	token := &oauth2.Token{
		AccessToken: config.Token,
		TokenType:   "Bearer"}

	config.tokenSource = authConfig.TokenSource(ctx, token)
	config.httpClient = oauth2.NewClient(ctx, config.tokenSource)

	return config
}

func getInfo(api string, httpClient *http.Client) (*Endpoint, error) {
	var endpoint Endpoint

	if api == "" {
		return DefaultEndpoint(), nil
	}

	resp, err := httpClient.Get(api + "/v2/info")
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	err = decodeBody(resp, &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, err
}

// NewRequest is used to create a new Request
func (c *Client) NewRequest(method, path string) *Request {
	r := &Request{
		method: method,
		url:    c.Config.ApiAddress + path,
		params: make(map[string][]string),
	}
	return r
}

// NewRequestWithBody is used to create a new request with
// arbigtrary body io.Reader.
func (c *Client) NewRequestWithBody(method, path string, body io.Reader) *Request {
	r := c.NewRequest(method, path)

	// Set request body
	r.body = body

	return r
}

// DoRequest runs a request with our client
func (c *Client) DoRequest(r *Request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// DoRequestWithoutRedirects executes the request without following redirects
func (c *Client) DoRequestWithoutRedirects(r *Request) (*http.Response, error) {
	prevCheckRedirect := c.Config.httpClient.CheckRedirect
	c.Config.httpClient.CheckRedirect = func(httpReq *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		c.Config.httpClient.CheckRedirect = prevCheckRedirect
	}()
	return c.DoRequest(r)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.Config.UserAgent)
	if req.Body != nil && req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err := c.Config.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return c.handleError(resp)
	}

	return resp, nil
}

func (c *Client) GetToken() (string, error) {
	if c.Config.tokenSourceDeadline != nil && c.Config.tokenSourceDeadline.Before(time.Now()) {
		if err := c.refreshEndpoint(); err != nil {
			return "", err
		}
	}

	token, err := c.Config.tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error getting bearer token: %w", err)
	}
	return "bearer " + token.AccessToken, nil
}

func (c *Client) GetSSHCode() (string, error) {
	authorizeUrl, err := url.Parse(c.Endpoint.TokenEndpoint)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", c.Endpoint.AppSSHOauthClient) // client_id，used by cf server

	authorizeUrl.Path = "/oauth/authorize"
	authorizeUrl.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", authorizeUrl.String(), nil)
	if err != nil {
		return "", err
	}

	token, err := c.GetToken()
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", token)
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			return ErrPreventRedirect
		},
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.Config.skipSSLValidation,
			},
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	resp, err := httpClient.Do(req)
	if err == nil {
		return "", errors.New("authorization server did not redirect with one time code")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if netErr, ok := err.(*url.Error); !ok || netErr.Err != ErrPreventRedirect {
		return "", fmt.Errorf("error requesting one time code from server: %w", err)
	}

	loc, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("error getting the redirected location: %w", err)
	}

	codes := loc.Query()["code"]
	if len(codes) != 1 {
		return "", errors.New("unable to acquire one time code from authorization response")
	}

	return codes[0], nil
}

func (c *Client) handleError(resp *http.Response) (*http.Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// Unmarshal v3 error response
	var errs resource.CloudFoundryErrors
	if err := json.Unmarshal(body, &errs); err != nil {
		return resp, CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// ensure we got an error back
	if len(errs.Errors) == 0 {
		return resp, CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// TODO handle 2+ errors
	return nil, errs.Errors[0]
}

func (c *Client) refreshEndpoint() error {
	// we want to keep the Timeout value from config.httpClient
	timeout := c.Config.httpClient.Timeout

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, c.Config.httpClient)

	endpoint, err := getInfo(c.Config.ApiAddress, oauth2.NewClient(ctx, nil))

	if err != nil {
		return fmt.Errorf("Could not get api /v2/info: %w", err)
	}

	switch {
	case c.Config.Token != "":
		c.Config = getUserTokenAuth(ctx, c.Config, endpoint)
	case c.Config.ClientID != "":
		c.Config = getClientAuth(ctx, c.Config, endpoint)
	default:
		c.Config, err = getUserAuth(ctx, c.Config, endpoint)
		if err != nil {
			return err
		}
	}
	// make sure original Timeout value will be used
	if c.Config.httpClient.Timeout != timeout {
		c.Config.httpClient.Timeout = timeout
	}

	c.Endpoint = *endpoint
	return nil
}

// toHTTP converts the request to an HTTP Request
func (r *Request) toHTTP() (*http.Request, error) {

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		b, err := encodeBody(r.obj)
		if err != nil {
			return nil, err
		}
		r.body = b
	}

	// Create the HTTP Request
	return http.NewRequest(r.method, r.url, r.body)
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
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

func (c *Client) delete(path string) error {
	req := c.NewRequest("DELETE", path)
	resp, err := c.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error deleting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// some endpoints return accepted and others return no content
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting %s, response code: %d", path, resp.StatusCode)
	}
	return nil
}

func (c *Client) get(path string, result any) error {
	req := c.NewRequest("GET", path)

	resp, err := c.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error getting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting %s, response code: %d", path, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return fmt.Errorf("error decoding %s get response JSON before '%s': %w", path, buf.String(), err)
	}
	return nil
}

func (c *Client) patch(path string, params any, result any) error {
	req := c.NewRequest("PATCH", path)
	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error updating %s: %w", path, err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding %s patch response JSON: %w", path, err)
	}
	return nil
}

func (c *Client) post(id, path string, params any, result any) error {
	req := c.NewRequest("POST", path)
	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error creating %s/%s: %w", path, id, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// This should always be created but some resources return OK on POST like starting an app
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating %s/%s, response code: %d", path, id, resp.StatusCode)
	}
	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return fmt.Errorf("error decoding %s/%s post response JSON: %w", path, id, err)
		}
	}
	return nil
}
