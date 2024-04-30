package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/cloudfoundry/go-cfclient/v3/internal/check"
	internal "github.com/cloudfoundry/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry/go-cfclient/v3/internal/ios"
	"github.com/cloudfoundry/go-cfclient/v3/internal/path"
)

// Client used to communicate with Cloud Foundry
type Client struct {
	Admin                     *AdminClient
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
	Organizations             *OrganizationClient
	OrganizationQuotas        *OrganizationQuotaClient
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
	*config.Config
}

type commonClient struct {
	client *Client
}

// New returns a new CF client
func New(config *config.Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	client := &Client{
		Config: config,
	}

	// populate sub-clients
	client.common.client = client
	client.Admin = (*AdminClient)(&client.common)
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
	client.Organizations = (*OrganizationClient)(&client.common)
	client.OrganizationQuotas = (*OrganizationQuotaClient)(&client.common)
	client.Packages = (*PackageClient)(&client.common)
	client.Processes = (*ProcessClient)(&client.common)
	client.Revisions = (*RevisionClient)(&client.common)
	client.ResourceMatches = (*ResourceMatchClient)(&client.common)
	client.Roles = (*RoleClient)(&client.common)
	client.Root = (*RootClient)(&client.common)
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

// ExecuteAuthRequest executes an HTTP request with authentication.
func (c *Client) ExecuteAuthRequest(req *http.Request) (*http.Response, error) {
	return c.executeHTTPRequest(req, true)
}

// ExecuteRequest executes an HTTP request without authentication.
func (c *Client) ExecuteRequest(req *http.Request) (*http.Response, error) {
	return c.executeHTTPRequest(req, false)
}

// SSHCode generates an SSH code that can be used by generic SSH clients to SSH into app instances
func (c *Client) SSHCode(ctx context.Context) (string, error) {
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", c.SSHOAuthClientID())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.AuthURL(path.Format("/oauth/authorize?%s", values)), nil)
	if err != nil {
		return "", fmt.Errorf("error creating SSH code request: %w", err)
	}

	resp, err := c.ExecuteAuthRequest(internal.IgnoreRedirect(req))
	if err != nil {
		return "", fmt.Errorf("error executing SSH code request: %w", err)
	}

	defer ios.Close(resp.Body)

	if !internal.IsResponseRedirect(resp.StatusCode) {
		return "", fmt.Errorf(
			"expected UAA to return a 302 location that contains the code, but instead got a %d", resp.StatusCode)
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

// delete does an HTTP DELETE to the specified endpoint and returns the job ID if any.
//
// This function takes the relative API resource path. If the resource returns an async job ID
// then the function returns the job GUID which the caller can reference via the job endpoint.
func (c *Client) delete(ctx context.Context, resourcePath string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.ApiURL(resourcePath), nil)
	if err != nil {
		return "", fmt.Errorf("creating DELETE request for %s failed: %w", resourcePath, err)
	}
	resp, err := c.ExecuteAuthRequest(req)
	if err != nil {
		return "", fmt.Errorf("executing DELETE request for %s failed: %w", resourcePath, err)
	}
	defer ios.Close(resp.Body)
	return internal.DecodeJobIDOrBody(resp, nil)
}

// get does an HTTP GET to the specified endpoint and automatically handles unmarshalling
// the result JSON body
func (c *Client) get(ctx context.Context, resourcePath string, result any) error {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return errors.New("expected result to be nil or a pointer type")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ApiURL(resourcePath), nil)
	if err != nil {
		return fmt.Errorf("error creating GET request for %s: %w", resourcePath, err)
	}

	resp, err := c.ExecuteAuthRequest(req)
	if err != nil {
		return fmt.Errorf("error executing GET request for %s: %w", resourcePath, err)
	}
	defer ios.Close(resp.Body)

	return internal.DecodeBody(resp, result)
}

// list does an HTTP GET to the specified endpoint and automatically handles unmarshalling the result JSON body.
// This is a utility function to support list functions.
func (c *Client) list(ctx context.Context, urlPathFormat string, queryStrFunc func() (url.Values, error), result any) error {
	params, err := queryStrFunc()
	if err != nil {
		return fmt.Errorf("error while generate query params: %w", err)
	}
	if len(params) > 0 {
		urlPathFormat = strings.TrimSuffix(urlPathFormat+"?"+params.Encode(), "?")
	}
	return c.get(ctx, urlPathFormat, result)
}

// patch does an HTTP PATCH to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to PATCH and an optional
// struct to unmarshall the result body. If the resource returns an async job ID instead of a
// response body, then the body won't be unmarshalled and the function returns the job GUID
// which the caller can reference via the job endpoint.
func (c *Client) patch(ctx context.Context, resourcePath string, params any, result any) (string, error) {
	return c.createOrUpdate(ctx, http.MethodPatch, resourcePath, params, result)
}

// post does an HTTP POST to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST and an optional
// struct to unmarshall the result body. If the resource returns an async job ID in the Location
// header then the job GUID is returned which the caller can reference via the job endpoint.
func (c *Client) post(ctx context.Context, resourcePath string, params, result any) (string, error) {
	return c.createOrUpdate(ctx, http.MethodPost, resourcePath, params, result)
}

// Download the bits of an existing package or droplet
// It is the caller's responsibility to close the io.ReadCloser
func (c *Client) download(ctx context.Context, resourcePath string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ApiURL(resourcePath), nil)
	if err != nil {
		return nil, fmt.Errorf("creating download request for %s failed: %w", resourcePath, err)
	}
	resp, err := c.ExecuteAuthRequest(internal.IgnoreRedirect(req))
	if err != nil {
		return nil, fmt.Errorf("executing download request for %s failed: %w", resourcePath, err)
	}
	ios.Close(resp.Body)
	if !internal.IsResponseRedirect(resp.StatusCode) {
		return nil, fmt.Errorf("error downloading `%s` bits, expected redirect to blobstore", resourcePath)
	}
	// get the full URL to the blobstore via the Location header
	blobStoreLocation := resp.Header.Get("Location")
	if blobStoreLocation == "" {
		return nil, errors.New("response redirect Location header was empty")
	}
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, blobStoreLocation, nil)
	if err != nil {
		return nil, fmt.Errorf("creating blob download request for %s failed: %w", blobStoreLocation, err)
	}
	resp, err = c.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("executing blob download request for %s failed: %w", blobStoreLocation, err)
	}
	return resp.Body, nil
}

