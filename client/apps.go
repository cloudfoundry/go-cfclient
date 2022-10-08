package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
)

type AppClient commonClient

const AppsPath = "/v3/apps"

const (
	StacksField            = "stacks"
	SpaceGUIDsField        = "space_guids"
	OrganizationGUIDsField = "organization_guids"
	NamesField             = "names"
	LifecycleTypeField     = "lifecycle_type"
)

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

func (l LifecycleType) ToQueryString() url.Values {
	v := url.Values{}
	if l != LifecycleNone {
		v.Set(LifecycleTypeField, l.String())
	}
	return v
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
		v.Set(IncludeField, a.String())
	}
	return v
}

type AppListOptions struct {
	*ListOptions

	GUIDs             Filter
	Names             Filter
	OrganizationGUIDs Filter
	SpaceGUIDs        Filter
	Stacks            Filter
	LifecycleType     LifecycleType
	Include           AppIncludeType
}

func NewAppListOptions() *AppListOptions {
	return &AppListOptions{
		ListOptions: NewListOptions(),
	}
}

func (a AppListOptions) ToQuerystring() url.Values {
	v := a.ListOptions.ToQueryString()
	v = appendQueryStrings(v, a.Stacks.ToQueryString(StacksField))
	v = appendQueryStrings(v, a.SpaceGUIDs.ToQueryString(SpaceGUIDsField))
	v = appendQueryStrings(v, a.OrganizationGUIDs.ToQueryString(OrganizationGUIDsField))
	v = appendQueryStrings(v, a.GUIDs.ToQueryString(GUIDsField))
	v = appendQueryStrings(v, a.Names.ToQueryString(NamesField))
	v = appendQueryStrings(v, a.LifecycleType.ToQueryString())
	v = appendQueryStrings(v, a.Include.ToQueryString())
	return v
}

func (c *AppClient) Create(r resource.CreateAppRequest) (*resource.App, error) {
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"space": resource.ToOneRelationship{
				Data: resource.Relationship{
					GUID: r.SpaceGUID,
				},
			},
		},
	}
	if len(r.EnvironmentVariables) > 0 {
		params["environment_variables"] = r.EnvironmentVariables
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	var app resource.App
	err := c.client.post(r.Name, AppsPath, params, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *AppClient) Delete(guid string) error {
	return c.client.delete(joinPath(AppsPath, guid))
}

func (c *AppClient) Get(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.get(joinPath(AppsPath, guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *AppClient) GetEnvironment(appGUID string) (*resource.AppEnvironment, error) {
	var appEnv resource.AppEnvironment
	err := c.client.get(joinPath(AppsPath, appGUID, "env"), &appEnv)
	if err != nil {
		return nil, err
	}
	return &appEnv, nil
}

func (c *AppClient) GetInclude(guid string, include AppIncludeType) (*resource.App, error) {
	var app resource.App
	err := c.client.get(joinPathAndQS(include.ToQueryString(), AppsPath, guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *AppClient) List(opts *AppListOptions) ([]*resource.App, *Pager, error) {
	var res resource.ListAppsResponse
	err := c.client.get(joinPathAndQS(opts.ToQuerystring(), AppsPath), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := &Pager{
		pagination: res.Pagination,
	}
	return res.Resources, pager, nil
}

func (c *AppClient) ListAll() ([]*resource.App, error) {
	opts := NewAppListOptions()
	var allApps []*resource.App
	for {
		apps, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allApps = append(allApps, apps...)
		if !pager.NextPage(opts.ListOptions) {
			break
		}
	}
	return allApps, nil
}

func (c *AppClient) SetEnvVariables(appGUID string, envRequest resource.EnvVar) (*resource.EnvVar, error) {
	var envVarResponse resource.EnvVarResponse
	err := c.client.patch(joinPath(AppsPath, appGUID, "environment_variables"), envRequest, &envVarResponse)
	if err != nil {
		return nil, err
	}
	return &envVarResponse.EnvVar, nil
}

func (c *AppClient) Start(guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.post(guid, joinPath(AppsPath, guid, "actions/start"), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *AppClient) Update(appGUID string, r resource.UpdateAppRequest) (*resource.App, error) {
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	var app resource.App
	err := c.client.patch("/v3/apps/"+appGUID, params, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
