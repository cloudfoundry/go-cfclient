package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	"github.com/pkg/errors"
)

func (c *Client) ListPackagesForApp(appGUID string, query url.Values) ([]v3.Package, error) {
	var packages []v3.Package
	requestURL := "/v3/apps/" + appGUID + "/packages"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.DoRequest(c.NewRequest("GET", requestURL))
		if err != nil {
			return nil, errors.Wrapf(err, "Error requesting packages for app %s", appGUID)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 app packages, response code: %d", resp.StatusCode)
		}

		var data v3.ListPackagesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 app packages")
		}

		packages = append(packages, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 packages")
		}
	}
	return packages, nil
}

// CopyPackage makes a copy of a package that is associated with one app
// and associates the copy with a new app.
func (c *Client) CopyPackage(packageGUID, appGUID string) (*v3.Package, error) {
	req := c.NewRequest("POST", "/v3/packages?source_guid="+packageGUID)
	req.obj = map[string]interface{}{
		"relationships": map[string]interface{}{
			"app": v3.ToOneRelationship{
				Data: v3.Relationship{
					GUID: appGUID,
				},
			},
		},
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while copying v3 package")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error copying v3 package %s, response code: %d", packageGUID, resp.StatusCode)
	}

	var pkg v3.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app package")
	}

	return &pkg, nil
}

// CreateDockerPackage creates a Docker package
func (c *Client) CreateDockerPackage(image string, appGUID string, dockerCredentials *v3.DockerCredentials) (*v3.Package, error) {
	req := c.NewRequest("POST", "/v3/packages")
	req.obj = v3.CreateDockerPackageRequest{
		Type: "docker",
		Relationships: map[string]v3.ToOneRelationship{
			"app": {Data: v3.Relationship{GUID: appGUID}},
		},
		Data: v3.DockerPackageData{
			Image:             image,
			DockerCredentials: dockerCredentials,
		},
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while copying v3 package")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 docker package, response code: %d", resp.StatusCode)
	}

	var pkg v3.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 app package")
	}

	return &pkg, nil
}
