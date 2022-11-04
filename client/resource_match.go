package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type ResourceMatchClient commonClient

// Create a list of cached resources from the input list
func (c *ResourceMatchClient) Create(toMatch *resource.ResourceMatches) (*resource.ResourceMatches, error) {
	var matched resource.ResourceMatches
	err := c.client.post("", "/v3/resource_matches", toMatch, &matched)
	if err != nil {
		return nil, err
	}
	return &matched, nil
}
