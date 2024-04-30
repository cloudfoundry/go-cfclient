package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	internalhttp "github.com/cloudfoundry/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry/go-cfclient/v3/internal/ios"
	"github.com/cloudfoundry/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type ManifestClient commonClient

// Generate the specified app manifest as a yaml text string
func (c *ManifestClient) Generate(ctx context.Context, appGUID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.client.ApiURL(path.Format("/v3/apps/%s/manifest", appGUID)), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create manifest request for app %s: %w", appGUID, err)
	}

	resp, err := c.client.ExecuteAuthRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute manifest request for app %s: %w", appGUID, err)
	}
	defer ios.Close(resp.Body)

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		return "", fmt.Errorf("failed to read manifest for app %s: %w", appGUID, err)
	}
	return buf.String(), nil
}

// ApplyManifest applies the changes specified in a manifest to the named apps and their underlying processes
// asynchronously and returns a jobGUID.
//
// The apps must reside in the space. These changes are additive and will not modify any unspecified
// properties or remove any existing environment variables, routes, or services.
func (c *ManifestClient) ApplyManifest(ctx context.Context, spaceGUID string, manifest string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.client.ApiURL(path.Format("/v3/spaces/%s/actions/apply_manifest", spaceGUID)), strings.NewReader(manifest))
	if err != nil {
		return "", fmt.Errorf("failed to create manifest apply request for space %s: %w", spaceGUID, err)
	}
	req.Header.Set("Content-Type", "application/x-yaml")

	resp, err := c.client.ExecuteAuthRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload manifest for space %s: %w", spaceGUID, err)
	}
	defer ios.Close(resp.Body)
	return internalhttp.DecodeJobID(resp), nil
}

// ManifestDiff compares the provided manifest against the current state of the space.
func (c *ManifestClient) ManifestDiff(ctx context.Context, spaceGUID string, manifest string) (*resource.ManifestDiff, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.client.ApiURL(path.Format("/v3/spaces/%s/manifest_diff", spaceGUID)), strings.NewReader(manifest))
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest diff request for space %s: %w", spaceGUID, err)
	}
	req.Header.Set("Content-Type", "application/x-yaml")

	resp, err := c.client.ExecuteAuthRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute manifest diff request for space %s: %w", spaceGUID, err)
	}
	defer ios.Close(resp.Body)
	var diff resource.ManifestDiff
	if err = internalhttp.DecodeBody(resp, &diff); err != nil {
		return nil, err
	}
	return &diff, nil
}
