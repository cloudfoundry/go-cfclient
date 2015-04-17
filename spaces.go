package cfclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type SpaceResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	Resources []SpaceResource `json:"resources"`
}

type SpaceResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Space `json:"entity"`
}

type Space struct {
	Guid    string      `json:"guid"`
	Name    string      `json:"name"`
	OrgURL  string      `json:"organization_url"`
	OrgData OrgResource `json:"organization"`
	c       *Client
}

func (s *Space) Org() Org {
	var orgResource OrgResource
	r := s.c.newRequest("GET", s.OrgURL)
	resp, err := s.c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting org %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading org request %v", resBody)
	}

	err = json.Unmarshal(resBody, &orgResource)
	if err != nil {
		log.Printf("Error unmarshaling org %v", err)
	}
	orgResource.Entity.Guid = orgResource.Meta.Guid
	orgResource.Entity.c = s.c
	return orgResource.Entity
}

func (c *Client) ListSpaces() []Space {
	var spaces []Space
	var spaceResp SpaceResponse
	r := c.newRequest("GET", "/v2/spaces")
	resp, err := c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting spaces %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		log.Printf("Error unmarshalling space %v", err)
	}
	for _, space := range spaceResp.Resources {
		space.Entity.Guid = space.Meta.Guid
		space.Entity.c = c
		spaces = append(spaces, space.Entity)
	}
	return spaces
}
