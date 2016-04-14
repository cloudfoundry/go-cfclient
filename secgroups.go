package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type SecGroupResponse struct {
	Count     int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	NextUrl   string             `json:"next_url"`
	Resources []SecGroupResource `json:"resources"`
}

type SecGroupResource struct {
	Meta   Meta     `json:"metadata"`
	Entity SecGroup `json:"entity"`
}

type SecGroup struct {
	Guid       string          `json:"guid"`
	Name       string          `json:"name"`
	Rules      []SecGroupRule  `json:"rules"`
	Running    bool            `json:"running_default"`
	Staging    bool            `json:"staging_default"`
	SpacesURL  string          `json:"spaces_url"`
	SpacesData []SpaceResource `json:"spaces"`
	c          *Client
}

type SecGroupRule struct {
	Protocol    string `json:"protocol"`
	Ports       string `json:"ports"`
	Destination string `json:"destination"`
}

func (c *Client) ListSecGroups() (secGroups []SecGroup, err error) {

	requestUrl := "/v2/security_groups?inline-relations-depth=1"
	for {
		var secGroupResp SecGroupResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)

		if err != nil {
			return nil, fmt.Errorf("Error requesting sec groups %v", err)
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading sec group request %v", resBody)
		}

		err = json.Unmarshal(resBody, &secGroupResp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling sec group %v", err)
		}

		for _, secGroup := range secGroupResp.Resources {
			secGroup.Entity.Guid = secGroup.Meta.Guid
			secGroup.Entity.c = c
			for i, space := range secGroup.Entity.SpacesData {
				space.Entity.Guid = space.Meta.Guid
				secGroup.Entity.SpacesData[i] = space
			}
			if len(secGroup.Entity.SpacesData) == 0 {
				spaces, err := secGroup.Entity.ListSpaceResources()
				if err != nil {
					return secGroups, nil
				}
				for _, space := range spaces {
					secGroup.Entity.SpacesData = append(secGroup.Entity.SpacesData, space)
				}
			}
			secGroups = append(secGroups, secGroup.Entity)
		}

		requestUrl = secGroupResp.NextUrl
		if requestUrl == "" {
			break
		}
		resp.Body.Close()
	}
	return secGroups, nil
}

func (secGroup *SecGroup) ListSpaceResources() ([]SpaceResource, error) {
	var spaceResources []SpaceResource
	requestUrl := secGroup.SpacesURL
	for {
		spaceResp, err := secGroup.c.getSpaceResponse(requestUrl)
		if err != nil {
			return []SpaceResource{}, err
		}
		for i, spaceRes := range spaceResp.Resources {
			spaceRes.Entity.Guid = spaceRes.Meta.Guid
			spaceResp.Resources[i] = spaceRes
		}
		spaceResources = append(spaceResources, spaceResp.Resources...)
		requestUrl = spaceResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return spaceResources, nil
}