// postFileUpload does an HTTP POST to the specified endpoint and automatically handles uploading the specified file
// and handling the result whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST and an optional
// struct to unmarshall the result body. If the resource returns an async job ID in the Location
// header then the job GUID is returned which the caller can reference via the job endpoint.
func (c *Client) postFileUpload(ctx context.Context, path, fieldName, fileName string, fileContent io.Reader, result any) (string, error) {
	// Validate input parameters
	if path == "" || fieldName == "" || fileName == "" {
		return "", errors.New("path, fieldName, and fileName are required")
	}
	if fileContent == nil {
		return "", fmt.Errorf("no content was provided for the %s file", fileName)
	}

	if !check.IsNil(result) && !check.IsPointer(result) {
		return "", errors.New("expected result to be a pointer type or nil")
	}
	//
	// Prepare multipart form data
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)
	part, err := formWriter.CreateFormFile(fieldName, filepath.Base(fileName))
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s: %w", path, err)
	}
	if _, err = io.Copy(part, fileContent); err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed on copy: %w", path, err)
	}
	if err = formWriter.Close(); err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed to close multipart form writer: %w", path, err)
	}

	// Create and execute the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ApiURL(path), body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", formWriter.FormDataContentType())

	resp, err := c.ExecuteAuthRequest(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer ios.Close(resp.Body) // Ensure closure of the response body

	// Decode the response
	return internal.DecodeJobIDAndBody(resp, result)
}

// createOrUpdate is a utility function for patch and post that does an HTTP POST or PATCH to the specified
// endpoint and automatically handles the result whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST/PATCH and an optional
// struct to unmarshall the result body. If the resource returns an async job ID in the Location
// header then the job GUID is returned which the caller can reference via the job endpoint.
func (c *Client) createOrUpdate(ctx context.Context, method, resourcePath string, params, result any) (string, error) {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return "", errors.New("expected result to be a pointer type, or nil")
	}
	body, err := internal.EncodeBody(params)
	if err != nil {
		return "", fmt.Errorf("failed to encode params: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.ApiURL(resourcePath), body)
	if err != nil {
		return "", fmt.Errorf("creating %s request for %s failed: %w", method, resourcePath, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.ExecuteAuthRequest(req)
	if err != nil {
		return "", fmt.Errorf("executing %s request for %s failed: %w", method, resourcePath, err)
	}
	defer ios.Close(resp.Body)
	return internal.DecodeJobIDOrBody(resp, result)
}

// executeHTTPRequest is the low level client function that handles executing the request against the
// correct http.Client.
func (c *Client) executeHTTPRequest(req *http.Request, includeAuthHeader bool) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", c.UserAgent())
	if includeAuthHeader {
		resp, err = c.HTTPAuthClient().Do(req)
	} else {
		resp, err = c.HTTPClient().Do(req)
	}

	if err != nil {
		return nil, fmt.Errorf("error executing request, failed during HTTP request send: %w", err)
	}
	if !internal.IsStatusSuccess(resp.StatusCode) {
		return nil, internal.DecodeError(resp)
	}

	return resp, err
}
