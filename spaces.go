package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type SpaceRequest struct {
	Name               string   `json:"name"`
	OrganizationGuid   string   `json:"organization_guid"`
	DeveloperGuid      []string `json:"developer_guids,omitempty"`
	ManagerGuid        []string `json:"manager_guids,omitempty"`
	AuditorGuid        []string `json:"auditor_guids,omitempty"`
	DomainGuid         []string `json:"domain_guids,omitempty"`
	SecurityGroupGuids []string `json:"security_group_guids,omitempty"`
	SpaceQuotaDefGuid  string   `json:"space_quota_definition_guid,omitempty"`
	AllowSSH           bool     `json:"allow_ssh,omitempty"`
}

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

type ServicePlanEntity struct {
	Name                string                  `json:"name"`
	Free                bool                    `json:"free"`
	Public              bool                    `json:"public"`
	Active              bool                    `json:"active"`
	Description         string                  `json:"description"`
	ServiceOfferingGUID string                  `json:"service_guid"`
	ServiceOffering     ServiceOfferingResource `json:"service"`
}

type ServiceOfferingExtra struct {
	DisplayName      string `json:"displayName"`
	DocumentationURL string `json:"documentationURL"`
	LongDescription  string `json:"longDescription"`
}

type ServiceOfferingEntity struct {
	Label        string
	Description  string
	Provider     string        `json:"provider"`
	BrokerGUID   string        `json:"service_broker_guid"`
	Requires     []string      `json:"requires"`
	ServicePlans []interface{} `json:"service_plans"`
	Extra        ServiceOfferingExtra
}

type ServiceOfferingResource struct {
	Metadata Meta
	Entity   ServiceOfferingEntity
}

type ServiceOfferingResponse struct {
	Count     int                       `json:"total_results"`
	Pages     int                       `json:"total_pages"`
	NextUrl   string                    `json:"next_url"`
	PrevUrl   string                    `json:"prev_url"`
	Resources []ServiceOfferingResource `json:"resources"`
}

type SpaceUserResponse struct {
	Count     int            `json:"total_results"`
	Pages     int            `json:"total_pages"`
	NextURL   string         `json:"next_url"`
	Resources []UserResource `json:"resources"`
}

