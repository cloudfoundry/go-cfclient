package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestCreateBuild(t *testing.T) {
	body := `{"metadata":{"labels":{"foo":"bar"}},"package":{"guid":"package-guid"}}`
	setup(MockRoute{"POST", "/v3/builds", []string{buildPayload}, "", http.StatusCreated, "", &body}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	bc := resource.NewBuildCreate("package-guid")
	bc.Metadata = &resource.Metadata{
		Labels: map[string]string{
			"foo": "bar",
		},
	}
	build, err := client.Builds.Create(bc)
	require.NoError(t, err)
	require.NotNil(t, build)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", build.GUID)
	require.Equal(t, "bill", build.CreatedBy.Name)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
}

func TestGetBuild(t *testing.T) {
	setup(MockRoute{"GET", "/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16", []string{buildPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	build, err := client.Builds.Get("585bc3c1-3743-497d-88b0-403ad6b56d16")
	require.NoError(t, err)
	require.NotNil(t, build)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", build.GUID)
	require.Equal(t, "bill", build.CreatedBy.Name)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
	require.Equal(t, resource.BuildStateStaging, build.State)
	require.Nil(t, build.Error)
	require.Equal(t, "buildpack", build.Lifecycle.Type)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
	require.Nil(t, build.Droplet)
	require.Equal(t, "https://api.example.org/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16", build.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396", build.Links["app"].Href)
}

func TestDeleteBuild(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Builds.Delete("585bc3c1-3743-497d-88b0-403ad6b56d16")
	require.NoError(t, err)
}

func TestUpdateBuild(t *testing.T) {
	setup(MockRoute{"PATCH", "/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16", []string{buildPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	u := resource.NewBuildUpdate()
	u.Metadata.Annotations["foo"] = "bar"
	u.Metadata.Labels["env"] = "dev"
	build, err := client.Builds.Update("585bc3c1-3743-497d-88b0-403ad6b56d16", u)
	require.NoError(t, err)
	require.NotNil(t, build)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", build.GUID)
	require.Equal(t, "bill", build.CreatedBy.Name)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
	require.Equal(t, resource.BuildStateStaging, build.State)
	require.Nil(t, build.Error)
	require.Equal(t, "buildpack", build.Lifecycle.Type)
	require.Equal(t, "8e4da443-f255-499c-8b47-b3729b5b7432", build.Package.GUID)
	require.Nil(t, build.Droplet)
	require.Equal(t, "https://api.example.org/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16", build.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396", build.Links["app"].Href)
}

func TestListBuilds(t *testing.T) {
	setup(MockRoute{"GET", "/v3/builds", []string{buildListPayloadPage1}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	opts := NewBuildListOptions()
	opts.PerPage = 1
	builds, _, err := client.Builds.List(opts)
	require.NoError(t, err)
	require.Len(t, builds, 1)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", builds[0].GUID)
	require.Equal(t, resource.BuildStateStaging, builds[0].State)
	require.Equal(t, "ruby_buildpack", builds[0].Lifecycle.BuildpackData.Buildpacks[0])
}

func TestListAllBuilds(t *testing.T) {
	mr := MockRoute{
		"GET",
		"/v3/builds",
		[]string{buildListPayloadPage1, buildListPayloadPage2},
		"",
		http.StatusOK,
		"",
		nil}
	setup(mr, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	builds, err := client.Builds.ListAll()
	require.NoError(t, err)

	require.Len(t, builds, 2)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", builds[0].GUID)
	require.Equal(t, "787bc3c1-3743-497d-88b0-403ad6b56d23", builds[1].GUID)
}

func TestListForAppBuilds(t *testing.T) {
	mr := MockRoute{
		"GET",
		"/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/builds",
		[]string{buildListPayloadPage1},
		"",
		http.StatusOK,
		"",
		nil}
	setup(mr, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	opts := NewBuildAppListOptions()
	opts.PerPage = 1
	builds, _, err := client.Builds.ListForApp("1cb006ee-fb05-47e1-b541-c34179ddc446", opts)
	require.NoError(t, err)

	require.Len(t, builds, 1)
	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", builds[0].GUID)
}
