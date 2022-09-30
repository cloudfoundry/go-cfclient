package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) GetDeployment(deploymentGUID string) (*v3.Deployment, error) {
	req := c.NewRequest("GET", "/v3/deployments/"+deploymentGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting deployment with GUID [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	var r v3.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *Client) CreateDeployment(appGUID string, optionalParams *v3.CreateDeploymentOptionalParameters) (*v3.Deployment, error) {
	// validate the params
	if optionalParams != nil {
		if optionalParams.Droplet != nil && optionalParams.Revision != nil {
			return nil, errors.New("droplet and revision cannot both be set")
		}
	}

	requestBody := v3.CreateDeploymentRequest{}
	requestBody.CreateDeploymentOptionalParameters = optionalParams

	requestBody.Relationships = struct {
		App v3.ToOneRelationship "json:\"app\""
	}{
		App: v3.ToOneRelationship{
			Data: v3.Relationship{
				GUID: appGUID,
			},
		},
	}

	req := c.NewRequest("POST", "/v3/deployments")
	req.obj = requestBody

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating deployment for app GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r v3.Deployment
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *Client) CancelDeployment(deploymentGUID string) error {
	req := c.NewRequest("POST", "/v3/deployments/"+deploymentGUID+"/actions/cancel")
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error canceling deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error canceling deployment [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	return nil
}
