package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) CreateApp(r v3.CreateAppRequest) (*v3.App, error) {
	req := c.NewRequest("POST", "/v3/apps")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"space": v3.ToOneRelationship{
				Data: v3.Relationship{
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
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating  app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating  app %s, response code: %d", r.Name, resp.StatusCode)
	}

	var app v3.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading  app JSON")
	}

	return &app, nil
}

func (c *Client) GetAppByGUID(guid string) (*v3.App, error) {
	req := c.NewRequest("GET", "/v3/apps/"+guid)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting  app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting  app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app v3.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading  app JSON")
	}

	return &app, nil
}

func (c *Client) StartApp(guid string) (*v3.App, error) {
	req := c.NewRequest("POST", "/v3/apps/"+guid+"/actions/start")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while starting  app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error starting  app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	var app v3.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading  app JSON")
	}

	return &app, nil
}

func (c *Client) DeleteApp(guid string) error {
	req := c.NewRequest("DELETE", "/v3/apps/"+guid)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting  app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting  app with GUID [%s], response code: %d", guid, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateApp(appGUID string, r v3.UpdateAppRequest) (*v3.App, error) {
	req := c.NewRequest("PATCH", "/v3/apps/"+appGUID)
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

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating  app")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating  app %s, response code: %d", appGUID, resp.StatusCode)
	}

	var app v3.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, errors.Wrap(err, "Error reading  app JSON")
	}

	return &app, nil
}

func (c *Client) ListApps() ([]v3.App, error) {
	return c.ListAppsByQuery(url.Values{})
}

func (c *Client) ListAppsByQuery(query url.Values) ([]v3.App, error) {
	var apps []v3.App
	requestURL := "/v3/apps"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  apps")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  apps, response code: %d", resp.StatusCode)
		}

		var data v3.ListAppsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  apps")
		}

		apps = append(apps, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for  apps")
		}
	}

	return apps, nil
}

func extractPathFromURL(requestURL string) (string, error) {
	url, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	result := url.Path
	if q := url.Query().Encode(); q != "" {
		result = result + "?" + q
	}
	return result, nil
}

func (c *Client) GetAppEnvironment(appGUID string) (v3.AppEnvironment, error) {
	var result v3.AppEnvironment

	resp, err := c.DoRequest(c.NewRequest("GET", "/v3/apps/"+appGUID+"/env"))
	if err != nil {
		return result, errors.Wrapf(err, "Error requesting app env for %s", appGUID)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, errors.Wrap(err, "Error parsing JSON for app env")
	}

	return result, nil
}

func (c *Client) SetAppEnvVariables(appGUID string, envRequest v3.EnvVar) (v3.EnvVar, error) {
	var result v3.EnvVarResponse

	req := c.NewRequest("PATCH", "/v3/apps/"+appGUID+"/environment_variables")
	req.obj = envRequest

	resp, err := c.DoRequest(req)
	if err != nil {
		return result.EnvVar, errors.Wrapf(err, "Error setting app env variables for %s", appGUID)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result.EnvVar, errors.Wrap(err, "Error parsing JSON for app env")
	}

	return result.EnvVar, nil
}
