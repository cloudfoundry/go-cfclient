package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

// ListSecurityGroupsByQuery retrieves security groups based on query
func (c *Client) ListSecurityGroupsByQuery(query url.Values) ([]resource.SecurityGroup, error) {
	var securityGroups []resource.SecurityGroup
	requestURL, err := url.Parse("/v3/security_groups")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  security groups")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  security groups, response code: %d", resp.StatusCode)
		}

		var data resource.ListSecurityGroupResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  security groups")
		}

		securityGroups = append(securityGroups, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return securityGroups, nil
}

// CreateSecurityGroup creates security group from CreateSecurityGroupRequest
func (c *Client) CreateSecurityGroup(r resource.CreateSecurityGroupRequest) (*resource.SecurityGroup, error) {
	req := c.NewRequest("POST", "/v3/security_groups")

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	req.body = buf

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating  security group %s, response code: %d", r.Name, resp.StatusCode)
	}

	var securitygroup resource.SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securitygroup); err != nil {
		return nil, errors.Wrap(err, "Error reading  security group JSON")
	}

	return &securitygroup, nil
}

// DeleteSecurityGroup deletes security group by GUID
func (c *Client) DeleteSecurityGroup(GUID string) error {
	req := c.NewRequest("DELETE", "/v3/security_groups/"+GUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting  security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting  security group with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}
	return nil
}

// UpdateSecurityGroup updates security group by GUID and from UpdateSecurityGroupRequest
func (c *Client) UpdateSecurityGroup(GUID string, r resource.UpdateSecurityGroupRequest) (*resource.SecurityGroup, error) {
	req := c.NewRequest("PATCH", "/v3/security_groups/"+GUID)
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	req.body = buf

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating  security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating  security group %s, response code: %d", GUID, resp.StatusCode)
	}

	var securityGroup resource.SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securityGroup); err != nil {
		return nil, errors.Wrap(err, "Error reading  security group JSON")
	}

	return &securityGroup, nil
}

// GetSecurityGroupByGUID retrieves security group base on provided GUID
func (c *Client) GetSecurityGroupByGUID(GUID string) (*resource.SecurityGroup, error) {
	req := c.NewRequest("GET", "/v3/security_groups/"+GUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting  security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting  security group with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}

	var securityGroup resource.SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securityGroup); err != nil {
		return nil, errors.Wrap(err, "Error reading  security group JSON")
	}

	return &securityGroup, nil
}
