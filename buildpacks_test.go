package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListBuildpacks(t *testing.T) {
	Convey("List buildpack", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/buildpacks", listBuildpacksPayload, "", 200, "", nil},
			{"GET", "/v2/buildpacksPage2", listBuildpacksPayload2, "", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		buildpacks, err := client.ListBuildpacks()
		So(err, ShouldBeNil)

		So(len(buildpacks), ShouldEqual, 6)
		So(buildpacks[0].Guid, ShouldEqual, "c92b6f5f-d2a4-413a-b515-647d059723aa")
		So(buildpacks[0].CreatedAt, ShouldEqual, "2016-06-08T16:41:31Z")
		So(buildpacks[0].UpdatedAt, ShouldEqual, "2016-06-08T16:41:26Z")
		So(buildpacks[0].Name, ShouldEqual, "name_1")
	})
}

func TestGetBuildpackByGuid(t *testing.T) {
	Convey("A buildpack", t, func() {
		setup(MockRoute{"GET", "/v2/buildpacks/c92b6f5f-d2a4-413a-b515-647d059723aa", buildpackPayload, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		buildpack, err := client.GetBuildpackByGuid("c92b6f5f-d2a4-413a-b515-647d059723aa")
		So(err, ShouldBeNil)

		So(buildpack.Guid, ShouldEqual, "c92b6f5f-d2a4-413a-b515-647d059723aa")
		So(buildpack.CreatedAt, ShouldEqual, "2016-06-08T16:41:31Z")
		So(buildpack.UpdatedAt, ShouldEqual, "2016-06-08T16:41:26Z")
		So(buildpack.Name, ShouldEqual, "name_1")
	})
}
