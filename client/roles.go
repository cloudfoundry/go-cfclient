package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type RoleClient commonClient

func (c *RoleClient) CreateSpaceRole(spaceGUID, userGUID, roleType string) (*resource.Role, error) {
	spaceRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: spaceGUID}}
	userRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: userGUID}}
	req := c.client.NewRequest("POST", "/v3/roles")
	req.obj = resource.CreateSpaceRoleRequest{
		RoleType:      roleType,
		Relationships: resource.SpaceUserRelationships{Space: spaceRel, User: userRel},
	}
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while creating  role: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating  role, response code: %d", resp.StatusCode)
	}

	var role resource.Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("error reading  role: %w", err)
	}

	return &role, nil
}

func (c *RoleClient) CreateOrganizationRole(orgGUID, userGUID, roleType string) (*resource.Role, error) {
	orgRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: orgGUID}}
	userRel := resource.ToOneRelationship{Data: resource.Relationship{GUID: userGUID}}
	req := c.client.NewRequest("POST", "/v3/roles")
	req.obj = resource.CreateOrganizationRoleRequest{
		RoleType:      roleType,
		Relationships: resource.OrgUserRelationships{Org: orgRel, User: userRel},
	}
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while creating  role: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating  role, response code: %d", resp.StatusCode)
	}

	var role resource.Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("error reading  role: %w", err)
	}

	return &role, nil
}

// ListRolesByQuery retrieves roles based on query
func (c *RoleClient) ListRolesByQuery(query url.Values) ([]resource.Role, error) {
	var roles []resource.Role
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.client.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting  space roles: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing  space roles, response code: %d", resp.StatusCode)
		}

		var data resource.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list  space roles: %w", err)
		}

		roles = append(roles, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, fmt.Errorf("error parsing next page URL: %w", err)
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, nil
}

func (c *RoleClient) ListRoleUsersByQuery(query url.Values) ([]resource.User, error) {
	var users []resource.User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.client.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting  roles: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing  roles, response code: %d", resp.StatusCode)
		}

		var data resource.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list  roles: %w", err)
		}

		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, fmt.Errorf("error parsing next page URL: %w", err)
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}

func (c *RoleClient) ListRoleAndUsersByQuery(query url.Values) ([]resource.Role, []resource.User, error) {
	var roles []resource.Role
	var users []resource.User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.client.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, nil, fmt.Errorf("error requesting  roles: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("error listing  roles, response code: %d", resp.StatusCode)
		}

		var data resource.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, nil, fmt.Errorf("error parsing JSON from list  roles: %w", err)
		}

		roles = append(roles, data.Resources...)
		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing next page URL: %w", err)
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, users, nil
}

// ListSpaceRolesByGUID retrieves roles based on query
func (c *RoleClient) ListSpaceRolesByGUID(spaceGUID string) ([]resource.Role, []resource.User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["include"] = []string{"user"}
	return c.ListRoleAndUsersByQuery(query)
}

// ListSpaceRolesByGUIDAndType retrieves roles based on query
func (c *RoleClient) ListSpaceRolesByGUIDAndType(spaceGUID string, roleType string) ([]resource.User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListRoleUsersByQuery(query)
}

// ListOrganizationRolesByGUIDAndType retrieves roles based on query
func (c *RoleClient) ListOrganizationRolesByGUIDAndType(orgGUID string, roleType string) ([]resource.User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListRoleUsersByQuery(query)
}

// ListOrganizationRolesByGUID retrieves roles based on query
func (c *RoleClient) ListOrganizationRolesByGUID(orgGUID string) ([]resource.Role, []resource.User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["include"] = []string{"user"}
	return c.ListRoleAndUsersByQuery(query)
}

func (c *RoleClient) Delete(roleGUID string) error {
	req := c.client.NewRequest("DELETE", "/v3/roles/"+roleGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error while deleting role: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting role with GUID [%s], response code: %d", roleGUID, resp.StatusCode)
	}

	return nil
}