type Space struct {
	Guid                string      `json:"guid"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	Name                string      `json:"name"`
	OrganizationGuid    string      `json:"organization_guid"`
	OrgURL              string      `json:"organization_url"`
	OrgData             OrgResource `json:"organization"`
	QuotaDefinitionGuid string      `json:"space_quota_definition_guid"`
	AllowSSH            bool        `json:"allow_ssh"`
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
		return Org{}, errors.Wrap(err, "Error requesting org")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Org{}, errors.Wrap(err, "Error reading org request")
	}

	err = json.Unmarshal(resBody, &orgResource)
	if err != nil {
		return Org{}, errors.Wrap(err, "Error unmarshaling org")
	}
	return s.c.mergeOrgResource(orgResource), nil
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
		return &SpaceQuota{}, errors.Wrap(err, "Error requesting space quota")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &SpaceQuota{}, errors.Wrap(err, "Error reading space quota body")
	}
	err = json.Unmarshal(resBody, &spaceQuotaResource)
	if err != nil {
		return &SpaceQuota{}, errors.Wrap(err, "Error unmarshalling space quota")
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
		return SpaceSummary{}, errors.Wrap(err, "Error requesting space summary")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceSummary{}, errors.Wrap(err, "Error reading space summary body")
	}
	err = json.Unmarshal(resBody, &spaceSummary)
	if err != nil {
		return SpaceSummary{}, errors.Wrap(err, "Error unmarshalling space summary")
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

func (c *Client) CreateSpace(req SpaceRequest) (Space, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return Space{}, err
	}
	r := c.NewRequestWithBody("POST", "/v2/spaces", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Space{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Space{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleSpaceResp(resp)
}

func (c *Client) UpdateSpace(spaceGUID string, req SpaceRequest) (Space, error) {
	space := Space{Guid: spaceGUID, c: c}
	return space.Update(req)
}

func (c *Client) DeleteSpace(guid string, recursive, async bool) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/spaces/%s?recursive=%t&async=%t", guid, recursive, async)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting space %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) ListSpaceManagersByQuery(spaceGUID string, query url.Values) ([]User, error) {
	return c.listSpaceUsersByRoleAndQuery(spaceGUID, "managers", query)
}

func (c *Client) ListSpaceManagers(spaceGUID string) ([]User, error) {
	return c.ListSpaceManagersByQuery(spaceGUID, nil)
}

func (c *Client) ListSpaceAuditorsByQuery(spaceGUID string, query url.Values) ([]User, error) {
	return c.listSpaceUsersByRoleAndQuery(spaceGUID, "auditors", query)
}

func (c *Client) ListSpaceAuditors(spaceGUID string) ([]User, error) {
	return c.ListSpaceAuditorsByQuery(spaceGUID, nil)
}

func (c *Client) ListSpaceDevelopersByQuery(spaceGUID string, query url.Values) ([]User, error) {
	return c.listSpaceUsersByRoleAndQuery(spaceGUID, "developers", query)
}

func (c *Client) listSpaceUsersByRoleAndQuery(spaceGUID, role string, query url.Values) ([]User, error) {
	var users []User
	requestURL := fmt.Sprintf("/v2/spaces/%s/%s?%s", spaceGUID, role, query.Encode())
	for {
		userResp, err := c.getUserResponse(requestURL)
		if err != nil {
			return []User{}, err
		}
		for _, u := range userResp.Resources {
			users = append(users, c.mergeUserResource(u))
		}
		requestURL = userResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) ListSpaceDevelopers(spaceGUID string) ([]User, error) {
	return c.ListSpaceDevelopersByQuery(spaceGUID, nil)
}

func (c *Client) AssociateSpaceDeveloperByUsername(spaceGUID, name string) (Space, error) {
	space := Space{Guid: spaceGUID, c: c}
	return space.AssociateDeveloperByUsername(name)
}

func (c *Client) RemoveSpaceDeveloperByUsername(spaceGUID, name string) error {
	space := Space{Guid: spaceGUID, c: c}
	return space.RemoveDeveloperByUsername(name)
}

func (c *Client) AssociateSpaceAuditorByUsername(spaceGUID, name string) (Space, error) {
	space := Space{Guid: spaceGUID, c: c}
	return space.AssociateAuditorByUsername(name)
}

func (c *Client) RemoveSpaceAuditorByUsername(spaceGUID, name string) error {
	space := Space{Guid: spaceGUID, c: c}
	return space.RemoveAuditorByUsername(name)
}

func (c *Client) AssociateSpaceManagerByUsername(spaceGUID, name string) (Space, error) {
	space := Space{Guid: spaceGUID, c: c}
	return space.AssociateManagerByUsername(name)
}

func (c *Client) RemoveSpaceManagerByUsername(spaceGUID, name string) error {
	space := Space{Guid: spaceGUID, c: c}
	return space.RemoveManagerByUsername(name)
}

func (s *Space) AssociateDeveloperByUsername(name string) (Space, error) {
	return s.associateUserByRole(name, "developers")
}

func (s *Space) RemoveDeveloperByUsername(name string) error {
	return s.removeUserByRole(name, "developers")
}
func (s *Space) AssociateAuditorByUsername(name string) (Space, error) {
	return s.associateUserByRole(name, "auditors")
}

func (s *Space) RemoveAuditorByUsername(name string) error {
	return s.removeUserByRole(name, "auditors")
}

func (s *Space) AssociateManagerByUsername(name string) (Space, error) {
	return s.associateUserByRole(name, "managers")
}

func (s *Space) RemoveManagerByUsername(name string) error {
	return s.removeUserByRole(name, "managers")
}

func (s *Space) associateUserByRole(name, role string) (Space, error) {
	requestUrl := fmt.Sprintf("/v2/spaces/%s/%s", s.Guid, role)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return Space{}, err
	}
	r := s.c.NewRequestWithBody("PUT", requestUrl, buf)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return Space{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Space{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return s.c.handleSpaceResp(resp)
}

func (s *Space) removeUserByRole(name, role string) error {
	requestUrl := fmt.Sprintf("/v2/spaces/%s/%s", s.Guid, role)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]string{"username": name})
	if err != nil {
		return err
	}
	r := s.c.NewRequestWithBody("DELETE", requestUrl, buf)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) ListSpaceSecGroups(spaceGUID string) (secGroups []SecGroup, err error) {
	space := Space{Guid: spaceGUID, c: c}
	return space.ListSecGroups()
}

func (s *Space) ListSecGroups() (secGroups []SecGroup, err error) {
	requestURL := fmt.Sprintf("/v2/spaces/%s/security_groups?inline-relations-depth=1", s.Guid)
	for requestURL != "" {
		var secGroupResp SecGroupResponse
		r := s.c.NewRequest("GET", requestURL)
		resp, err := s.c.DoRequest(r)

		if err != nil {
			return nil, errors.Wrap(err, "Error requesting sec groups")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading sec group response body")
		}

		err = json.Unmarshal(resBody, &secGroupResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling sec group")
		}

		for _, secGroup := range secGroupResp.Resources {
			secGroup.Entity.Guid = secGroup.Meta.Guid
			secGroup.Entity.c = s.c
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

func (s *Space) GetServiceOfferings() (ServiceOfferingResponse, error) {
	var response ServiceOfferingResponse
	requestURL := fmt.Sprintf("/v2/spaces/%s/services", s.Guid)
	req := s.c.NewRequest("GET", requestURL)

	resp, err := s.c.DoRequest(req)
	if err != nil {
		return ServiceOfferingResponse{}, errors.Wrap(err, "Error requesting service offerings")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ServiceOfferingResponse{}, errors.Wrap(err, "Error reading service offering response")
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return ServiceOfferingResponse{}, errors.Wrap(err, "Error unmarshalling service offering response")
	}

	return response, nil
}

func (s *Space) Update(req SpaceRequest) (Space, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return Space{}, err
	}
	r := s.c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/spaces/%s", s.Guid), buf)
	resp, err := s.c.DoRequest(r)
	if err != nil {
		return Space{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Space{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return s.c.handleSpaceResp(resp)
}

func (c *Client) ListSpacesByQuery(query url.Values) ([]Space, error) {
	return c.fetchSpaces("/v2/spaces?" + query.Encode())
}

func (c *Client) ListSpaces() ([]Space, error) {
	return c.ListSpacesByQuery(nil)
}

func (c *Client) fetchSpaces(requestUrl string) ([]Space, error) {
	var spaces []Space
	for {
		spaceResp, err := c.getSpaceResponse(requestUrl)
		if err != nil {
			return []Space{}, err
		}
		for _, space := range spaceResp.Resources {
			spaces = append(spaces, c.mergeSpaceResource(space))
		}
		requestUrl = spaceResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return spaces, nil
}

func (c *Client) GetSpaceByName(spaceName string, orgGuid string) (space Space, err error) {
	query := url.Values{}
	query.Add("q", fmt.Sprintf("organization_guid:%s", orgGuid))
	query.Add("q", fmt.Sprintf("name:%s", spaceName))
	spaces, err := c.ListSpacesByQuery(query)
	if err != nil {
		return
	}

	if len(spaces) == 0 {
		return space, fmt.Errorf("No space found with name: `%s` in org with GUID: `%s`", spaceName, orgGuid)
	}

	return spaces[0], nil

}

func (c *Client) GetSpaceByGuid(spaceGUID string) (Space, error) {
	requestUrl := fmt.Sprintf("/v2/spaces/%s", spaceGUID)
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Space{}, errors.Wrap(err, "Error requesting space info")
	}
	return c.handleSpaceResp(resp)
}

func (c *Client) getSpaceResponse(requestUrl string) (SpaceResponse, error) {
	var spaceResp SpaceResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return SpaceResponse{}, errors.Wrap(err, "Error requesting spaces")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceResponse{}, errors.Wrap(err, "Error reading space request")
	}
	err = json.Unmarshal(resBody, &spaceResp)
	if err != nil {
		return SpaceResponse{}, errors.Wrap(err, "Error unmarshalling space")
	}
	return spaceResp, nil
}

func (c *Client) getSpaceRolesResponse(requestUrl string) (SpaceRoleResponse, error) {
	var roleResp SpaceRoleResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return roleResp, errors.Wrap(err, "Error requesting space roles")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return roleResp, errors.Wrap(err, "Error reading space roles request")
	}
	err = json.Unmarshal(resBody, &roleResp)
	if err != nil {
		return roleResp, errors.Wrap(err, "Error unmarshalling space roles")
	}
	return roleResp, nil
}

func (c *Client) handleSpaceResp(resp *http.Response) (Space, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Space{}, err
	}
	var spaceResource SpaceResource
	err = json.Unmarshal(body, &spaceResource)
	if err != nil {
		return Space{}, err
	}
	return c.mergeSpaceResource(spaceResource), nil
}

func (c *Client) mergeSpaceResource(space SpaceResource) Space {
	space.Entity.Guid = space.Meta.Guid
	space.Entity.CreatedAt = space.Meta.CreatedAt
	space.Entity.UpdatedAt = space.Meta.UpdatedAt
	space.Entity.c = c
	return space.Entity
}

type serviceOfferingExtra ServiceOfferingExtra

func (resource *ServiceOfferingExtra) UnmarshalJSON(rawData []byte) error {
	if string(rawData) == "null" {
		return nil
	}

	extra := serviceOfferingExtra{}

	unquoted, err := strconv.Unquote(string(rawData))
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(unquoted), &extra)
	if err != nil {
		return err
	}

	*resource = ServiceOfferingExtra(extra)

	return nil
}
