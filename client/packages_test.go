package client

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestListPackagesForApp(t *testing.T) {
	setup(MockRoute{"GET", "/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages", []string{listPackagesForAppPayloadPage1, listPackagesForAppPayloadPage2}, "", http.StatusOK, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	packages, err := client.Packages.ListForApp("f2efe391-2b5b-4836-8518-ad93fa9ebf69", nil)
	require.NoError(t, err)
	require.Len(t, packages, 2)

	require.Equal(t, "bits", packages[0].Type)
	require.Equal(t, "READY", string(packages[0].State))
	require.Equal(t, "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69", packages[0].Links["app"].Href)
	require.Equal(t, "https://api.example.org/v3/packages/752edab0-2147-4f58-9c25-cd72ad8c3561/download", packages[0].Links["download"].Href)
	require.Equal(t, "bits", string(packages[1].Type))
	require.Equal(t, "READY", string(packages[1].State))
	require.Equal(t, "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69", packages[1].Links["app"].Href)
	require.Equal(t, "https://api.example.org/v3/packages/2345ab-2147-4f58-9c25-cd72ad8c3561/download", packages[1].Links["download"].Href)

}

func TestCopyPackage(t *testing.T) {
	expectedBody := `{"relationships":{"app":{"data":{"guid":"app-guid"}}}}`
	setup(MockRoute{"POST", "/v3/packages", []string{copyPackagePayload}, "", http.StatusCreated, "source_guid=package-guid", expectedBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	pkg, err := client.Packages.Copy("package-guid", "app-guid")
	require.NoError(t, err)

	require.Equal(t, "COPYING", string(pkg.State))
	require.Equal(t, "docker", pkg.Type)
	require.Equal(t, "fec72fc1-e453-4463-a86d-5df426f337a3", pkg.GUID)

	docker, err := pkg.DockerData()
	require.NoError(t, err)
	require.Equal(t, "http://awesome-sauce.example.org", docker.Image)
}

func TestPackageDataDockerErrorsWhenTypeIsBits(t *testing.T) {
	p := resource.Package{Type: "bits"}
	_, err := p.DockerData()
	require.NotNil(t, err)
}

func TestPackageDataDocker(t *testing.T) {
	p := resource.Package{
		Type: "docker",
		Data: json.RawMessage(`{"image":"nginx","username":"admin","password":"password"}`),
	}
	d, err := p.DockerData()
	require.NoError(t, err)
	require.Equal(t, "nginx", d.Image)
	require.Equal(t, "admin", d.Username)
	require.Equal(t, "password", d.Password)
}

func TestPackageDataBitsErrorsWhenTypeIsDocker(t *testing.T) {
	p := resource.Package{Type: "docker"}
	_, err := p.BitsData()
	require.NotNil(t, err)
}

func TestPackageDataBits(t *testing.T) {
	p := resource.Package{
		Type: "bits",
		Data: json.RawMessage(`{"error":"None","checksum":{"type":"sha256","value":"foo"}}`),
	}
	b, err := p.BitsData()
	require.NoError(t, err)
	require.Equal(t, "None", b.Error)
	require.Equal(t, "sha256", b.Checksum.Type)
}
