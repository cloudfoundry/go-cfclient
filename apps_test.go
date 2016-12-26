package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListApps(t *testing.T) {
	Convey("List Apps", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps", listAppsPayload, "Test-golang"},
			{"GET", "/v2/appsPage2", listAppsPayloadPage2, "Test-golang"},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		apps, err := client.ListApps()
		So(err, ShouldBeNil)

		So(len(apps), ShouldEqual, 2)
		So(apps[0].Guid, ShouldEqual, "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c")
		So(apps[0].Name, ShouldEqual, "app-test")
		So(apps[0].Environment["FOOBAR"], ShouldEqual, "QUX")
		So(apps[1].Guid, ShouldEqual, "f9ad202b-76dd-44ec-b7c2-fd2417a561e8")
		So(apps[1].Name, ShouldEqual, "app-test2")
	})
}

func TestAppByGuid(t *testing.T) {
	Convey("App By GUID", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", appPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.AppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
	})

	Convey("App By GUID with environment variables with different types", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", appPayloadWithEnvironment_json, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.AppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(app.Environment["string"], ShouldEqual, "string")
		So(app.Environment["int"], ShouldEqual, 1)
	})
}

func TestGetAppInstances(t *testing.T) {
	Convey("App completely running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances", appInstancePayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appInstances, err := client.GetAppInstances("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appInstances["0"].State, ShouldEqual, "RUNNING")
		So(appInstances["1"].State, ShouldEqual, "RUNNING")
	})

	Convey("App partially running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances", appInstanceUnhealthyPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appInstances, err := client.GetAppInstances("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appInstances["0"].State, ShouldEqual, "RUNNING")
		So(appInstances["1"].State, ShouldEqual, "STARTING")
	})
}

func TestKillAppInstance(t *testing.T) {
	Convey("Kills an app instance", t, func() {
		setup(MockRoute{"DELETE", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances/0", "", ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		So(client.KillAppInstance("9902530c-c634-4864-a189-71d763cb12e2", "0"), ShouldBeNil)
	})
}

func TestAppSpace(t *testing.T) {
	Convey("Find app space", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/foobar", spacePayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app := &App{
			Guid:     "123",
			Name:     "test app",
			SpaceURL: "/v2/spaces/foobar",
			c:        client,
		}
		space, err := app.Space()
		So(err, ShouldBeNil)

		So(space.Name, ShouldEqual, "test-space")
		So(space.Guid, ShouldEqual, "a72fa1e8-c694-47b3-85f2-55f61fd00d73")
	})
}
