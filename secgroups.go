package cfclient

import (
	"encoding/json"
	"fmt"
	"io"
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
	Type        string `json:"type,omitempty"`        //ICMP type. Only valid if Protocol=="icmp"
	Ports       string `json:"ports"`                 //e.g. "4000-5000,9142"
	Destination string `json:"destination"`           //CIDR Format
	Description string `json:"description,omitempty"` //Optional description
	Log         bool   `json:"log,omitempty"`         //If true, log this rule
}

func (c *Client) ListSecGroups() (secGroups []SecGroup, err error) {
	requestURL := "/v2/security_groups?inline-relations-depth=1"
	for requestURL != "" {
		var secGroupResp SecGroupResponse
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)

		if err != nil {
			return nil, fmt.Errorf("Error requesting sec groups %v", err)
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading sec group request %v", string(resBody))
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
					return nil, err
				}
				for _, space := range spaces {
					secGroup.Entity.SpacesData = append(secGroup.Entity.SpacesData, space)
				}
			}
			secGroups = append(secGroups, secGroup.Entity)
		}

		requestURL = secGroupResp.NextUrl
		resp.Body.Close()
	}
	return secGroups, nil
}

func (secGroup *SecGroup) ListSpaceResources() ([]SpaceResource, error) {
	var spaceResources []SpaceResource
	requestURL := secGroup.SpacesURL
	for requestURL != "" {
		spaceResp, err := secGroup.c.getSpaceResponse(requestURL)
		if err != nil {
			return []SpaceResource{}, err
		}
		for i, spaceRes := range spaceResp.Resources {
			spaceRes.Entity.Guid = spaceRes.Meta.Guid
			spaceResp.Resources[i] = spaceRes
		}
		spaceResources = append(spaceResources, spaceResp.Resources...)
		requestURL = spaceResp.NextUrl
	}
	return spaceResources, nil
}

/*
CreateSecGroup contacts the CF endpoint for creating a new security group.
name: the name to give to the created security group
rules: A slice of rule objects that describe the rules that this security group enforces.
	This can technically be nil or an empty slice - we won't judge you
spaceGuids: The security group will be associated with the spaces specified by the contents of this slice.
	If nil, the security group will not be associated with any spaces initially.
*/
func (c *Client) CreateSecGroup(name string, rules []SecGroupRule, spaceGuids []string) (*SecGroup, error) {
	return c.secGroupCreateHelper("/v2/security_groups", "POST", name, rules, spaceGuids)
}

/*
UpdateSecGroup contacts the CF endpoint to update an existing security group.
guid: identifies the security group that you would like to update.
name: the new name to give to the security group
rules: A slice of rule objects that describe the rules that this security group enforces.
	If this is left nil, the rules will not be changed.
spaceGuids: The security group will be associated with the spaces specified by the contents of this slice.
	If nil, the space associations will not be changed.
*/
func (c *Client) UpdateSecGroup(guid, name string, rules []SecGroupRule, spaceGuids []string) (*SecGroup, error) {
	return c.secGroupCreateHelper("/v2/security_groups/"+guid, "PUT", name, rules, spaceGuids)
}

/*
DeleteSecGroup contacts the CF endpoint to delete an existing security group.
guid: Indentifies the security group to be deleted.
*/
func (c *Client) DeleteSecGroup(guid string) error {
	//Perform the DELETE and check for errors
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/security_groups/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 { //204 No Content
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

/*
GetSecGroup contacts the CF endpoint for fetching the info for a particular security group.
guid: Identifies the security group to fetch information from
*/
func (c *Client) GetSecGroup(guid string) (*SecGroup, error) {
	//Perform the GET and check for errors
	resp, err := c.DoRequest(c.NewRequest("GET", "/v2/security_groups/"+guid))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	//get the json out of the response body
	return respBodyToSecGroup(resp.Body, c)
}

/*
BindSecGroup contacts the CF endpoint to associate a space with a security group
secGUID: identifies the security group to add a space to
spaceGUID: identifies the space to associate
*/
func (c *Client) BindSecGroup(secGUID, spaceGUID string) error {
	//Perform the PUT and check for errors
	resp, err := c.DoRequest(c.NewRequest("PUT", fmt.Sprintf("/v2/security_groups/%s/spaces/%s", secGUID, spaceGUID)))
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 { //201 Created
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

/*
UnbindSecGroup contacts the CF endpoint to dissociate a space from a security group
secGUID: identifies the security group to remove a space from
spaceGUID: identifies the space to dissociate from the security group
*/
func (c *Client) UnbindSecGroup(secGUID, spaceGUID string) error {
	//Perform the DELETE and check for errors
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/security_groups/%s/spaces/%s", secGUID, spaceGUID)))
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 { //204 No Content
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

//Reads most security group response bodies into a SecGroup object
func respBodyToSecGroup(body io.ReadCloser, c *Client) (*SecGroup, error) {
	//get the json from the response body
	bodyRaw, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("Could not read response body: %s", err.Error())
	}
	jStruct := SecGroupResource{}
	//make it a SecGroup
	err = json.Unmarshal([]byte(bodyRaw), &jStruct)
	if err != nil {
		return nil, fmt.Errorf(`Could not unmarshal response body as json.
		body: %s
		error: %s`, bodyRaw, err.Error())
	}
	//pull a few extra fields from other places
	ret := jStruct.Entity
	ret.Guid = jStruct.Meta.Guid
	ret.c = c
	return &ret, nil
}

//Create and Update secGroup pretty much do the same thing, so this function abstracts those out.
func (c *Client) secGroupCreateHelper(url, method, name string, rules []SecGroupRule, spaceGuids []string) (*SecGroup, error) {
	req := c.NewRequest(method, url)
	//set up request body
	req.obj = map[string]interface{}{
		"name":        name,
		"rules":       rules,
		"space_guids": spaceGuids,
	}
	//fire off the request and check for problems
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 { // Both create and update should give 201 CREATED
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	//get the json from the response body
	return respBodyToSecGroup(resp.Body, c)
}
