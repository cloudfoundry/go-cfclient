package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type V3ProcessReference struct {
	GUID string `json:"guid"`
	Type string `type:"type"`
}

type V3DeploymentStatus struct {
	Value   string            `json:"value"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details"`
}

type V3Deployment struct {
	GUID            string               `json:"guid"`
	State           string               `json:"state"`
	Status          V3DeploymentStatus   `json:"status"`
	Strategy        string               `json:"strategy"`
	Droplet         V3Relationship       `json:"droplet"`
	PreviousDroplet V3Relationship       `json:"previous_droplet"`
	NewProcesses    []V3ProcessReference `json:"new_processes"`
	Revision        struct {
		GUID    string `json:"guid"`
		Version int    `json:"version"`
	} `json:"revision"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
	Metadata      V3Metadata                     `json:"metadata,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
}

func (c *Client) GetV3Deployment(deploymentGUID string) (*V3Deployment, error) {
	req := c.NewRequest("GET", "/v3/deployments/"+deploymentGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting deployment")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting deployment with GUID [%s], response code: %d", deploymentGUID, resp.StatusCode)
	}

	var r V3Deployment
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error reading deployment response JSON")
	}

	return &r, nil
}
