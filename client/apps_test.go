package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreateApp(t *testing.T) {
	expectedBody := `{"environment_variables":{"FOO":"BAR"},"name":"my-app","relationships":{"space":{"data":{"guid":"space-guid"}}}}`
	setup(MockRoute{"POST", "/v3/apps", []string{createAppPayload}, "", http.StatusCreated, "", &expectedBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	r := resource.NewAppCreate("my-app", "space-guid")
	r.EnvironmentVariables = map[string]string{"FOO": "BAR"}
	app, err := client.Applications.Create(r)
	require.NoError(t, err)
	require.NotNil(t, app)

	require.Equal(t, "app-guid", app.GUID)
	require.Equal(t, "space-guid", app.Relationships["space"].Data.GUID)
	require.Equal(t, "buildpack", app.Lifecycle.Type)
	require.Len(t, app.Lifecycle.BuildpackData.Buildpacks, 1)
	require.Equal(t, "java_buildpack", app.Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "cflinuxfs2", app.Lifecycle.BuildpackData.Stack)
	require.Equal(t, "https://api.example.org/v3/spaces/space-guid", app.Links["space"].Href)
	require.Len(t, app.Metadata.Annotations, 0)
}

func TestGetApp(t *testing.T) {
	setup(MockRoute{"GET", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{getAppPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	app, err := client.Applications.Get("1cb006ee-fb05-47e1-b541-c34179ddc446")
	require.NoError(t, err)
	require.NotNil(t, app)

	require.Equal(t, "1cb006ee-fb05-47e1-b541-c34179ddc446", app.GUID)
	require.Equal(t, "my_app", app.Name)
	require.Equal(t, "buildpack", app.Lifecycle.Type)
	require.Len(t, app.Lifecycle.BuildpackData.Buildpacks, 1)
	require.Equal(t, "java_buildpack", app.Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "cflinuxfs2", app.Lifecycle.BuildpackData.Stack)
	require.Equal(t, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576", app.Links["space"].Href)
	require.Len(t, app.Metadata.Annotations, 1)
	require.Equal(t, "Bill tel(1111111) email(bill@fixme), Bob tel(222222) pager(3333333#555) email(bob@fixme)", app.Metadata.Annotations["contacts"])
}

func TestGetAppEnv(t *testing.T) {
	setup(MockRoute{"GET", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/env", []string{getAppEnvPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	env, err := client.Applications.GetEnvironment("1cb006ee-fb05-47e1-b541-c34179ddc446")
	require.NoError(t, err)

	require.Len(t, env.EnvVars, 1)
	require.Equal(t, "production", env.EnvVars["RAILS_ENV"])

	require.Len(t, env.StagingEnv, 1)
	require.Equal(t, "http://gem-cache.example.org", env.StagingEnv["GEM_CACHE"])

	require.Len(t, env.RunningEnv, 1)
	require.Equal(t, "http://proxy.example.org", env.RunningEnv["HTTP_PROXY"])

	require.Len(t, env.SystemEnvVars, 1)
	require.Contains(t, env.SystemEnvVars, "VCAP_SERVICES")
	require.Len(t, env.AppEnvVars, 1)
	require.Contains(t, env.AppEnvVars, "VCAP_APPLICATION")
}

func TestSetAppEnvVariables(t *testing.T) {
	setup(MockRoute{"PATCH", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables", []string{setAppEnvironmentVariablesPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	falseVar := "false"
	env, err := client.Applications.SetEnvVariables("1cb006ee-fb05-47e1-b541-c34179ddc446",
		resource.EnvVar{Var: map[string]*string{
			"DEBUG": &falseVar,
			"USER":  nil,
		}},
	)
	require.NoError(t, err)
	require.Len(t, env.Var, 2)
	require.Equal(t, "production", *env.Var["RAILS_ENV"])
	require.Equal(t, "false", *env.Var["DEBUG"])
}

func TestStartApp(t *testing.T) {
	setup(MockRoute{"POST", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start", []string{startAppPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	app, err := client.Applications.Start("1cb006ee-fb05-47e1-b541-c34179ddc446")
	require.NoError(t, err)
	require.NotNil(t, app)

	require.Equal(t, "STARTED", app.State)
	require.Equal(t, "1cb006ee-fb05-47e1-b541-c34179ddc446", app.GUID)
	require.Equal(t, "my_app", app.Name)
	require.Equal(t, "buildpack", app.Lifecycle.Type)
	require.Len(t, app.Lifecycle.BuildpackData.Buildpacks, 1)
	require.Equal(t, "java_buildpack", app.Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "cflinuxfs2", app.Lifecycle.BuildpackData.Stack)
	require.Equal(t, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576", app.Links["space"].Href)
	require.Len(t, app.Metadata.Annotations, 0)
}

func TestDeleteApp(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Applications.Delete("1cb006ee-fb05-47e1-b541-c34179ddc446")
	require.NoError(t, err)
}

func TestUpdateApp(t *testing.T) {
	setup(MockRoute{"PATCH", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{updateAppPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	app, err := client.Applications.Update("1cb006ee-fb05-47e1-b541-c34179ddc446", &resource.AppUpdate{})
	require.NoError(t, err)
	require.NotNil(t, app)

	require.Equal(t, "1cb006ee-fb05-47e1-b541-c34179ddc446", app.GUID)
	require.Equal(t, "STARTED", app.State)
	require.Equal(t, "my_app", app.Name)
	require.Equal(t, "buildpack", app.Lifecycle.Type)
	require.Len(t, app.Lifecycle.BuildpackData.Buildpacks, 1)
	require.Equal(t, "java_buildpack", app.Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "cflinuxfs2", app.Lifecycle.BuildpackData.Stack)
	require.Equal(t, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576", app.Links["space"].Href)
	require.Len(t, app.Metadata.Annotations, 0)
	require.Len(t, app.Metadata.Labels, 2)
	require.Equal(t, "production", app.Metadata.Labels["environment"])
	require.Equal(t, "false", app.Metadata.Labels["internet-facing"])
}

func TestListApps(t *testing.T) {
	setup(MockRoute{"GET", "/v3/apps", []string{listAppsPayloadPage1}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	opts := NewAppListOptions()
	apps, _, err := client.Applications.List(opts)
	require.NoError(t, err)
	require.Len(t, apps, 1)

	require.Equal(t, "1cb006ee-fb05-47e1-b541-c34179ddc446", apps[0].GUID)
	require.Equal(t, "my_app", apps[0].Name)
	require.Equal(t, "java_buildpack", apps[0].Lifecycle.BuildpackData.Buildpacks[0])
}

func TestListAllApps(t *testing.T) {
	mr := MockRoute{
		"GET",
		"/v3/apps",
		[]string{listAppsPayloadPage1, listAppsPayloadPage2},
		"",
		http.StatusOK,
		"",
		nil}
	setup(mr, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	apps, err := client.Applications.ListAll()
	require.NoError(t, err)

	require.Len(t, apps, 2)

	require.Equal(t, "my_app", apps[0].Name)
	require.Equal(t, "my_app2", apps[1].Name)

	require.Equal(t, "STOPPED", apps[1].State)

	require.Equal(t, "java_buildpack", apps[0].Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "ruby_buildpack", apps[1].Lifecycle.BuildpackData.Buildpacks[0])
	require.Equal(t, "staticfile_buildpack", apps[1].Lifecycle.BuildpackData.Buildpacks[1])
}
