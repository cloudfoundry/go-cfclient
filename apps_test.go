package cfclient

import (
	"testing"
	"time"

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
		So(apps[0].Memory, ShouldEqual, 256)
		So(apps[0].Instances, ShouldEqual, 1)
		So(apps[0].DiskQuota, ShouldEqual, 1024)
		So(apps[0].SpaceGuid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(apps[0].StackGuid, ShouldEqual, "2c531037-68a2-4e2c-a9e0-71f9d0abf0d4")
		So(apps[0].State, ShouldEqual, "STARTED")
		So(apps[0].Command, ShouldEqual, "")
		So(apps[0].Buildpack, ShouldEqual, "https://github.com/cloudfoundry/buildpack-go.git")
		So(apps[0].DetectedBuildpack, ShouldEqual, "")
		So(apps[0].DetectedBuildpackGuid, ShouldEqual, "0d22f6a1-76c5-417f-ac6c-d9d21463ecbc")
		So(apps[0].HealthCheckHttpEndpoint, ShouldEqual, "")
		So(apps[0].HealthCheckType, ShouldEqual, "port")
		So(apps[0].HealthCheckTimeout, ShouldEqual, 0)
		So(apps[0].Diego, ShouldEqual, true)
		So(apps[0].EnableSSH, ShouldEqual, true)
		So(apps[0].DetectedStartCommand, ShouldEqual, "app-launching-service-broker")
		So(apps[0].DockerImage, ShouldEqual, "")
		So(apps[0].DockerCredentials["redacted_message"], ShouldEqual, "[PRIVATE DATA HIDDEN]")
		So(apps[0].Environment["FOOBAR"], ShouldEqual, "QUX")
		So(apps[0].StagingFailedReason, ShouldEqual, "")
		So(apps[0].StagingFailedDescription, ShouldEqual, "")
		So(len(apps[0].Ports), ShouldEqual, 1)
		So(apps[0].Ports[0], ShouldEqual, 8080)

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

		var d0 float64 = 1455210430.5104606
		var d1 float64 = 1455210430.3912115
		date0 := time.Unix(int64(d0), 0)
		date1 := time.Unix(int64(d1), 0)

		So(appInstances["0"].Since.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appInstances["1"].Since.Format(time.UnixDate), ShouldEqual, date1.Format(time.UnixDate))
		So(appInstances["0"].Since.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appInstances["1"].Since.ToTime(), ShouldHaveSameTypeAs, date1)

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

		var d0 float64 = 1455210430.5104606
		var d1 float64 = 1455210430.3912115
		date0 := time.Unix(int64(d0), 0)
		date1 := time.Unix(int64(d1), 0)

		So(appInstances["0"].Since.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appInstances["1"].Since.Format(time.UnixDate), ShouldEqual, date1.Format(time.UnixDate))
		So(appInstances["0"].Since.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appInstances["1"].Since.ToTime(), ShouldHaveSameTypeAs, date1)

	})
}

func TestGetAppStats(t *testing.T) {
	Convey("App stats completely running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/stats", appStatsPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appStats, err := client.GetAppStats("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appStats["0"].State, ShouldEqual, "RUNNING")
		So(appStats["1"].State, ShouldEqual, "RUNNING")

		date0, _ := time.Parse("2006-01-02 15:04:05 -0700", "2016-09-17 15:46:17 +0000")
		date1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2016-09-17 15:46:17 +0000")

		So(appStats["0"].Stats.Usage.Time.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appStats["1"].Stats.Usage.Time.Format(time.UnixDate), ShouldEqual, date1.Format(time.UnixDate))
		So(appStats["0"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appStats["1"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date1)
		So(appStats["0"].Stats.Usage.CPU, ShouldEqual, 0.36580239597146486)
		So(appStats["1"].Stats.Usage.CPU, ShouldEqual, 0.33857742931636664)
		So(appStats["0"].Stats.Usage.Mem, ShouldEqual, 518123520)
		So(appStats["1"].Stats.Usage.Mem, ShouldEqual, 530731008)
		So(appStats["0"].Stats.Usage.Disk, ShouldEqual, 151150592)
		So(appStats["1"].Stats.Usage.Disk, ShouldEqual, 151150592)

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
