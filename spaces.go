package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SpaceResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
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

func (s *Space) Org() (Org, error) {
	var orgResource OrgResource
	r := s.c.NewRequest("GET", s.OrgURL)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return Org{}, fmt.Errorf("Error requesting org %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Org{}, fmt.Errorf("Error reading org request %v", err)
	}

	err = json.Unmarshal(resBody, &orgResource)
	if err != nil {
		return Org{}, fmt.Errorf("Error unmarshaling org %v", err)
	}
	orgResource.Entity.Guid = orgResource.Meta.Guid
	orgResource.Entity.c = s.c
	return orgResource.Entity, nil
}

func (c *Client) ListSpaces() ([]Space, error) {
	var spaces []Space
	requestUrl := "/v2/spaces"
	for {
		spaceResp, err := c.getSpaceResponse(requestUrl)
		if err != nil {
			return []Space{}, err
		}
		for _, space := range spaceResp.Resources {
			space.Entity.Guid = space.Meta.Guid
			space.Entity.c = c
			spaces = append(spaces, space.Entity)
		}
		requestUrl = spaceResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return spaces, nil
}

func (c *Client) getSpaceResponse(requestUrl string) (SpaceResponse, error) {
	var spaceResp SpaceResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error requesting spaces %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error reading space request %v", err)
	}
	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		return SpaceResponse{}, fmt.Errorf("Error unmarshalling space %v", err)
	}
	return spaceResp, nil
}
