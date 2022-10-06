package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type ServiceInstanceClient commonClient

func (c *ServiceInstanceClient) List() ([]resource.ServiceInstance, error) {
	return c.ListByQuery(nil)
}

func (c *ServiceInstanceClient) ListByQuery(query url.Values) ([]resource.ServiceInstance, error) {
	var svcInstances []resource.ServiceInstance
	requestURL := "/v3/service_instances"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.client.NewRequest("GET", requestURL)
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting service instances: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing service instances, response code: %d", resp.StatusCode)
		}

		var data resource.ListServiceInstancesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list service instances: %w", err)
		}

		svcInstances = append(svcInstances, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for service instances: %w", err)
		}
	}

	return svcInstances, nil
}
