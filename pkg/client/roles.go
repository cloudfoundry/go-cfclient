package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) CreateSpaceRole(spaceGUID, userGUID, roleType string) (*v3.Role, error) {
	spaceRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: spaceGUID}}
	userRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: userGUID}}
	req := c.NewRequest("POST", "/v3/roles")
	req.obj = v3.CreateSpaceRoleRequest{
		RoleType:      roleType,
		Relationships: v3.SpaceUserRelationships{Space: spaceRel, User: userRel},
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating  role, response code: %d", resp.StatusCode)
	}

	var role v3.Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, errors.Wrap(err, "Error reading  role")
	}

	return &role, nil
}

func (c *Client) CreateOrganizationRole(orgGUID, userGUID, roleType string) (*v3.Role, error) {
	orgRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: orgGUID}}
	userRel := v3.ToOneRelationship{Data: v3.Relationship{GUID: userGUID}}
	req := c.NewRequest("POST", "/v3/roles")
	req.obj = v3.CreateOrganizationRoleRequest{
		RoleType:      roleType,
		Relationships: v3.OrgUserRelationships{Org: orgRel, User: userRel},
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating  role, response code: %d", resp.StatusCode)
	}

	var role v3.Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, errors.Wrap(err, "Error reading  role")
	}

	return &role, nil
}

// ListRolesByQuery retrieves roles based on query
func (c *Client) ListRolesByQuery(query url.Values) ([]v3.Role, error) {
	var roles []v3.Role
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  space roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  space roles, response code: %d", resp.StatusCode)
		}

		var data v3.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  space roles")
		}

		roles = append(roles, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, nil
}

func (c *Client) ListRoleUsersByQuery(query url.Values) ([]v3.User, error) {
	var users []v3.User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  roles, response code: %d", resp.StatusCode)
		}

		var data v3.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  roles")
		}

		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}

func (c *Client) ListRoleAndUsersByQuery(query url.Values) ([]v3.Role, []v3.User, error) {
	var roles []v3.Role
	var users []v3.User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error requesting  roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("Error listing  roles, response code: %d", resp.StatusCode)
		}

		var data v3.ListRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, nil, errors.Wrap(err, "Error parsing JSON from list  roles")
		}

		roles = append(roles, data.Resources...)
		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, users, nil
}

// ListSpaceRolesByGUID retrieves roles based on query
func (c *Client) ListSpaceRolesByGUID(spaceGUID string) ([]v3.Role, []v3.User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["include"] = []string{"user"}
	return c.ListRoleAndUsersByQuery(query)
}

// ListSpaceRolesByGUIDAndType retrieves roles based on query
func (c *Client) ListSpaceRolesByGUIDAndType(spaceGUID string, roleType string) ([]v3.User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListRoleUsersByQuery(query)
}

// ListSpaceRolesByGUIDAndType retrieves roles based on query
func (c *Client) ListOrganizationRolesByGUIDAndType(orgGUID string, roleType string) ([]v3.User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListRoleUsersByQuery(query)
}

// ListOrganizationRolesByGUID retrieves roles based on query
func (c *Client) ListOrganizationRolesByGUID(orgGUID string) ([]v3.Role, []v3.User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["include"] = []string{"user"}
	return c.ListRoleAndUsersByQuery(query)
}

func (c *Client) DeleteRole(roleGUID string) error {
	req := c.NewRequest("DELETE", "/v3/roles/"+roleGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting  role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting  role with GUID [%s], response code: %d", roleGUID, resp.StatusCode)
	}

	return nil
}
