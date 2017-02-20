package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type Org struct {
	Guid                string `json:"guid"`
	Name                string `json:"name"`
	QuotaDefinitionGuid string `json:"quota_definition_guid"`
	c                   *Client
}

type OrgSummary struct {
	Guid   string             `json:"guid"`
	Name   string             `json:"name"`
	Status string             `json:"status"`
	Spaces []OrgSummarySpaces `json:"spaces"`
}

type OrgSummarySpaces struct {
	Guid         string `json:"guid"`
	Name         string `json:"name"`
	ServiceCount int    `json:"service_count"`
	AppCount     int    `json:"app_count"`
	MemDevTotal  int    `json:"mem_dev_total"`
	MemProdTotal int    `json:"mem_prod_total"`
}

type OrgRequest struct {
	Name                string `json:"name"`
	Status              string `json:"status,omitempty"`
	QuotaDefinitionGuid string `json:"quota_definition_guid,omitempty"`
}

func (c *Client) ListOrgsByQuery(query url.Values) ([]Org, error) {
	var orgs []Org
	requestUrl := "/v2/organizations?" + query.Encode()
	for {
		orgResp, err := c.getOrgResponse(requestUrl)
		if err != nil {
			return []Org{}, err
		}
		for _, org := range orgResp.Resources {
			org.Entity.Guid = org.Meta.Guid
			org.Entity.c = c
			orgs = append(orgs, org.Entity)
		}
		requestUrl = orgResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return orgs, nil
}

func (c *Client) ListOrgs() ([]Org, error) {
	return c.ListOrgsByQuery(nil)
}

func (c *Client) GetOrgByName(name string) (Org, error) {
	var org Org
	q := url.Values{}
	q.Set("q", "name:"+name)
	orgs, err := c.ListOrgsByQuery(q)
	if err != nil {
		return org, err
	}
	if len(orgs) == 0 {
		return org, fmt.Errorf("Unable to find org %s", name)
	}
	return orgs[0], nil
}

func (c *Client) GetOrgByGuid(guid string) (Org, error) {
	var orgRes OrgResource
	r := c.NewRequest("GET", "/v2/organizations/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Org{}, err
	}
	err = json.Unmarshal(body, &orgRes)
	if err != nil {
		return Org{}, err
	}
	orgRes.Entity.Guid = orgRes.Meta.Guid
	orgRes.Entity.c = c
	return orgRes.Entity, nil
}

func (c *Client) getOrgResponse(requestUrl string) (OrgResponse, error) {
	var orgResp OrgResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return OrgResponse{}, fmt.Errorf("Error requesting orgs %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return OrgResponse{}, fmt.Errorf("Error reading org request %v", err)
	}
	err = json.Unmarshal(resBody, &orgResp)
	if err != nil {
		return OrgResponse{}, fmt.Errorf("Error unmarshalling org %v", err)
	}
	return orgResp, nil
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

func (o *Org) Summary() (OrgSummary, error) {
	var orgSummary OrgSummary
	requestUrl := fmt.Sprintf("/v2/organizations/%s/summary", o.Guid)
	r := o.c.NewRequest("GET", requestUrl)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return OrgSummary{}, fmt.Errorf("Error requesting org summary %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return OrgSummary{}, fmt.Errorf("Error reading org summary body %v", err)
	}
	err = json.Unmarshal(resBody, &orgSummary)
	if err != nil {
		return OrgSummary{}, fmt.Errorf("Error unmarshalling org summary %v", err)
	}
	return orgSummary, nil
}

func (o *Org) Quota() (*OrgQuota, error) {
	var orgQuota *OrgQuota
	var orgQuotaResource OrgQuotasResource
	if o.QuotaDefinitionGuid == "" {
		return nil, nil
	}
	requestUrl := fmt.Sprintf("/v2/quota_definitions/%s", o.QuotaDefinitionGuid)
	r := o.c.NewRequest("GET", requestUrl)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return &OrgQuota{}, fmt.Errorf("Error requesting org quota %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &OrgQuota{}, fmt.Errorf("Error reading org quota body %v", err)
	}
	err = json.Unmarshal(resBody, &orgQuotaResource)
	if err != nil {
		return &OrgQuota{}, fmt.Errorf("Error unmarshalling org quota %v", err)
	}
	orgQuota = &orgQuotaResource.Entity
	orgQuota.Guid = orgQuotaResource.Meta.Guid
	orgQuota.c = o.c
	return orgQuota, nil
}

func (o *Org) AssociateManager(userGUID string) (Org, error) {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/managers/%s", o.Guid, userGUID)
	r := o.c.NewRequest("PUT", requestUrl)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) AssociateManagerByUsername(name string) (Org, error) {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/managers", o.Guid)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return Org{}, err
	}
	r := o.c.NewRequestWithBody("PUT", requestUrl, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) AssociateAuditorByUsername(name string) (Org, error) {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/auditors", o.Guid)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return Org{}, err
	}
	r := o.c.NewRequestWithBody("PUT", requestUrl, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) RemoveManager(userGUID string) error {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/managers/%s", o.Guid, userGUID)
	r := o.c.NewRequest("DELETE", requestUrl)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (o *Org) RemoveManagerByUsername(name string) error {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/managers", o.Guid)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return err
	}
	r := o.c.NewRequestWithBody("DELETE", requestUrl, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (o *Org) RemoveAuditorByUsername(name string) error {
	requestUrl := fmt.Sprintf("/v2/organizations/%s/auditors", o.Guid)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return err
	}
	r := o.c.NewRequestWithBody("DELETE", requestUrl, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrg(req OrgRequest) (Org, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return Org{}, err
	}
	r := c.NewRequestWithBody("POST", "/v2/organizations", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleOrgResp(resp)
}

func (c *Client) DeleteOrg(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/organizations/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) handleOrgResp(resp *http.Response) (Org, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Org{}, err
	}
	var orgResource OrgResource
	err = json.Unmarshal(body, &orgResource)
	if err != nil {
		return Org{}, err
	}
	org := orgResource.Entity
	org.Guid = orgResource.Meta.Guid
	org.c = c
	return org, nil
}
