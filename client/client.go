package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"io"
	http2 "net/http"
	"net/url"
	"strings"
)

var ErrPreventRedirect = errors.New("prevent-redirect")

// Client used to communicate with Cloud Foundry
type Client struct {
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
	Organizations             *OrgClient
	OrganizationQuotas        *OrgQuotaClient
	Packages                  *PackageClient
	Processes                 *ProcessClient
	Revisions                 *RevisionClient
	ResourceMatches           *ResourceMatchClient
	Roles                     *RoleClient
	Root                      *RootClient
	Routes                    *RouteClient
	SecurityGroups            *SecurityGroupClient
	ServiceBrokers            *ServiceBrokerClient
	ServiceCredentialBindings *ServiceCredentialBindingClient
	ServiceInstances          *ServiceInstanceClient
	ServiceOfferings          *ServiceOfferingClient
	ServicePlans              *ServicePlanClient
	ServicePlansVisibility    *ServicePlanVisibilityClient
	ServiceRouteBindings      *ServiceRouteBindingClient
	ServiceUsageEvents        *ServiceUsageClient
	Sidecars                  *SidecarClient
	Spaces                    *SpaceClient
	SpaceFeatures             *SpaceFeatureClient
	SpaceQuotas               *SpaceQuotaClient
	Stacks                    *StackClient
	Tasks                     *TaskClient
	Users                     *UserClient

	common commonClient // Reuse a single struct instead of allocating one for each commonClient on the heap.
	config *config.Config

	authenticatedHTTPExecutor   *http.Executor
	authenticatedClientProvider *http.OAuthSessionManager
}

type commonClient struct {
	client *Client
}

// New returns a new CF client
func New(config *config.Config) (*Client, error) {
	// construct an unauthenticated root client
	unauthenticatedClientProvider := http.NewUnauthenticatedClientProvider(config.BaseHTTPClient)
	unauthenticatedHTTPExecutor := http.NewExecutor(unauthenticatedClientProvider, config.APIEndpointURL, config.UserAgent)
	rootClient := NewRootClient(unauthenticatedHTTPExecutor)
	err := authServiceDiscovery(config, rootClient)
	if err != nil {
		return nil, err
	}

	// create the client instance
	authenticatedClientProvider := http.NewOAuthSessionManager(config)
	authenticatedHTTPExecutor := http.NewExecutor(authenticatedClientProvider, config.APIEndpointURL, config.UserAgent)
	client := &Client{
		config:                      config,
		authenticatedHTTPExecutor:   authenticatedHTTPExecutor,
		authenticatedClientProvider: authenticatedClientProvider,
	}

	// populate sub-clients
	client.common.client = client
	client.Root = rootClient
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
	client.Organizations = (*OrgClient)(&client.common)
	client.OrganizationQuotas = (*OrgQuotaClient)(&client.common)
	client.Packages = (*PackageClient)(&client.common)
	client.Processes = (*ProcessClient)(&client.common)
	client.Revisions = (*RevisionClient)(&client.common)
	client.ResourceMatches = (*ResourceMatchClient)(&client.common)
	client.Roles = (*RoleClient)(&client.common)
	client.Routes = (*RouteClient)(&client.common)
	client.SecurityGroups = (*SecurityGroupClient)(&client.common)
	client.ServiceBrokers = (*ServiceBrokerClient)(&client.common)
	client.ServiceCredentialBindings = (*ServiceCredentialBindingClient)(&client.common)
	client.ServiceInstances = (*ServiceInstanceClient)(&client.common)
	client.ServiceOfferings = (*ServiceOfferingClient)(&client.common)
	client.ServicePlans = (*ServicePlanClient)(&client.common)
	client.ServicePlansVisibility = (*ServicePlanVisibilityClient)(&client.common)
	client.ServiceRouteBindings = (*ServiceRouteBindingClient)(&client.common)
	client.ServiceUsageEvents = (*ServiceUsageClient)(&client.common)
	client.Sidecars = (*SidecarClient)(&client.common)
	client.Spaces = (*SpaceClient)(&client.common)
	client.SpaceQuotas = (*SpaceQuotaClient)(&client.common)
	client.SpaceFeatures = (*SpaceFeatureClient)(&client.common)
	client.Stacks = (*StackClient)(&client.common)
	client.Tasks = (*TaskClient)(&client.common)
	client.Users = (*UserClient)(&client.common)
	return client, nil
}

