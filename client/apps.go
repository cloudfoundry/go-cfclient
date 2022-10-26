package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type AppClient commonClient

// LifecycleType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#list-apps
type LifecycleType int

const (
	LifecycleNone LifecycleType = iota
	LifecycleBuildpack
	LifecycleDocker
)

func (l LifecycleType) String() string {
	switch l {
	case LifecycleBuildpack:
		return "buildpack"
	case LifecycleDocker:
		return "docker"
	}
	return ""
}

// AppIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type AppIncludeType int

const (
	AppIncludeNone AppIncludeType = iota
	AppIncludeSpace
	AppIncludeSpaceOrganization
)

func (a AppIncludeType) String() string {
	switch a {
	case AppIncludeSpace:
		return "space"
	case AppIncludeSpaceOrganization:
		return "space.organization"
	}
	return ""
}

func (a AppIncludeType) ToQueryString() url.Values {
	v := url.Values{}
	if a != AppIncludeNone {
		v.Set("include", a.String())
	}
	return v
}

// AppListOptions list filters
type AppListOptions struct {
	*ListOptions

	GUIDs             Filter         `filter:"guids,omitempty"`
	Names             Filter         `filter:"names,omitempty"`
	OrganizationGUIDs Filter         `filter:"organization_guids,omitempty"`
	SpaceGUIDs        Filter         `filter:"space_guids,omitempty"`
	Stacks            Filter         `filter:"stacks,omitempty"`
	LifecycleType     LifecycleType  `filter:"lifecycle_type,omitempty"`
	Include           AppIncludeType `filter:"include,omitempty"`
}

// NewAppListOptions creates new options to pass to list
func NewAppListOptions() *AppListOptions {
	return &AppListOptions{
		ListOptions: NewListOptions(),
	}
}

// Create a new app
func (c *AppClient) Create(r *resource.AppCreate) (*resource.App, error) {
	var app resource.App
	err := c.client.post(r.Name, "/v3/apps", r, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Delete the specified app
func (c *AppClient) Delete(guid string) error {
	return c.client.delete(path("/v3/apps/%s", guid))
}

// Get the specified app
func (c *AppClient) Get(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.get(path("/v3/apps/%s", guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetEnvironment retrieves the environment variables that will be provided to an app at runtime.
// It will include environment variables for Environment Variable Groups and Service Bindings.
func (c *AppClient) GetEnvironment(appGUID string) (*resource.AppEnvironment, error) {
	var appEnv resource.AppEnvironment
	err := c.client.get(path("/v3/apps/%s/env", appGUID), &appEnv)
	if err != nil {
		return nil, err
	}
	return &appEnv, nil
}

// GetAndInclude allows callers to fetch an app and include information of parent objects in the response
func (c *AppClient) GetAndInclude(guid string, include AppIncludeType) (*resource.App, error) {
	var app resource.App

	err := c.client.get(path("/v3/apps/%s?%s", guid, include.ToQueryString()), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List all apps the user has access to in paged results
func (c *AppClient) List(opts *AppListOptions) ([]*resource.App, *Pager, error) {
	var res resource.AppList
	err := c.client.get(path("/v3/apps?%s", opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all apps the user has access to
func (c *AppClient) ListAll() ([]*resource.App, error) {
	opts := NewAppListOptions()
	var allApps []*resource.App
	for {
		apps, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allApps = append(allApps, apps...)
		if !pager.HasNextPage() {
			break
		}
		opts.ListOptions = pager.NextPage(opts.ListOptions)
	}
	return allApps, nil
}

// Permissions gets the current user’s permissions for the given app.
// If a user can see an app, then they can see its basic data.
// Only admin, read-only admins, and space developers can read sensitive data.
func (c *AppClient) Permissions(guid string) (*resource.AppPermissions, error) {
	var appPerms resource.AppPermissions
	err := c.client.get(path("/v3/apps/%s/permissions", guid), &appPerms)
	if err != nil {
		return nil, err
	}
	return &appPerms, nil
}

// Restart will synchronously stop and start an application.
// Unlike the start and stop actions, this endpoint will error if the app is not successfully stopped in the runtime.
// For restarting applications without downtime, see the Deployments resource.
func (c *AppClient) Restart(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.post(guid, path("/v3/apps/%s/actions/restart", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// SetEnvVariables updates the environment variables associated with the given app.
// The variables given in the request will be merged with the existing app environment variables.
// Any requested variables with a value of null will be removed from the app.
//
// Environment variable names may not start with VCAP_
// PORT is not a valid environment variable.
func (c *AppClient) SetEnvVariables(appGUID string, envRequest resource.EnvVar) (*resource.EnvVar, error) {
	var envVarResponse resource.EnvVarResponse
	err := c.client.patch(path("/v3/apps/%s/environment_variables", appGUID), envRequest, &envVarResponse)
	if err != nil {
		return nil, err
	}
	return &envVarResponse.EnvVar, nil
}

// Start the app if not already started
func (c *AppClient) Start(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.post(guid, path("/v3/apps/%s/actions/start", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Stop the app if not already stopped
func (c *AppClient) Stop(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.post(guid, path("/v3/apps/%s/actions/stop", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Update the specified attributes of the app
func (c *AppClient) Update(guid string, r *resource.AppUpdate) (*resource.App, error) {
	var app resource.App
	err := c.client.patch(path("/v3/apps/%s", guid), r, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// SSHEnabled returns if an application’s runtime environment will accept ssh connections.
// If ssh is disabled, the reason field will describe whether it is disabled globally,
// at the space level, or at the app level.
func (c *AppClient) SSHEnabled(guid string) (*resource.AppSSHEnabled, error) {
	var appSSH resource.AppSSHEnabled
	err := c.client.get(path("/v3/apps/%s/ssh_enabled", guid), &appSSH)
	if err != nil {
		return nil, err
	}
	return &appSSH, nil
}
