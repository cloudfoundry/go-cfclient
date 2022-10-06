package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type DropletClient commonClient

func (c *DropletClient) SetCurrentForApp(appGUID, dropletGUID string) (*resource.CurrentDropletResponse, error) {
	req := c.client.NewRequest("PATCH", "/v3/apps/"+appGUID+"/relationships/current_droplet")
	req.obj = resource.ToOneRelationship{Data: resource.Relationship{GUID: dropletGUID}}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error setting droplet for v3 app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error setting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.CurrentDropletResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error reading droplet response JSON: %w", err)
	}

	return &r, nil
}

func (c *DropletClient) GetCurrentForApp(appGUID string) (*resource.Droplet, error) {
	req := c.client.NewRequest("GET", "/v3/apps/"+appGUID+"/droplets/current")
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting droplet for v3 app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.Droplet
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error reading droplet response JSON: %w", err)
	}

	return &r, nil
}

func (c *DropletClient) Delete(dropletGUID string) error {
	req := c.client.NewRequest("DELETE", "/v3/droplets/"+dropletGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error deleting droplet %s: %w", dropletGUID, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting droplet %s with response code %d", dropletGUID, resp.StatusCode)
	}

	return nil
}
