package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type DeploymentClient commonClient

func (c *DeploymentClient) Get(deploymentGUID string) (*resource.Deployment, error) {
	req := c.client.NewRequest("GET", "/v3/deployments/"+deploymentGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting deployment: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting deployment with GUID [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	var r resource.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error reading deployment response JSON: %w", err)
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
		return nil, fmt.Errorf("error creating deployment: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating deployment for app GUID [%s], response code: %d", appGUID, resp.StatusCode)
	}

	var r resource.Deployment
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error reading deployment response JSON: %w", err)
	}

	return &r, nil
}

func (c *DeploymentClient) Cancel(deploymentGUID string) error {
	req := c.client.NewRequest("POST", "/v3/deployments/"+deploymentGUID+"/actions/cancel")
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error canceling deployment: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error canceling deployment [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	return nil
}
