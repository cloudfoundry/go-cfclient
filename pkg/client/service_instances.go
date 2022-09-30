package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) ListServiceInstances() ([]v3.ServiceInstance, error) {
	return c.ListServiceInstancesByQuery(nil)
}

func (c *Client) ListServiceInstancesByQuery(query url.Values) ([]v3.ServiceInstance, error) {
	var svcInstances []v3.ServiceInstance
	requestURL := "/v3/service_instances"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  service instances")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing  service instances, response code: %d", resp.StatusCode)
		}

		var data v3.ListServiceInstancesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  service instances")
		}

		svcInstances = append(svcInstances, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for  service instances")
		}
	}

	return svcInstances, nil
}
