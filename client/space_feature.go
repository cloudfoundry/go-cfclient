package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SpaceFeatureClient commonClient

// EnableSSH toggles the SSH feature for a space
func (c *SpaceFeatureClient) EnableSSH(spaceGUID string, enable bool) error {
	r := resource.SpaceFeatureUpdate{
		Enabled: enable,
	}
	_, err := c.client.patch(path("/v3/spaces/%s/features/ssh", spaceGUID), r, nil)
	return err
}
