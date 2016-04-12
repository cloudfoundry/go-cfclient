package cfclient

import (
	"encoding/json"
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

func (c *Client) ListSecGroups() []SecGroup {
	var secGroups []SecGroup

	requestUrl := "/v2/security_groups?inline-relations-depth=1"
	for {
		var secGroupResp SecGroupResponse
		r := c.newRequest("GET", requestUrl)
		resp, err := c.doRequest(r)

		if err != nil {
			log.Printf("Error requesting sec groups %v", err)
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading sec group request %v", resBody)
		}

		err = json.Unmarshal(resBody, &secGroupResp)
		if err != nil {
			log.Printf("Error unmarshaling sec group %v", err)
		}

		for _, secGroup := range secGroupResp.Resources {
			secGroup.Entity.Guid = secGroup.Meta.Guid
			secGroup.Entity.c = c
			for i, space := range secGroup.Entity.SpacesData {
				space.Entity.Guid = space.Meta.Guid
				secGroup.Entity.SpacesData[i] = space
			}
			if len(secGroup.Entity.SpacesData) == 0 {
				spaces := secGroup.Entity.ListSpaceResources()
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
	return secGroups
}

func (secGroup *SecGroup) ListSpaceResources() []SpaceResource {
	var spaceResources []SpaceResource
	requestUrl := secGroup.SpacesURL
	for {
		var spaceResp = secGroup.c.getSpaceResponse(requestUrl)
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
	return spaceResources
}
