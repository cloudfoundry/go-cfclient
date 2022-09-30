package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) SetCurrentDropletForApp(appGUID, dropletGUID string) (*v3.CurrentDropletResponse, error) {
	req := c.NewRequest("PATCH", "/v3/apps/"+appGUID+"/relationships/current_droplet")
	req.obj = v3.ToOneRelationship{Data: v3.Relationship{GUID: dropletGUID}}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error setting droplet for v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error setting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r v3.CurrentDropletResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading droplet response JSON")
	}

	return &r, nil
}

func (c *Client) GetCurrentDropletForApp(appGUID string) (*v3.Droplet, error) {
	req := c.NewRequest("GET", "/v3/apps/"+appGUID+"/droplets/current")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting droplet for v3 app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting droplet for v3 app with GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r v3.Droplet
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading droplet response JSON")
	}

	return &r, nil
}

func (c *Client) DeleteDroplet(dropletGUID string) error {
	req := c.NewRequest("DELETE", "/v3/droplets/"+dropletGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrapf(err, "Error deleting droplet %s", dropletGUID)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting droplet %s with response code %d", dropletGUID, resp.StatusCode)
	}

	return nil
}
