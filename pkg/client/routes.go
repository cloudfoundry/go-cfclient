package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) ListRoutes() ([]v3.Route, error) {
	return c.ListRoutesByQuery(nil)
}

func (c *Client) ListRoutesByQuery(query url.Values) ([]v3.Route, error) {
	var routes []v3.Route
	requestURL := "/v3/routes"
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

		var data v3.ListRouteResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  service instances")
		}

		routes = append(routes, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for  service instances")
		}
	}

	return routes, nil
}

func (c *Client) CreateRoute(
	spaceGUID string,
	domainGUID string,
	opt *v3.CreateRouteOptionalParameters,
) (*v3.Route, error) {

	spaceRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: spaceGUID}}
	domainRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: domainGUID}}

	req := c.NewRequest("POST", "/v3/routes")
	req.obj = v3.CreateRouteRequest{
		Relationships:                 v3.RouteRelationships{Space: spaceRel, Domain: domainRel},
		CreateRouteOptionalParameters: opt,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  route")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating  route, response code: %d", resp.StatusCode)
	}

	var route v3.Route
	if err := json.NewDecoder(resp.Body).Decode(&route); err != nil {
		return nil, errors.Wrap(err, "Error reading  app package")
	}

	return &route, nil
}
