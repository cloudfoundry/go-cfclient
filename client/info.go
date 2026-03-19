package client

import (
	"context"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type InfoClient commonClient

// Get retrieves information about the Cloud Foundry deployment
//
// This endpoint returns metadata about the Cloud Foundry deployment including
// version, build info, CLI version requirements, and operator-configured custom metadata.
//
// Authentication: No authentication required
func (c *InfoClient) Get(ctx context.Context) (*resource.Info, error) {
	var info resource.Info
	err := c.client.get(ctx, "/v3/info", &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetUsageSummary retrieves platform-wide usage statistics
//
// This endpoint returns usage information across the entire Cloud Foundry deployment
// including started instances, memory usage, routes, service instances, and more.
//
// Authentication: Requires authentication with cloud_controller.admin or cloud_controller.admin_read_only scope
func (c *InfoClient) GetUsageSummary(ctx context.Context) (*resource.InfoUsageSummary, error) {
	var usageSummary resource.InfoUsageSummary
	err := c.client.get(ctx, "/v3/info/usage_summary", &usageSummary)
	if err != nil {
		return nil, err
	}
	return &usageSummary, nil
}
