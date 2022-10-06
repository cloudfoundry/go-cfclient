package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestCreateBuild(t *testing.T) {
	body := `{"metadata":{"labels":{"foo":"bar"}},"package":{"guid":"package-guid"}}`
	setup(MockRoute{"POST", "/v3/builds", []string{createBuildPayload}, "", http.StatusCreated, "", &body}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	build, err := client.CreateBuild("package-guid", nil,
		&resource.Metadata{Labels: map[string]string{"foo": "bar"}})
	require.NoError(t, err)
	require.NotNil(t, build)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", build.GUID)
	require.Equal(t, "bill", build.CreatedBy.Name)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
}
