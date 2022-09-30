package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) GetBuildByGUID(buildGUID string) (*v3.Build, error) {
	resp, err := c.DoRequest(c.NewRequest("GET", "/v3/builds/"+buildGUID))
	if err != nil {
		return nil, errors.Wrap(err, "Error getting  build")
	}
	defer resp.Body.Close()

	var build v3.Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "Error reading  build JSON")
	}

	return &build, nil
}

func (c *Client) CreateBuild(packageGUID string, lifecycle *v3.Lifecycle, metadata *v3.Metadata) (*v3.Build, error) {
	req := c.NewRequest("POST", "/v3/builds")
	params := map[string]interface{}{
		"package": map[string]interface{}{
			"guid": packageGUID,
		},
	}
	if lifecycle != nil {
		params["lifecycle"] = lifecycle
	}
	if metadata != nil {
		params["metadata"] = metadata
	}
	req.obj = params

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 build")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 build, response code: %d", resp.StatusCode)
	}

	var build v3.Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "Error reading  Build JSON")
	}

	return &build, nil
}
