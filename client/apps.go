package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type AppClient commonClient

func (c *AppClient) Create(r resource.CreateAppRequest) (*resource.App, error) {
	req := c.client.NewRequest("POST", "/v3/apps")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"space": resource.ToOneRelationship{
				Data: resource.Relationship{
					GUID: r.SpaceGUID,
				},
			},
		},
	}
	if len(r.EnvironmentVariables) > 0 {
		params["environment_variables"] = r.EnvironmentVariables
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while creating app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating app %s, response code: %d", r.Name, resp.StatusCode)
	}

	var app resource.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("error reading app JSON: %w", err)
	}

	return &app, nil
}

func (c *AppClient) Delete(guid string) error {
	req := c.client.NewRequest("DELETE", "/v3/apps/"+guid)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("error while deleting app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error deleting app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	return nil
}

func (c *AppClient) Get(guid string) (*resource.App, error) {
	req := c.client.NewRequest("GET", "/v3/apps/"+guid)

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while getting app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app resource.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("error reading app JSON: %w", err)
	}

	return &app, nil
}

func (c *AppClient) GetEnvironment(appGUID string) (resource.AppEnvironment, error) {
	var result resource.AppEnvironment

	resp, err := c.client.DoRequest(c.client.NewRequest("GET", "/v3/apps/"+appGUID+"/env"))
	if err != nil {
		return result, fmt.Errorf("error requesting app env for %s: %w", appGUID, err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("error parsing JSON for app env: %w", err)
	}

	return result, nil
}

func (c *AppClient) List() ([]resource.App, error) {
	return c.ListByQuery(url.Values{})
}

func (c *AppClient) ListByQuery(query url.Values) ([]resource.App, error) {
	var apps []resource.App
	requestURL := "/v3/apps"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.client.NewRequest("GET", requestURL)
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting  apps: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing apps, response code: %d", resp.StatusCode)
		}

		var data resource.ListAppsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list apps: %w", err)
		}

		apps = append(apps, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for apps: %w", err)
		}
	}

	return apps, nil
}

func (c *AppClient) SetEnvVariables(appGUID string, envRequest resource.EnvVar) (resource.EnvVar, error) {
	var result resource.EnvVarResponse

	req := c.client.NewRequest("PATCH", "/v3/apps/"+appGUID+"/environment_variables")
	req.obj = envRequest

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return result.EnvVar, fmt.Errorf("error setting app env variables for %s: %w", appGUID, err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result.EnvVar, fmt.Errorf("error parsing JSON for app env: %w", err)
	}

	return result.EnvVar, nil
}

func (c *AppClient) Start(guid string) (*resource.App, error) {
	req := c.client.NewRequest("POST", "/v3/apps/"+guid+"/actions/start")
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while starting app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error starting app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app resource.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("error reading app JSON: %w", err)
	}

	return &app, nil
}

func (c *AppClient) Update(appGUID string, r resource.UpdateAppRequest) (*resource.App, error) {
	req := c.client.NewRequest("PATCH", "/v3/apps/"+appGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Lifecycle != nil {
		params["lifecycle"] = r.Lifecycle
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while updating app: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error updating app %s, response code: %d", appGUID, resp.StatusCode)
	}

	var app resource.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("error reading app JSON: %w", err)
	}

	return &app, nil
}
