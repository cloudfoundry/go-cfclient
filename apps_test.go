package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListApps(t *testing.T) {
	Convey("List Apps", t, func() {
		setup("GET", "/v2/apps", listAppsPayload)
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: server.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		apps := client.ListApps()
		So(len(apps), ShouldEqual, 1)
		So(apps[0].Guid, ShouldEqual, "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c")
		So(apps[0].Name, ShouldEqual, "app-test")
		So(apps[0].Environment["FOOBAR"], ShouldEqual, "QUX")
	})
}

func TestAppByGuid(t *testing.T) {
	Convey("App By GUID", t, func() {
		setup("GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", appPayload)
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: server.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		app := client.AppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
	})
}
