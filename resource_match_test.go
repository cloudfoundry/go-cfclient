package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResourceMatch(t *testing.T) {
	Convey("Resource Match", t, func() {
		mocks := []MockRoute{
			{"PUT", "/v2/resource_match", []string{listValidResources}, "Test-golang", 200, "", nil},
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
		ResourceList := []Resource{}
		ResourceList = append(ResourceList, Resource{Sha1: "002d760bea1be268e27077412e11a320d0f164d3", Size: 36})
		ResourceList = append(ResourceList, Resource{Sha1: "a9993e364706816aba3e25717850c26c9cd0d89d", Size: 1})

		ResApps, err := client.ResourceMatch(ResourceList)
		assertResourceList(ResApps, err)
	})
}

func assertResourceList(resList []Resource, err error) {
	So(err, ShouldBeNil)
	So(len(resList), ShouldEqual, 2)
	So(resList[0].Sha1, ShouldEqual, "002d760bea1be268e27077412e11a320d0f164d3")
	So(resList[0].Size, ShouldEqual, 36)
	So(resList[1].Sha1, ShouldEqual, "a9993e364706816aba3e25717850c26c9cd0d89d")
	So(resList[1].Size, ShouldEqual, 1)

}
