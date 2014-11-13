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
	Guid string `json:"guid"`
	Name string `json:"name"`
}

func (c *Client) ListSpaces() []Space {
	var spaces []Space
	var spaceResp SpaceResponse
	r := c.newRequest("GET", "/v2/Spaces")
	resp, err := c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting Spaces %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		log.Printf("Error unmarshalling Space %v", err)
	}
	for _, space := range spaceResp.Resources {
		space.Entity.Guid = space.Meta.Guid
		spaces = append(spaces, space.Entity)
	}
	return spaces
}
