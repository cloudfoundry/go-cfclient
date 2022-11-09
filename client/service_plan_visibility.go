package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServicePlanVisibilityClient commonClient

// Get the specified service plan visibility
func (c *ServicePlanVisibilityClient) Get(guid string) (resource.ServicePlanVisibilityType, error) {
	var s resource.ServicePlanVisibility
	err := c.client.get(path("/v3/service_plans/%s/visibility", guid), &s)
	if err != nil {
		return resource.ServicePlanVisibilityNone, err
	}
	return resource.ParseServicePlanVisibilityType(s.Type)
}

// Update a service plan visibility. It behaves similar to Apply service plan visibility endpoint
// but this endpoint will replace the existing list of organizations when the service plan is
// organization visible
func (c *ServicePlanVisibilityClient) Update(guid string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error) {
	var res resource.ServicePlanVisibility
	_, err := c.client.patch(path("/v3/service_plans/%s/visibility", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Apply a service plan visibility. It behaves similar to the Update service plan visibility endpoint
// but this endpoint will append to the existing list of organizations when the service plan is
// organization visible
func (c *ServicePlanVisibilityClient) Apply(guid string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error) {
	var res resource.ServicePlanVisibility
	_, err := c.client.post(path("/v3/service_plans/%s/visibility", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete an organization from a service plan visibility list of organizations
// It is only defined for service plans which are org-restricted
func (c *ServicePlanVisibilityClient) Delete(guid, orgGUID string) error {
	_, err := c.client.delete(path("/v3/service_plans/%s/visibility/%s", guid, orgGUID))
	return err
}
