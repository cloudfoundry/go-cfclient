package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) CreateOrganization(r v3.CreateOrganizationRequest) (*v3.Organization, error) {
	req := c.NewRequest("POST", "/v3/organizations")
	params := map[string]interface{}{
		"name": r.Name,
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 organization %s, response code: %d", r.Name, resp.StatusCode)
	}

	var organization v3.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (c *Client) GetOrganizationByGUID(organizationGUID string) (*v3.Organization, error) {
	req := c.NewRequest("GET", "/v3/organizations/"+organizationGUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization v3.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (c *Client) DeleteOrganization(organizationGUID string) error {
	req := c.NewRequest("DELETE", "/v3/organizations/"+organizationGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateOrganization(organizationGUID string, r v3.UpdateOrganizationRequest) (*v3.Organization, error) {
	req := c.NewRequest("PATCH", "/v3/organizations/"+organizationGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 organization %s, response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization v3.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (c *Client) ListOrganizationsByQuery(query url.Values) ([]v3.Organization, error) {
	var organizations []v3.Organization
	requestURL := "/v3/organizations"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 organizations")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 organizations, response code: %d", resp.StatusCode)
		}

		var data v3.ListOrganizationsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 organizations")
		}

		organizations = append(organizations, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 organizations")
		}
	}

	return organizations, nil
}
