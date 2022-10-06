package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

type DropletClient commonClient

func (c *DropletClient) SetCurrentForApp(appGUID, dropletGUID string) (*resource.CurrentDropletResponse, error) {
	req := c.client.NewRequest("PATCH", "/v3/apps/"+appGUID+"/relationships/current_droplet")
	req.obj = resource.ToOneRelationship{Data: resource.Relationship{GUID: dropletGUID}}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error setting droplet for v3 app")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error setting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.CurrentDropletResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading droplet response JSON")
	}

	return &r, nil
}

func (c *DropletClient) GetCurrentForApp(appGUID string) (*resource.Droplet, error) {
	req := c.client.NewRequest("GET", "/v3/apps/"+appGUID+"/droplets/current")
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting droplet for v3 app")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.Droplet
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading droplet response JSON")
	}

	return &r, nil
}

func (c *DropletClient) Delete(dropletGUID string) error {
	req := c.client.NewRequest("DELETE", "/v3/droplets/"+dropletGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return errors.Wrapf(err, "Error deleting droplet %s", dropletGUID)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting droplet %s with response code %d", dropletGUID, resp.StatusCode)
	}

	return nil
}
