package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServiceInstanceResource struct {
	Meta   Meta            `json:"metadata"`
	Entity ServiceInstance `json:"entity"`
}

type ServiceInstance struct {
	Guid               string `json:"guid"`
	ServicePlanGuid    string `json:"service_plan_guid"`
	Name               string `json:"name"`
	SpaceGuid          string `json:"space_guid"`
	DashboardUrl       string `json:"dashboard_url"`
	Type               string `json:"type"`
	SpaceUrl           string `json:"space_url"`
	ServicePlanUrl     string `json:"service_plan_url"`
	ServiceBindingsUrl string `json:"service_bindings_url"`
	ServiceKeysUrl     string `json:"service_keys_url"`
	RoutesUrl          string `json:"routes_url"`
	c                  *Client
}

func (c *Client) ServiceInstanceByGuid(guid string) (ServiceInstance, error) {
	var sir ServiceInstanceResource
	req := c.NewRequest("GET", "/v2/service_instances/"+guid)
	res, err := c.DoRequest(req)
	if err != nil {
		return ServiceInstance{}, fmt.Errorf("Error requesting service instance: %s", err)
	}
	if res.StatusCode >= 400 {
		return ServiceInstance{}, fmt.Errorf("Error requesting service instance '%s': %s", guid, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ServiceInstance{}, fmt.Errorf("Error reading service instance response: %s", err)
	}
	err = json.Unmarshal(data, &sir)
	if err != nil {
		return ServiceInstance{}, fmt.Errorf("Error JSON parsing service instance response: %s", err)
	}
	sir.Entity.Guid = sir.Meta.Guid
	sir.Entity.c = c
	return sir.Entity, nil
}
