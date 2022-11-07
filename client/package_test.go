package client

import (
	"encoding/json"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestBitsMarshalling(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	rawPkg := g.Package()

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
	g := test.NewObjectJSONGenerator(1)
	rawPkg := g.PackageDocker()

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
	g := test.NewObjectJSONGenerator(1)
	pkg := g.Package()
	pkg2 := g.Package()
	pkg3 := g.Package()
	pkg4 := g.Package()

	tests := []RouteTest{
		{
			Description: "Create package",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/packages",
				Output:   []string{pkg},
				Status:   http.StatusCreated,
				PostForm: `{ "type": "bits", "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewPackageCreate("8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
				return c.Packages.Create(r)
			},
		},
		{
			Description: "Copy package",
			Route: MockRoute{
				Method:      "POST",
				Endpoint:    "/v3/packages",
				QueryString: "source_guid=66e89f29-475e-4baf-9675-40c6096c017b",
				Output:      []string{pkg},
				Status:      http.StatusCreated,
				PostForm:    `{ "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Copy("66e89f29-475e-4baf-9675-40c6096c017b", "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
			},
		},
		{
			Description: "Delete package",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Packages.Delete("66e89f29-475e-4baf-9675-40c6096c017b")
			},
		},
		{
			Description: "Download package",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b/download",
				Output:   []string{"package bits..."},
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				reader, err := c.Packages.Download("66e89f29-475e-4baf-9675-40c6096c017b")
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
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/66e89f29-475e-4baf-9675-40c6096c017b",
				Output:   []string{pkg},
				Status:   http.StatusOK,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Get("66e89f29-475e-4baf-9675-40c6096c017b")
			},
		},
		{
			Description: "List first page of packages",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages",
				Output:   g.Paged([]string{pkg}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Packages.List(nil)
				return ds, err
			},
		},
		{
			Description: "List all packages",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages",
				Output:   g.Paged([]string{pkg, pkg2}, []string{pkg3, pkg4}),
				Status:   http.StatusOK},
			Expected: g.Array(pkg, pkg2, pkg3, pkg4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.ListAll(nil)
			},
		},
		{
			Description: "List first page of packages for an app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b/packages",
				Output:   g.Paged([]string{pkg}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Packages.ListForApp("8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", nil)
				return ds, err
			},
		},
		{
			Description: "List all packages for an app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b/packages",
				Output:   g.Paged([]string{pkg, pkg2}, []string{pkg3}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(pkg, pkg2, pkg3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.ListForAppAll("8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", nil)
			},
		},
		{
			Description: "Update package",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/packages/8d1f1d2e-08b1-4a10-a8df-471a1418cb8b",
				Output:   []string{pkg},
				Status:   http.StatusOK,
			},
			Expected: pkg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Packages.Update("8d1f1d2e-08b1-4a10-a8df-471a1418cb8b", &resource.PackageUpdate{})
			},
		},
	}
	executeTests(tests, t)
}