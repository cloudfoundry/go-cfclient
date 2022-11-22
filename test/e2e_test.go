//go:build integration
// +build integration

package test

import (
	"archive/zip"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/operation"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	OrgName   = "go-cfclient-e2e"
	SpaceName = "go-cfclient-e2e"
	AppName   = "go-cfclient-hello-world"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	c := createClient(t)

	// get the org with the access token
	org := getOrg(t, ctx, c)
	fmt.Printf("Targeting org %s\n", org.Name)

	// try to get the space
	space := getSpace(t, ctx, c, org)
	fmt.Printf("Targeting space %s\n", space.Name)

	// push an app
	app := pushApp(t, ctx, c, org, space)
	fmt.Printf("Successfully pushed %s\n", app.Name)

	// curl the app
	curlApp(t, ctx, c, app)
}

func curlApp(t *testing.T, ctx context.Context, c *client.Client, app *resource.App) {
	routes, err := c.Routes.ListForAppAll(ctx, app.GUID, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(routes), 1, "expected 1+ routes")
	route := routes[0]

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	appClient := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/go-cfclient-e2e-test", route.URL), nil)
	require.NoError(t, err)
	resp, err := appClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, 200, resp.StatusCode)
}

func pushApp(t *testing.T, ctx context.Context, c *client.Client, org *resource.Organization, space *resource.Space) *resource.App {
	appRoute := getAppRoute(t, ctx, c, org)

	zipPath := createZipFile(t)
	zipReader, err := os.Open(zipPath)
	require.NoError(t, err)

	manifest := operation.NewAppManifest("go-cfclient-hello-world")
	manifest.Buildpacks = []string{
		"go_buildpack",
	}
	manifest.Routes = []operation.AppManifestRoutes{
		{
			Route: appRoute,
		},
	}
	manifest.HealthCheckType = "http"
	manifest.Memory = "64M"

	p := operation.NewAppPushOperation(c, org.Name, space.Name)
	app, err := p.Push(ctx, manifest, zipReader)
	require.NoError(t, err)
	return app
}

func getAppRoute(t *testing.T, ctx context.Context, c *client.Client, org *resource.Organization) string {
	domains, err := c.Domains.ListForOrgAll(ctx, org.GUID, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(domains), 1, "expected 1+ available domains for org %s", org.Name)

	domain := ""
	for _, d := range domains {
		if !strings.Contains(d.Name, "internal") {
			domain = d.Name
			break
		}
	}
	require.NotEmpty(t, domain, "expected to find 1 non-internal domain")
	return fmt.Sprintf("%s.%s", AppName, domain)
}

func createZipFile(t *testing.T) string {
	fmt.Println("creating zip archive...")
	archive, err := os.CreateTemp("", "go-cfclient-hello-world-*.zip")
	require.NoError(t, err)
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	err = writeFileToZip(zipWriter, "helloworld/go.mod")
	require.NoError(t, err)
	err = writeFileToZip(zipWriter, "helloworld/main.go")
	require.NoError(t, err)

	return archive.Name()
}

func writeFileToZip(zipWriter *zip.Writer, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zipWriter.Create(path.Base(filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}

func getOrg(t *testing.T, ctx context.Context, c *client.Client) *resource.Organization {
	opts := client.NewOrgListOptions()
	opts.Names = client.Filter{
		Values: []string{OrgName},
	}
	orgs, _, err := c.Organizations.List(ctx, opts)
	require.NoError(t, err)

	var org *resource.Organization
	if len(orgs) > 0 {
		org = orgs[0]
	} else {
		oc := &resource.OrganizationCreate{
			Name: OrgName,
		}
		org, err = c.Organizations.Create(ctx, oc)
		require.NoError(t, err)
	}
	require.Equal(t, OrgName, org.Name)
	require.NotEmpty(t, org.GUID)
	return org
}

func getSpace(t *testing.T, ctx context.Context, c *client.Client, org *resource.Organization) *resource.Space {
	opts := client.NewSpaceListOptions()
	opts.Names = client.Filter{
		Values: []string{SpaceName},
	}
	spaces, _, err := c.Spaces.List(ctx, opts)
	require.NoError(t, err)

	var space *resource.Space
	if len(spaces) > 0 {
		space = spaces[0]
	} else {
		sc := resource.NewSpaceCreate(SpaceName, org.GUID)
		space, err = c.Spaces.Create(ctx, sc)
		require.NoError(t, err)
	}
	require.Equal(t, SpaceName, space.Name)
	require.NotEmpty(t, space.GUID)
	return space
}

func createClient(t *testing.T) *client.Client {
	cfg, err := config.NewFromCFHome()
	require.NoError(t, err)
	cfg.WithSkipTLSValidation(true)
	c, err := client.New(cfg)
	require.NoError(t, err)
	return c
}
