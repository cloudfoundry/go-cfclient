package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// V3Role implements role object. Roles control access to resources in organizations and spaces. Roles are assigned to users.
type V3Role struct {
	GUID          string                         `json:"guid,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Type          string                         `json:"type,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
}

type listV3SpaceRolesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Role   `json:"resources,omitempty"`
}

// ListV3RolesByQuery retrieves roles based on query
func (c *Client) ListV3RolesByQuery(query url.Values) ([]V3Role, error) {
	var roles []V3Role
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 space roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 space roles, response code: %d", resp.StatusCode)
		}

		var data listV3SpaceRolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 space roles")
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
