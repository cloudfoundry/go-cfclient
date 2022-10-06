package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type SpaceClient commonClient

func (c *SpaceClient) Create(r resource.CreateSpaceRequest) (*resource.Space, error) {
	req := c.client.NewRequest("POST", "/v3/spaces")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"organization": resource.ToOneRelationship{
				Data: resource.Relationship{
					GUID: r.OrgGUID,
				},
			},
		},
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while creating space: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating space %s, response code: %d", r.Name, resp.StatusCode)
	}

	var space resource.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, fmt.Errorf("error reading space JSON: %w", err)
	}

	return &space, nil
}

func (c *SpaceClient) Get(spaceGUID string) (*resource.Space, error) {
	req := c.client.NewRequest("GET", "/v3/spaces/"+spaceGUID)

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while getting space: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	var space resource.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, fmt.Errorf("error reading space JSON: %w", err)
	}

	return &space, nil
}

func (c *SpaceClient) Delete(spaceGUID string) error {
	req := c.client.NewRequest("DELETE", "/v3/spaces/"+spaceGUID)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error while deleting space: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	return nil
}

func (c *SpaceClient) Update(spaceGUID string, r resource.UpdateSpaceRequest) (*resource.Space, error) {
	req := c.client.NewRequest("PATCH", "/v3/spaces/"+spaceGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while updating space: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error updating space %s, response code: %d", spaceGUID, resp.StatusCode)
	}

	var space resource.Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, fmt.Errorf("error reading space JSON: %w", err)
	}

	return &space, nil
}

func (c *SpaceClient) ListByQuery(query url.Values) ([]resource.Space, error) {
	var spaces []resource.Space
	requestURL := "/v3/spaces"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.client.NewRequest("GET", requestURL)
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting spaces: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing spaces, response code: %d", resp.StatusCode)
		}

		var data resource.ListSpacesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list spaces: %w", err)
		}

		spaces = append(spaces, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for spaces: %w", err)
		}
	}

	return spaces, nil
}

// ListUsers lists users by space GUID
func (c *SpaceClient) ListUsers(spaceGUID string) ([]resource.User, error) {
	var users []resource.User
	requestURL := "/v3/spaces/" + spaceGUID + "/users"

	for {
		r := c.client.NewRequest("GET", requestURL)
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting space users: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing space users, response code: %d", resp.StatusCode)
		}

		var data resource.ListSpaceUsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list space users: %w", err)
		}
		users = append(users, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for space users: %w", err)
		}
	}

	return users, nil
}