// SSHCode generates an SSH code that can be used by generic SSH clients to SSH into app instances
func (c *Client) SSHCode() (string, error) {
	// need this to grab the SSH client id, should probably be cached in config
	r, err := c.Root.Get()
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", r.Links.AppSSH.Meta.OauthClient) // client_id，used by cf server

	token, err := c.authenticatedClientProvider.AccessToken()
	if err != nil {
		return "", err
	}

	req := http.NewRequest("GET", path.Format("/oauth/authorize", values)).
		WithHeader("authorization", token)

	nonRedirectingHTTPClient := &http2.Client{
		CheckRedirect: func(req *http2.Request, _ []*http2.Request) error {
			return ErrPreventRedirect
		},
		Timeout:   c.config.BaseHTTPClient.Timeout,
		Transport: c.config.BaseHTTPClient.Transport,
	}

	unauthenticatedClientProvider := http.NewUnauthenticatedClientProvider(nonRedirectingHTTPClient)
	unauthenticatedHTTPExecutor := http.NewExecutor(unauthenticatedClientProvider, c.config.UAAEndpointURL, c.config.UserAgent)
	resp, err := unauthenticatedHTTPExecutor.ExecuteRequest(req)
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

// delete does an HTTP DELETE to the specified endpoint and returns the job ID if any
//
// This function takes the relative API resource path. If the resource returns an async job ID
// then the function returns the job GUID which the caller can reference via the job endpoint.
func (c *Client) delete(path string) (string, error) {
	req := http.NewRequest("DELETE", path)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error deleting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// some endpoints return accepted and others return no content
	if resp.StatusCode != http2.StatusAccepted && resp.StatusCode != http2.StatusNoContent {
		return "", c.handleError(resp)
	}
	return c.decodeBodyOrJobID(resp, nil)
}

// get does an HTTP GET to the specified endpoint and automatically handles unmarshalling
// the result JSON body
func (c *Client) get(path string, result any) error {
	req := http.NewRequest("GET", path)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return fmt.Errorf("error getting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK {
		return c.handleError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return fmt.Errorf("error decoding %s get response JSON before '%s': %w", path, buf.String(), err)
	}
	return nil
}

// patch does an HTTP PATCH to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to PATCH and an optional
// struct to unmarshall the result body. If the resource returns an async job ID instead of a
// response body, then the body won't be unmarshalled and the function returns the job GUID
// which the caller can reference via the job endpoint.
func (c *Client) patch(path string, params any, result any) (string, error) {
	req := http.NewRequest("PATCH", path).WithObject(params)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error updating %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK && resp.StatusCode != http2.StatusAccepted {
		return "", c.handleError(resp)
	}
	return c.decodeBodyOrJobID(resp, &result)
}

// post does an HTTP POST to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST and an optional
// struct to unmarshall the result body. If the resource returns an async job ID instead of a
// response body, then the body won't be unmarshalled and the function returns the job GUID
// which the caller can reference via the job endpoint.
func (c *Client) post(path string, params, result any) (string, error) {
	req := http.NewRequest("POST", path).WithObject(params)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error creating %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// Endpoints return different status codes for posts
	if resp.StatusCode != http2.StatusCreated && resp.StatusCode != http2.StatusOK && resp.StatusCode != http2.StatusAccepted {
		return "", c.handleError(resp)
	}
	return c.decodeBodyOrJobID(resp, result)
}

// handleError attempts to unmarshall the response body as a CF error
func (c *Client) handleError(resp *http2.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CloudFoundryHTTPError{
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
		return CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// ensure we got an error back
	if len(errs.Errors) == 0 {
		return CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// TODO handle 2+ errors
	return errs.Errors[0]
}

// decodeBodyOrJobID returns the jobID if specified in the Location response header, otherwise it
// unmarshalls the JSON response to result
func (c *Client) decodeBodyOrJobID(resp *http2.Response, result any) (string, error) {
	var jobID string
	location, err := resp.Location()
	if err == nil && strings.Contains(location.Path, "jobs") {
		p := strings.Split(location.Path, "/")
		jobID = p[len(p)-1]
	} else if result != nil {
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return "", fmt.Errorf("error decoding response JSON: %w", err)
		}
	}
	return jobID, nil
}

// authServiceDiscovery sets the UAA and Login endpoint if the user didn't configure these manually
func authServiceDiscovery(config *config.Config, rootClient *RootClient) error {
	if config.UAAEndpointURL != "" && config.LoginEndpointURL != "" {
		return nil
	}
	root, err := rootClient.Get()
	if err != nil {
		return err
	}
	config.UAAEndpointURL = root.Links.Uaa.Href
	config.LoginEndpointURL = root.Links.Login.Href
	return nil
}
