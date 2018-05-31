package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MappingRequest struct {
	AppGUID   string `json:"app_guid"`
	RouteGUID string `json:"route_guid"`
	AppPort   int    `json:"app_port"`
}

type Mapping struct {
	Guid      string `json:"guid"`
	AppPort   int    `json:"app_port"`
	AppGUID   string `json:"app_guid"`
	RouteGUID string `json:"route_guid"`
	AppUrl    string `json:"app_url"`
	RouteUrl  string `json:"route_url"`
	c         *Client
}

type MappingResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Mapping `json:"entity"`
}

func (c *Client) MappingAppAndRoute(req MappingRequest) (*Mapping, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("POST", "/v2/route_mappings", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleMappingResp(resp)
}

func (c *Client) handleMappingResp(resp *http.Response) (*Mapping, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var mappingResource MappingResource
	err = json.Unmarshal(body, &mappingResource)
	if err != nil {
		return nil, err
	}
	return c.mergeMappingResource(mappingResource), nil
}

func (c *Client) mergeMappingResource(mapping MappingResource) *Mapping {
	mapping.Entity.Guid = mapping.Meta.Guid
	mapping.Entity.c = c
	return &mapping.Entity
}
