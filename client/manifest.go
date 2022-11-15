package client

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client/http"
	"io"
	http2 "net/http"
	"strings"
)

type ManifestClient commonClient

// Generate the specified app manifest as a yaml text string
func (c *ManifestClient) Generate(appGUID string) (string, error) {
	p := path("/v3/apps/%s/manifest", appGUID)
	req := http.NewRequest("GET", p)

	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %w", p, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK {
		return "", c.client.handleError(resp)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response %s: %w", p, err)
	}
	return buf.String(), nil
}

func (c *ManifestClient) ApplyManifest(spaceGUID string, manifest string) (string, error) {
	reader := strings.NewReader(manifest)
	req := http.NewRequest("POST", path("/v3/spaces/%s/actions/apply_manifest", spaceGUID)).
		WithContentType("application/x-yaml").
		WithBody(reader)

	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error uploading manifest %s bits: %w", spaceGUID, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http2.StatusAccepted {
		return "", c.client.handleError(resp)
	}

	jobID, err := c.client.decodeBodyOrJobID(resp, nil)
	if err != nil {
		return "", fmt.Errorf("error reading jobID: %w", err)
	}
	return jobID, nil
}
