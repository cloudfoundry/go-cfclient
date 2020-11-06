package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListStacks(t *testing.T) {
	Convey("List Stacks", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/stacks", []string{listStacksPayloadPage1}, "", 200, "", nil},
			{"GET", "/v2/stacks_page_2", []string{listStacksPayloadPage2}, "", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		stacks, err := client.ListStacks()
		So(err, ShouldBeNil)

		So(len(stacks), ShouldEqual, 2)
		So(stacks[0].Guid, ShouldEqual, "67e019a3-322a-407a-96e0-178e95bd0e55")
		So(stacks[0].Name, ShouldEqual, "cflinuxfs2")
		So(stacks[0].Description, ShouldEqual, "Cloud Foundry Linux-based filesystem")
	})
}

func TestGetStackByGuid(t *testing.T) {
	Convey("Get Stack By Guid", t, func() {
		setup(MockRoute{"GET", "/v2/stacks/a9be2e10-0164-401d-94e0-88455d614844", []string{stackByGuidPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		stack, err := client.GetStackByGuid("a9be2e10-0164-401d-94e0-88455d614844")
		So(err, ShouldBeNil)

		So(stack.Guid, ShouldEqual, "a9be2e10-0164-401d-94e0-88455d614844")
		So(stack.Name, ShouldEqual, "windows2012R2")
		So(stack.Description, ShouldEqual, "Experimental Windows runtime")
	})
}
