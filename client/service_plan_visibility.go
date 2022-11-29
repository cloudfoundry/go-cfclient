package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServicePlanVisibilityClient commonClient

// Get the specified service plan visibility
func (c *ServicePlanVisibilityClient) Get(ctx context.Context, guid string) (resource.ServicePlanVisibilityType, error) {
	var s resource.ServicePlanVisibility
	err := c.client.get(ctx, path.Format("/v3/service_plans/%s/visibility", guid), &s)
	if err != nil {
		return resource.ServicePlanVisibilityNone, err
	}
	return resource.ParseServicePlanVisibilityType(s.Type)
}

// Update a service plan visibility. It behaves similar to Apply service plan visibility endpoint
// but this endpoint will replace the existing list of organizations when the service plan is
// organization visible
func (c *ServicePlanVisibilityClient) Update(ctx context.Context, guid string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error) {
	var res resource.ServicePlanVisibility
	_, err := c.client.patch(ctx, path.Format("/v3/service_plans/%s/visibility", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Apply a service plan visibility. It behaves similar to the Update service plan visibility endpoint
// but this endpoint will append to the existing list of organizations when the service plan is
// organization visible
func (c *ServicePlanVisibilityClient) Apply(ctx context.Context, guid string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error) {
	var res resource.ServicePlanVisibility
	_, err := c.client.post(ctx, path.Format("/v3/service_plans/%s/visibility", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete an organization from a service plan visibility list of organizations
// It is only defined for service plans which are organization restricted
func (c *ServicePlanVisibilityClient) Delete(ctx context.Context, guid, organizationGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/service_plans/%s/visibility/%s", guid, organizationGUID))
	return err
}
