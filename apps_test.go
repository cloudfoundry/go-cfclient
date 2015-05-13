package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListApps(t *testing.T) {
	Convey("List Apps", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps", listAppsPayload},
			{"GET", "/v2/appsPage2", listAppsPayloadPage2},
		}
		setupMultiple(mocks)
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		apps := client.ListApps()
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
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", appPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		app := client.AppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
	})

	Convey("App By GUID with environment variables with different types", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", appPayloadWithEnvironment_json})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		app := client.AppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Environment["string"], ShouldEqual, "string")
		So(app.Environment["int"], ShouldEqual, 1)
	})
}

func TestAppSpace(t *testing.T) {
	Convey("Find app space", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/foobar", spacePayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		app := &App{
			Guid:     "123",
			Name:     "test app",
			SpaceURL: "/v2/spaces/foobar",
			c:        client,
		}
		space := app.Space()
		So(space.Name, ShouldEqual, "test-space")
		So(space.Guid, ShouldEqual, "a72fa1e8-c694-47b3-85f2-55f61fd00d73")
	})
}
