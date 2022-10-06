package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

type BuildClient commonClient

func (c *BuildClient) GetByGUID(buildGUID string) (*resource.Build, error) {
	resp, err := c.client.DoRequest(c.client.NewRequest("GET", "/v3/builds/"+buildGUID))
	if err != nil {
		return nil, errors.Wrap(err, "error getting  build")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var build resource.Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "error reading  build JSON")
	}

	return &build, nil
}

func (c *BuildClient) Create(packageGUID string, lifecycle *resource.Lifecycle, metadata *resource.Metadata) (*resource.Build, error) {
	req := c.client.NewRequest("POST", "/v3/builds")
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

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating v3 build")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 build, response code: %d", resp.StatusCode)
	}

	var build resource.Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, errors.Wrap(err, "error reading  Build JSON")
	}

	return &build, nil
}
