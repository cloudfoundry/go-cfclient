package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ManifestClient commonClient

// Generate the specified app manifest as a yaml text string
func (c *ManifestClient) Generate(appGUID string) (string, error) {
	p := path("/v3/apps/%s/manifest", appGUID)
	req := c.client.NewRequest("GET", p)

	resp, err := c.client.DoRequest(req)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %w", p, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting %s, response code: %d", p, resp.StatusCode)
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
	r := c.client.NewRequestWithBody("POST", path("/v3/spaces/%s/actions/apply_manifest", spaceGUID), reader)
	req, err := r.toHTTP()
	if err != nil {
		return "", fmt.Errorf("error posting manifest %s bits: %w", spaceGUID, err)
	}
	req.Header.Set("Content-Type", "application/x-yaml")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error uploading manifest %s bits: %w", spaceGUID, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("error uploading manifest %s bits, response code: %d", spaceGUID, resp.StatusCode)
	}

	jobID, err := c.client.decodeBodyOrJobID(resp, nil)
	if err != nil {
		return "", fmt.Errorf("error reading jobID: %w", err)
	}
	return jobID, nil
}
