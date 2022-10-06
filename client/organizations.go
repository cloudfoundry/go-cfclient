package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

type OrgClient commonClient

func (o *OrgClient) Create(r resource.CreateOrganizationRequest) (*resource.Organization, error) {
	req := o.client.NewRequest("POST", "/v3/organizations")
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
	resp, err := o.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 organization")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 organization %s, response code: %d", r.Name, resp.StatusCode)
	}

	var organization resource.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (o *OrgClient) GetByGUID(organizationGUID string) (*resource.Organization, error) {
	req := o.client.NewRequest("GET", "/v3/organizations/"+organizationGUID)

	resp, err := o.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 organization")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization resource.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (o *OrgClient) Delete(organizationGUID string) error {
	req := o.client.NewRequest("DELETE", "/v3/organizations/"+organizationGUID)
	resp, err := o.client.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 organization")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	return nil
}

func (o *OrgClient) Update(organizationGUID string, r resource.UpdateOrganizationRequest) (*resource.Organization, error) {
	req := o.client.NewRequest("PATCH", "/v3/organizations/"+organizationGUID)
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

	resp, err := o.client.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 organization")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 organization %s, response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization resource.Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (o *OrgClient) ListByQuery(query url.Values) ([]resource.Organization, error) {
	var organizations []resource.Organization
	requestURL := "/v3/organizations"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := o.client.NewRequest("GET", requestURL)
		resp, err := o.client.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 organizations")
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 organizations, response code: %d", resp.StatusCode)
		}

		var data resource.ListOrganizationsResponse
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
