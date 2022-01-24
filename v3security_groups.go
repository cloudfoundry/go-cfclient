package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type V3SecurityGroup struct {
	Name            string                           `json:"name,omitempty"`
	GUID            string                           `json:"guid,omitempty"`
	CreatedAt       string                           `json:"created_at,omitempty"`
	UpdatedAt       string                           `json:"updated_at,omitempty"`
	GloballyEnabled V3GloballyEnabled                `json:"globally_enabled,omitempty"`
	Rules           []V3Rule                         `json:"rules,omitempty"`
	Relationships   map[string]V3ToManyRelationships `json:"relationships,omitempty"`
	Links           map[string]Link                  `json:"links,omitempty"`
}

type V3GloballyEnabled struct {
	Running bool `json:"running"`
	Staging bool `json:"staging"`
}

type V3Rule struct {
	Protocol    string `json:"protocol,omitempty"`
	Destination string `json:"destination,omitempty"`
	Ports       string `json:"ports,omitempty"`
	Type        int    `json:"type,omitempty"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Log         bool   `json:"log,omitempty"`
}

type listV3SecurityGroupResponse struct {
	Pagination Pagination        `json:"pagination,omitempty"`
	Resources  []V3SecurityGroup `json:"resources,omitempty"`
}

func (c *Client) ListV3SecGroupsByQuery(query url.Values) ([]V3SecurityGroup, error) {
	var security_groups []V3SecurityGroup
	requestURL, err := url.Parse("/v3/security_groups")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 security groups")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 security groups, response code: %d", resp.StatusCode)
		}

		var data listV3SecurityGroupResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 security groups")
		}

		security_groups = append(security_groups, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return security_groups, nil
}
