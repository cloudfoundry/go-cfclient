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

var t *testing.T
var ctx context.Context
var cf *client.Client

func TestEndToEnd(tt *testing.T) {
	t = tt
	ctx = context.Background()
	cf = createClient()

	// get the org with the access token
	org := getOrg()
	fmt.Printf("Targeting org %s\n", org.Name)

	// try to get the space
	space := getSpace(org)
	fmt.Printf("Targeting space %s\n", space.Name)

	// push an app
	app := pushApp(org, space)
	fmt.Printf("Successfully pushed %s\n", app.Name)

	// curl the app
	curlApp(app)
	fmt.Printf("Successfully curled %s\n", app.Name)

	// download the current app droplet
	dropletFile := downloadDroplet(app)
	fmt.Printf("Successfully downloaded %s\n", dropletFile)

	// create a new droplet
	createNewDroplet(app, dropletFile)
	fmt.Println("Successfully uploaded new droplet")
}

func createNewDroplet(app *resource.App, dropletFile string) {
	c := resource.NewDropletCreate(app.GUID)
	d, err := cf.Droplets.Create(ctx, c)
	require.NoError(t, err)

	tgzFile, err := os.Open(dropletFile)
	require.NoError(t, err)
	defer tgzFile.Close()

	jobGUID, _, err := cf.Droplets.Upload(ctx, d.GUID, tgzFile)
	require.NoError(t, err)

	err = cf.Jobs.PollComplete(ctx, jobGUID, nil)
	require.NoError(t, err)
}

func downloadDroplet(app *resource.App) string {
	opts := client.NewDropletAppListOptions()
	opts.Current = true
	droplets, err := cf.Droplets.ListForAppAll(ctx, app.GUID, opts)
	require.NoError(t, err)
	require.Equal(t, 1, len(droplets))
	droplet := droplets[0]

	r, err := cf.Droplets.Download(ctx, droplet.GUID)
	require.NoError(t, err)

	f, err := os.CreateTemp("", "droplet-*.tgz")
	require.NoError(t, err)
	defer f.Close()

	_, err = io.Copy(f, r)
	require.NoError(t, err)

	return f.Name()
}

func curlApp(app *resource.App) {
	routes, err := cf.Routes.ListForAppAll(ctx, app.GUID, nil)
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

func pushApp(org *resource.Organization, space *resource.Space) *resource.App {
	appRoute := getAppRoute(org)

	zipPath := createZipFile()
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

	fmt.Printf("Pushing app %s...\n", AppName)
	p := operation.NewAppPushOperation(cf, org.Name, space.Name)
	app, err := p.Push(ctx, manifest, zipReader)
	require.NoError(t, err)
	return app
}

func getAppRoute(org *resource.Organization) string {
	domains, err := cf.Domains.ListForOrganizationAll(ctx, org.GUID, nil)
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

func createZipFile() string {
	fmt.Println("Creating zip archive")
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

func getOrg() *resource.Organization {
	opts := client.NewOrganizationListOptions()
	opts.Names = client.Filter{
		Values: []string{OrgName},
	}
	orgs, _, err := cf.Organizations.List(ctx, opts)
	require.NoError(t, err)

	var org *resource.Organization
	if len(orgs) > 0 {
		org = orgs[0]
	} else {
		oc := &resource.OrganizationCreate{
			Name: OrgName,
		}
		org, err = cf.Organizations.Create(ctx, oc)
		require.NoError(t, err)
	}
	require.Equal(t, OrgName, org.Name)
	require.NotEmpty(t, org.GUID)
	return org
}

func getSpace(org *resource.Organization) *resource.Space {
	opts := client.NewSpaceListOptions()
	opts.Names = client.Filter{
		Values: []string{SpaceName},
	}
	spaces, _, err := cf.Spaces.List(ctx, opts)
	require.NoError(t, err)

	var space *resource.Space
	if len(spaces) > 0 {
		space = spaces[0]
	} else {
		sc := resource.NewSpaceCreate(SpaceName, org.GUID)
		space, err = cf.Spaces.Create(ctx, sc)
		require.NoError(t, err)
	}
	require.Equal(t, SpaceName, space.Name)
	require.NotEmpty(t, space.GUID)
	return space
}

func createClient() *client.Client {
	cfg, err := config.NewFromCFHome()
	require.NoError(t, err)
	cfg.WithSkipTLSValidation(true)
	c, err := client.New(cfg)
	require.NoError(t, err)
	return c
}
