package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateApp(t *testing.T) {
	Convey("Update app", t, func() {
		setup(MockRoute{"PUT", "/v2/apps/97f7e56b-addf-4d26-be82-998a06600011", []string{AppUpdatePayload}, "", 201, "", nil}, t)
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		aur := AppUpdateResource{Name: "NewName", DiskQuota: 1024, Instances: 1, Memory: 65}
		ret, err := client.UpdateApp("97f7e56b-addf-4d26-be82-998a06600011", aur)
		So(err, ShouldBeNil)
		So(ret.Entity.Memory, ShouldEqual, 65)
		So(ret.Entity.Instances, ShouldEqual, 1)
		So(ret.Entity.DiskQuota, ShouldEqual, 1024)
		So(ret.Entity.Name, ShouldEqual, "NewName")
	})
}

func TestRestageApp(t *testing.T) {
	Convey("Restage app", t, func() {
		setup(MockRoute{"POST", "/v2/apps/97f7e56b-addf-4d26-be82-998a06600011/restage", []string{appRestagePayload}, "", 201, "", nil}, t)
		client, err := NewClient(&Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		})
		So(err, ShouldBeNil)

		resp, err := client.RestageApp("97f7e56b-addf-4d26-be82-998a06600011")
		So(err, ShouldBeNil)
		So(resp.Metadata.Guid, ShouldEqual, "97f7e56b-addf-4d26-be82-998a06600011")
		So(resp.Entity.Name, ShouldEqual, "name-2047")
		So(resp.Entity.EnableSSH, ShouldBeTrue)
	})
}
