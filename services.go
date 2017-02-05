package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type servicesResponse struct {
	Count     int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	Resources []servicesResource `json:"resources"`
}

type servicesResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Service `json:"entity"`
}

type Service struct {
	Guid  string `json:"guid"`
	Label string `json:"label"`
	c     *Client
}

type ServiceSummary struct {
	Guid          string `json:"guid"`
	Name          string `json:"name"`
	BoundAppCount int    `json:"bound_app_count"`
}

func (c *Client) ListServicesByQuery(query url.Values) ([]Service, error) {
	var services []Service
	var serviceResp servicesResponse
	r := c.NewRequest("GET", "/v2/services?"+query.Encode())
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting services %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading services request: %v", err)
	}

	err = json.Unmarshal(resBody, &serviceResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling services %v", err)
	}
	for _, service := range serviceResp.Resources {
		service.Entity.Guid = service.Meta.Guid
		service.Entity.c = c
		services = append(services, service.Entity)
	}
	return services, nil
}

func (c *Client) ListServices() ([]Service, error) {
	return c.ListServicesByQuery(nil)
}
