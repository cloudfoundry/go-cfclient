package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type RoutesResponse struct {
	Count     int              `json:"total_results"`
	Pages     int              `json:"total_pages"`
	NextUrl   string           `json:"next_url"`
	Resources []RoutesResource `json:"resources"`
}

type RoutesResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Route `json:"entity"`
}

type Route struct {
	Guid                string `json:"guid"`
	Host                string `json:"host"`
	Path                string `json:"path"`
	DomainGuid          string `json:"domain_guid"`
	SpaceGuid           string `json:"space_guid"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
	Port                int    `json:"port"`
	c                   *Client
}

func (c *Client) ListRoutesByQuery(query url.Values) ([]Route, error) {
	var routes []Route
	requestUrl := "/v2/routes?" + query.Encode()
	for {
		routesResp, err := c.getRoutesResponse(requestUrl)
		if err != nil {
			return []Route{}, err
		}
		for _, route := range routesResp.Resources {
			route.Entity.Guid = route.Meta.Guid
			route.Entity.c = c
			routes = append(routes, route.Entity)
		}
		requestUrl = routesResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return routes, nil
}

func (c *Client) ListRoutes() ([]Route, error) {
	return c.ListRoutesByQuery(nil)
}

func (c *Client) getRoutesResponse(requestUrl string) (RoutesResponse, error) {
	var routesResp RoutesResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return RoutesResponse{}, fmt.Errorf("Error requesting routes %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RoutesResponse{}, fmt.Errorf("Error reading routes body %v", err)
	}
	err = json.Unmarshal(resBody, &routesResp)
	if err != nil {
		return RoutesResponse{}, fmt.Errorf("Error unmarshalling routes %v", err)
	}
	return routesResp, nil
}
