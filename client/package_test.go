package client

import (
	"context"
	"encoding/json"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBitsMarshalling(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	rawPkg := g.Package("PROCESSING_UPLOAD").JSON

	var pkg resource.Package
	err := json.Unmarshal([]byte(rawPkg), &pkg)
	require.NoError(t, err)
	require.Nil(t, pkg.Data.Docker)
	require.NotNil(t, pkg.Data.Bits)
	require.Equal(t, "sha256", pkg.Data.Bits.Checksum.Type)
	require.Nil(t, pkg.Data.Bits.Error)

	b, err := json.Marshal(&pkg)
	require.NoError(t, err)
	require.JSONEq(t, rawPkg, string(b))
}

func TestDockerMarshalling(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	rawPkg := g.PackageDocker().JSON

	var pkg resource.Package
	err := json.Unmarshal([]byte(rawPkg), &pkg)
	require.NoError(t, err)
	require.NotNil(t, pkg.Data.Docker)
	require.Nil(t, pkg.Data.Bits)
	require.Equal(t, "registry/image:latest", pkg.Data.Docker.Image)
	require.Equal(t, "username", pkg.Data.Docker.Username)
	require.Equal(t, "secret", pkg.Data.Docker.Password)

	b, err := json.Marshal(&pkg)
	require.NoError(t, err)
	require.JSONEq(t, rawPkg, string(b))
}

func TestPackages(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	pkg := g.Package("PROCESSING_UPLOAD").JSON
	pkg2 := g.Package("PROCESSING_UPLOAD").JSON
	pkg3 := g.Package("PROCESSING_UPLOAD").JSON
	pkg4 := g.Package("PROCESSING_UPLOAD").JSON

	blobstoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("package bits..."))
	}))
	defer blobstoreServer.Close()

	tests := []RouteTest{
		{
			Description: "Create package",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/packages",
				Output:   g.Single(pkg),
				Status:   http.StatusCreated,
				PostForm: `{ "type": "bits", "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewPackageCreate("8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
				return c.Packages.Create(context.Background(), r)
			},
		},
		{
			Description: "Copy package",
			Route: testutil.MockRoute{
				Method:      "POST",
				Endpoint:    "/v3/packages",
				QueryString: "source_guid=66e89f29-475e-4baf-9675-40c6096c017b",
				Output:      g.Single(pkg),
				Status:      http.StatusCreated,
				PostForm:    `{ "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Copy(context.Background(), "66e89f29-475e-4baf-9675-40c6096c017b", "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
			},
		},
		{
			Description: "Delete package",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Delete(context.Background(), "66e89f29-475e-4baf-9675-40c6096c017b")
			},
		},
		{
			Description: "Download package",
			Route: testutil.MockRoute{
				Method:           "GET",
				Endpoint:         "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b/download",
				Status:           http.StatusFound,
				RedirectLocation: blobstoreServer.URL,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				reader, err := c.Packages.Download(context.Background(), "66e89f29-475e-4baf-9675-40c6096c017b")
				require.NoError(t, err)
				buf := new(strings.Builder)
				_, err = io.Copy(buf, reader)
				require.NoError(t, err)
				require.Equal(t, "package bits...", buf.String())
				return nil, nil
			},
		},
		{
			Description: "Get package",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b",
				Output:   g.Single(pkg),
				Status:   http.StatusOK,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Get(context.Background(), "66e89f29-475e-4baf-9675-40c6096c017b")
			},
		},
		{
			Description: "List first page of packages",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages",
				Output:   g.SinglePaged(pkg),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Packages.List(context.Background(), nil)
				return ds, err
			},
		},
		{
			Description: "List all packages",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages",
				Output:   g.Paged([]string{pkg, pkg2}, []string{pkg3, pkg4}),
				Status:   http.StatusOK},
			Expected: g.Array(pkg, pkg2, pkg3, pkg4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List first page of packages for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b/packages",
				Output:   g.SinglePaged(pkg),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Packages.ListForApp(context.Background(), "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", nil)
				return ds, err
			},
		},
		{
			Description: "List all packages for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b/packages",
				Output:   g.Paged([]string{pkg, pkg2}, []string{pkg3}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg, pkg2, pkg3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.ListForAppAll(context.Background(), "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", nil)
			},
		},
		{
			Description: "Update package",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/packages/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b",
				Output:   g.Single(pkg),
				Status:   http.StatusOK,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Update(context.Background(), "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", &resource.PackageUpdate{})
			},
		},
		{
			Description: "Upload package",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/packages/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b/upload",
				Output:   g.Single(pkg),
				Status:   http.StatusOK,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				zipFile := strings.NewReader("package")
				return c.Packages.Upload(context.Background(), "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", zipFile)
			},
		},
	}
	ExecuteTests(tests, t)
}
