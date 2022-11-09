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
