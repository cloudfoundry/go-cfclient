package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type RouteClient commonClient

func (c *RouteClient) List() ([]resource.Route, error) {
	return c.ListByQuery(nil)
}

func (c *RouteClient) ListByQuery(query url.Values) ([]resource.Route, error) {
	var routes []resource.Route
	requestURL := "/v3/routes"
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

		var data resource.ListRouteResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list service instances: %w", err)
		}

		routes = append(routes, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for service instances: %w", err)
		}
	}

	return routes, nil
}

func (c *RouteClient) Create(
	spaceGUID string,
	domainGUID string,
	opt *resource.CreateRouteOptionalParameters,
) (*resource.Route, error) {

	spaceRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: spaceGUID}}
	domainRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: domainGUID}}

	req := c.client.NewRequest("POST", "/v3/routes")
	req.obj = resource.CreateRouteRequest{
		Relationships:                 resource.RouteRelationships{Space: spaceRel, Domain: domainRel},
		CreateRouteOptionalParameters: opt,
	}
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while creating route: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating route, response code: %d", resp.StatusCode)
	}

	var route resource.Route
	if err := json.NewDecoder(resp.Body).Decode(&route); err != nil {
		return nil, fmt.Errorf("error reading app package: %w", err)
	}

	return &route, nil
}
