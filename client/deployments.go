package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

type DeploymentClient commonClient

func (c *DeploymentClient) Get(deploymentGUID string) (*resource.Deployment, error) {
	req := c.client.NewRequest("GET", "/v3/deployments/"+deploymentGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting deployment")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting deployment with GUID [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	var r resource.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *DeploymentClient) Create(appGUID string, optionalParams *resource.CreateDeploymentOptionalParameters) (*resource.Deployment, error) {
	// validate the params
	if optionalParams != nil {
		if optionalParams.Droplet != nil && optionalParams.Revision != nil {
			return nil, errors.New("droplet and revision cannot both be set")
		}
	}

	requestBody := resource.CreateDeploymentRequest{}
	requestBody.CreateDeploymentOptionalParameters = optionalParams

	requestBody.Relationships = struct {
		App resource.ToOneRelationship "json:\"app\""
	}{
		App: resource.ToOneRelationship{
			Data: resource.Relationship{
				GUID: appGUID,
			},
		},
	}

	req := c.client.NewRequest("POST", "/v3/deployments")
	req.obj = requestBody

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating deployment")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating deployment for app GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.Deployment
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}

func (c *DeploymentClient) Cancel(deploymentGUID string) error {
	req := c.client.NewRequest("POST", "/v3/deployments/"+deploymentGUID+"/actions/cancel")
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error canceling deployment")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error canceling deployment [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	return nil
}
