package cfclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ServiceResponse struct {
	Count     int               `json:"total_results"`
	Pages     int               `json:"total_pages"`
	Resources []ServiceResource `json:"resources"`
}

type ServiceResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Service `json:"entity"`
}

type Service struct {
	Guid  string `json:"guid"`
	Label string `json:"label"`
	c     *Client
}

func (c *Client) ListServices() []Service {
	var services []Service
	var serviceResp ServiceResponse
	r := c.newRequest("GET", "/v2/services")
	resp, err := c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting services %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading services request %v", resBody)
	}

	err = json.Unmarshal(resBody, &serviceResp)
	if err != nil {
		log.Printf("Error unmarshaling services %v", err)
	}
	for _, service := range serviceResp.Resources {
		service.Entity.Guid = service.Meta.Guid
		service.Entity.c = c
		services = append(services, service.Entity)
	}
	return services
}
