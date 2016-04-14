package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type Org struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
	c    *Client
}

func (c *Client) ListOrgs() ([]Org, error) {
	var orgs []Org
	var orgResp OrgResponse
	r := c.NewRequest("GET", "/v2/organizations")
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting organizations %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading organization request %v", resBody)
	}

	err = json.Unmarshal(resBody, &orgResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling organization %v", err)
	}
	for _, org := range orgResp.Resources {
		org.Entity.Guid = org.Meta.Guid
		org.Entity.c = c
		orgs = append(orgs, org.Entity)
	}
	return orgs, nil
}

func (c *Client) OrgSpaces(guid string) ([]Space, error) {
	var spaces []Space
	var spaceResp SpaceResponse
	path := fmt.Sprintf("/v2/organizations/%s/spaces", guid)
	r := c.NewRequest("GET", path)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting space %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		return nil, fmt.Errorf("Error space organization %v", err)
	}
	for _, space := range spaceResp.Resources {
		space.Entity.Guid = space.Meta.Guid
		spaces = append(spaces, space.Entity)
	}

	return spaces, nil

}
