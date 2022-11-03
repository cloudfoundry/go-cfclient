package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type JobClient commonClient

// Get the specified job
func (c *JobClient) Get(guid string) (*resource.Job, error) {
	var job resource.Job
	err := c.client.get(path("/v3/jobs/%s", guid), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
