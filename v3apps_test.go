package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateV3App(t *testing.T) {
	Convey("Create V3 App", t, func() {
		expectedBody := `{"environment_variables":{"FOO":"BAR"},"name":"my-app","relationships":{"space":{"data":{"guid":"space-guid"}}}}`
		setup(MockRoute{"POST", "/v3/apps", []string{createV3AppPayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.CreateV3App(CreateV3AppRequest{
			Name:                 "my-app",
			SpaceGUID:            "space-guid",
			EnvironmentVariables: map[string]string{"FOO": "BAR"},
		})
		So(err, ShouldBeNil)
		So(app, ShouldNotBeNil)

		So(app.GUID, ShouldEqual, "app-guid")
		So(app.Relationships["space"].Data.GUID, ShouldEqual, "space-guid")
		So(app.Lifecycle.Type, ShouldEqual, "buildpack")
		So(app.Lifecycle.BuildpackData.Buildpacks, ShouldHaveLength, 1)
		So(app.Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "java_buildpack")
		So(app.Lifecycle.BuildpackData.Stack, ShouldEqual, "cflinuxfs2")
		So(app.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/space-guid")
		So(app.Metadata.Annotations, ShouldHaveLength, 0)
	})
}

func TestGetV3App(t *testing.T) {
	Convey("Get V3 App", t, func() {
		setup(MockRoute{"GET", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{getV3AppPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.GetV3AppByGUID("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)
		So(app, ShouldNotBeNil)

		So(app.GUID, ShouldEqual, "1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(app.Name, ShouldEqual, "my_app")
		So(app.Lifecycle.Type, ShouldEqual, "buildpack")
		So(app.Lifecycle.BuildpackData.Buildpacks, ShouldHaveLength, 1)
		So(app.Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "java_buildpack")
		So(app.Lifecycle.BuildpackData.Stack, ShouldEqual, "cflinuxfs2")
		So(app.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576")
		So(app.Metadata.Annotations, ShouldHaveLength, 1)
		So(app.Metadata.Annotations["contacts"], ShouldEqual, "Bill tel(1111111) email(bill@fixme), Bob tel(222222) pager(3333333#555) email(bob@fixme)")
	})
}

func TestGetV3AppEnv(t *testing.T) {
	Convey("Get V3 App Environment", t, func() {
		setup(MockRoute{"GET", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/env", []string{getV3AppEnvPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		env, err := client.GetV3AppEnvironment("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)

		So(env.EnvVars, ShouldHaveLength, 1)
		So(env.EnvVars["RAILS_ENV"], ShouldEqual, "production")

		So(env.StagingEnv, ShouldHaveLength, 1)
		So(env.StagingEnv["GEM_CACHE"], ShouldEqual, "http://gem-cache.example.org")

		So(env.RunningEnv, ShouldHaveLength, 1)
		So(env.RunningEnv["HTTP_PROXY"], ShouldEqual, "http://proxy.example.org")

		So(env.SystemEnvVars, ShouldHaveLength, 1)
		So(env.SystemEnvVars, ShouldContainKey, "VCAP_SERVICES")

		So(env.AppEnvVars, ShouldHaveLength, 1)
		So(env.AppEnvVars, ShouldContainKey, "VCAP_APPLICATION")
	})
}

func TestStartV3App(t *testing.T) {
	Convey("Start V3 App", t, func() {
		setup(MockRoute{"POST", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start", []string{startV3AppPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.StartV3App("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)
		So(app, ShouldNotBeNil)

		So(app.State, ShouldEqual, "STARTED")
		So(app.GUID, ShouldEqual, "1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(app.Name, ShouldEqual, "my_app")
		So(app.Lifecycle.Type, ShouldEqual, "buildpack")
		So(app.Lifecycle.BuildpackData.Buildpacks, ShouldHaveLength, 1)
		So(app.Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "java_buildpack")
		So(app.Lifecycle.BuildpackData.Stack, ShouldEqual, "cflinuxfs2")
		So(app.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576")
		So(app.Metadata.Annotations, ShouldHaveLength, 0)
	})
}

func TestDeleteV3App(t *testing.T) {
	Convey("Delete V3 App", t, func() {
		setup(MockRoute{"DELETE", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteV3App("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)
	})
}

func TestUpdateV3App(t *testing.T) {
	Convey("Update V3 App", t, func() {
		setup(MockRoute{"PATCH", "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{updateV3AppPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.UpdateV3App("1cb006ee-fb05-47e1-b541-c34179ddc446", UpdateV3AppRequest{})
		So(err, ShouldBeNil)
		So(app, ShouldNotBeNil)

		So(app.GUID, ShouldEqual, "1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(app.State, ShouldEqual, "STARTED")
		So(app.Name, ShouldEqual, "my_app")
		So(app.Lifecycle.Type, ShouldEqual, "buildpack")
		So(app.Lifecycle.BuildpackData.Buildpacks, ShouldHaveLength, 1)
		So(app.Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "java_buildpack")
		So(app.Lifecycle.BuildpackData.Stack, ShouldEqual, "cflinuxfs2")
		So(app.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576")
		So(app.Metadata.Annotations, ShouldHaveLength, 0)
		So(app.Metadata.Labels, ShouldHaveLength, 2)
		So(app.Metadata.Labels["environment"], ShouldEqual, "production")
		So(app.Metadata.Labels["internet-facing"], ShouldEqual, "false")
	})
}

func TestListV3AppsByQuery(t *testing.T) {
	Convey("List V3 Apps", t, func() {
		setup(MockRoute{"GET", "/v3/apps", []string{listV3AppsPayload, listV3AppsPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		apps, err := client.ListV3AppsByQuery(nil)
		So(err, ShouldBeNil)
		So(apps, ShouldHaveLength, 2)

		So(apps[0].Name, ShouldEqual, "my_app")
		So(apps[1].Name, ShouldEqual, "my_app2")

		So(apps[1].State, ShouldEqual, "STOPPED")

		So(apps[0].Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "java_buildpack")
		So(apps[1].Lifecycle.BuildpackData.Buildpacks[0], ShouldEqual, "ruby_buildpack")
		So(apps[1].Lifecycle.BuildpackData.Buildpacks[1], ShouldEqual, "staticfile_buildpack")
	})
}
