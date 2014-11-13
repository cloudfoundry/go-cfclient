package cf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// cfAppsResponse describes the Cloud Controller API result for a list of apps
type OrganizationResponse struct {
	Count     int                    `json:"total_results"`
	Pages     int                    `json:"total_pages"`
	Resources []OrganizationResource `json:"resources"`
}

type OrganizationResource struct {
	Meta   Meta         `json:"metadata"`
	Entity Organization `json:"entity"`
}

type Organization struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

func (c *Client) ListOrganizations() []Organization {
	var organizations []Organization
	var orgResp OrganizationResponse
	r := c.newRequest("GET", "/v2/organizations")
	resp, err := c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting organizations %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading organization request %v", resBody)
	}

	err = json.Unmarshal(resBody, &orgResp)
	if err != nil {
		log.Printf("Error unmarshalling organization %v", err)
	}
	for _, org := range orgResp.Resources {
		org.Entity.Guid = org.Meta.Guid
		organizations = append(organizations, org.Entity)
	}
	return organizations
}
