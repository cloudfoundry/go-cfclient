package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type PackageClient commonClient

func (c *PackageClient) ListForApp(appGUID string, query url.Values) ([]resource.Package, error) {
	var packages []resource.Package
	requestURL := "/v3/apps/" + appGUID + "/packages"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.client.DoRequest(c.client.NewRequest("GET", requestURL))
		if err != nil {
			return nil, fmt.Errorf("error requesting packages for app %s: %w", appGUID, err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing v3 app packages, response code: %d", resp.StatusCode)
		}

		var data resource.ListPackagesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list v3 app packages: %w", err)
		}

		packages = append(packages, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for v3 packages: %w", err)
		}
	}
	return packages, nil
}

// Copy makes a copy of a package that is associated with one app
// and associates the copy with a new app.
func (c *PackageClient) Copy(packageGUID, appGUID string) (*resource.Package, error) {
	req := c.client.NewRequest("POST", "/v3/packages?source_guid="+packageGUID)
	req.obj = map[string]interface{}{
		"relationships": map[string]interface{}{
			"app": resource.ToOneRelationship{
				Data: resource.Relationship{
					GUID: appGUID,
				},
			},
		},
	}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while copying v3 package: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error copying v3 package %s, response code: %d", packageGUID, resp.StatusCode)
	}

	var pkg resource.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("error reading v3 app package: %w", err)
	}

	return &pkg, nil
}

// CreateDocker creates a Docker package
func (c *PackageClient) CreateDocker(image string, appGUID string, dockerCredentials *resource.DockerCredentials) (*resource.Package, error) {
	req := c.client.NewRequest("POST", "/v3/packages")
	req.obj = resource.CreateDockerPackageRequest{
		Type: "docker",
		Relationships: map[string]resource.ToOneRelationship{
			"app": {Data: resource.Relationship{GUID: appGUID}},
		},
		Data: resource.DockerPackageData{
			Image:             image,
			DockerCredentials: dockerCredentials,
		},
	}

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error while copying v3 package: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 docker package, response code: %d", resp.StatusCode)
	}

	var pkg resource.Package
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("error reading v3 app package: %w", err)
	}

	return &pkg, nil
}
