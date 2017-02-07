package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
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
	Guid                string      `json:"guid"`
	Name                string      `json:"name"`
	OrgURL              string      `json:"organization_url"`
	OrgData             OrgResource `json:"organization"`
	QuotaDefinitionGuid string      `json:"space_quota_definition_guid"`
	c                   *Client
}

type SpaceSummary struct {
	Guid     string           `json:"guid"`
	Name     string           `json:"name"`
	Apps     []AppSummary     `json:"apps"`
	Services []ServiceSummary `json:"services"`
}

type SpaceRoleResponse struct {
	Count     int                 `json:"total_results"`
	Pages     int                 `json:"total_pages"`
	NextUrl   string              `json:"next_url"`
	Resources []SpaceRoleResource `json:"resources"`
}

type SpaceRoleResource struct {
	Meta   Meta      `json:"metadata"`
	Entity SpaceRole `json:"entity"`
}

type SpaceRole struct {
	Guid                           string   `json:"guid"`
	Admin                          bool     `json:"admin"`
	Active                         bool     `json:"active"`
	DefaultSpaceGuid               string   `json:"default_space_guid"`
	Username                       string   `json:"username"`
	SpaceRoles                     []string `json:"space_roles"`
	SpacesUrl                      string   `json:"spaces_url"`
	OrganizationsUrl               string   `json:"organizations_url"`
	ManagedOrganizationsUrl        string   `json:"managed_organizations_url"`
	BillingManagedOrganizationsUrl string   `json:"billing_managed_organizations_url"`
	AuditedOrganizationsUrl        string   `json:"audited_organizations_url"`
	ManagedSpacesUrl               string   `json:"managed_spaces_url"`
	AuditedSpacesUrl               string   `json:"audited_spaces_url"`
	c                              *Client
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

func (s *Space) Quota() (*SpaceQuota, error) {
	var spaceQuota *SpaceQuota
	var spaceQuotaResource SpaceQuotasResource
	if s.QuotaDefinitionGuid == "" {
		return nil, nil
	}
	requestUrl := fmt.Sprintf("/v2/space_quota_definitions/%s", s.QuotaDefinitionGuid)
	r := s.c.NewRequest("GET", requestUrl)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return &SpaceQuota{}, fmt.Errorf("Error requesting space quota %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &SpaceQuota{}, fmt.Errorf("Error reading space quota body %v", err)
	}
	err = json.Unmarshal(resBody, &spaceQuotaResource)
	if err != nil {
		return &SpaceQuota{}, fmt.Errorf("Error unmarshalling space quota %v", err)
	}
	spaceQuota = &spaceQuotaResource.Entity
	spaceQuota.Guid = spaceQuotaResource.Meta.Guid
	spaceQuota.c = s.c
	return spaceQuota, nil
}

func (s *Space) Summary() (SpaceSummary, error) {
	var spaceSummary SpaceSummary
	requestUrl := fmt.Sprintf("/v2/spaces/%s/summary", s.Guid)
	r := s.c.NewRequest("GET", requestUrl)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return SpaceSummary{}, fmt.Errorf("Error requesting space summary %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceSummary{}, fmt.Errorf("Error reading space summary body %v", err)
	}
	err = json.Unmarshal(resBody, &spaceSummary)
	if err != nil {
		return SpaceSummary{}, fmt.Errorf("Error unmarshalling space summary %v", err)
	}
	return spaceSummary, nil
}

func (s *Space) Roles() ([]SpaceRole, error) {
	var roles []SpaceRole
	requestUrl := fmt.Sprintf("/v2/spaces/%s/user_roles", s.Guid)
	for {
		rolesResp, err := s.c.getSpaceRolesResponse(requestUrl)
		if err != nil {
			return roles, err
		}
		for _, role := range rolesResp.Resources {
			role.Entity.Guid = role.Meta.Guid
			role.Entity.c = s.c
			roles = append(roles, role.Entity)
		}
		requestUrl = rolesResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return roles, nil
}

func (c *Client) ListSpacesByQuery(query url.Values) ([]Space, error) {
	var spaces []Space
	requestUrl := "/v2/spaces?" + query.Encode()
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

func (c *Client) ListSpaces() ([]Space, error) {
	return c.ListSpacesByQuery(nil)
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

func (c *Client) getSpaceRolesResponse(requestUrl string) (SpaceRoleResponse, error) {
	var roleResp SpaceRoleResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return roleResp, fmt.Errorf("Error requesting space roles %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return roleResp, fmt.Errorf("Error reading space roles request %v", err)
	}
	err = json.Unmarshal(resBody, &roleResp)
	if err != nil {
		return roleResp, fmt.Errorf("Error unmarshalling space roles %v", err)
	}
	return roleResp, nil
}
