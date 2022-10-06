package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type ServiceCredentialBindingClient commonClient

// List retrieves all service credential bindings
func (c *ServiceCredentialBindingClient) List() ([]resource.ServiceCredentialBindings, error) {
	return c.ListByQuery(nil)
}

// ListByQuery retrieves service credential bindings using a query
func (c *ServiceCredentialBindingClient) ListByQuery(query url.Values) ([]resource.ServiceCredentialBindings, error) {
	var svcCredentialBindings []resource.ServiceCredentialBindings
	requestURL := "/v3/service_credential_bindings"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.client.NewRequest("GET", requestURL)
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting service credential bindings: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing service credential bindings, response code: %d", resp.StatusCode)
		}

		var data resource.ListServiceCredentialBindingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list service credential bindings: %w", err)
		}

		svcCredentialBindings = append(svcCredentialBindings, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for service credential bindings: %w", err)
		}
	}

	return svcCredentialBindings, nil
}

// Get retrieves the service credential binding based on the provided guid
func (c *ServiceCredentialBindingClient) Get(GUID string) (*resource.ServiceCredentialBindings, error) {
	requestURL := fmt.Sprintf("/v3/service_credential_bindings/%s", GUID)
	req := c.client.NewRequest("GET", requestURL)
	resp, err := c.client.DoRequest(req)

	if err != nil {
		return nil, fmt.Errorf("error while getting service credential binding: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting service credential binding with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}

	var svcCredentialBindings resource.ServiceCredentialBindings
	if err := json.NewDecoder(resp.Body).Decode(&svcCredentialBindings); err != nil {
		return nil, fmt.Errorf("error reading service credential binding JSON: %w", err)
	}

	return &svcCredentialBindings, nil
}
