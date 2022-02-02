package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type V3Route struct {
	Guid          string                         `json:"guid"`
	Host          string                         `json:"host"`
	Path          string                         `json:"path"`
	Url           string                         `json:"url"`
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
	Metadata      Metadata                       `json:"metadata"`
	Relationships map[string]V3ToOneRelationship `json:"relationships"`
	Links         map[string]Link                `json:"links"`
}

type listV3RouteResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Route  `json:"resources,omitempty"`
}

func (c *Client) ListV3Routes() ([]V3Route, error) {
	return c.ListV3RoutesByQuery(nil)
}

func (c *Client) ListV3RoutesByQuery(query url.Values) ([]V3Route, error) {
	var routes []V3Route
	requestURL := "/v3/routes"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 service instances")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing v3 service instances, response code: %d", resp.StatusCode)
		}

		var data listV3RouteResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 service instances")
		}

		routes = append(routes, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 service instances")
		}
	}

	return routes, nil
}
