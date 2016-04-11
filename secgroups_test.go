package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSecGroups(t *testing.T) {
	Convey("List SecGroups", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/security_groups", listSecGroupsPayload},
			{"GET", "/v2/security_groupsPage2", listSecGroupsPayloadPage2},
		}
		setupMultiple(mocks)
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		SecGroups := client.ListSecGroups()
		So(len(SecGroups), ShouldEqual, 2)
		So(SecGroups[0].Guid, ShouldEqual, "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c")
		So(SecGroups[0].Name, ShouldEqual, "secgroup-test")
		So(SecGroups[0].Running, ShouldEqual, true)
		So(SecGroups[0].Staging, ShouldEqual, true)
		So(SecGroups[0].Rules[0].Protocol, ShouldEqual, "tcp")
		So(SecGroups[0].Rules[0].Ports, ShouldEqual, "443,4443")
		So(SecGroups[0].Rules[0].Destination, ShouldEqual, "1.1.1.1")
		So(SecGroups[0].Rules[1].Protocol, ShouldEqual, "udp")
		So(SecGroups[0].Rules[1].Ports, ShouldEqual, "1111")
		So(SecGroups[0].Rules[1].Destination, ShouldEqual, "1.2.3.4")
		So(SecGroups[0].SpacesURL, ShouldEqual, "/v2/security_groups/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/spaces")
		So(SecGroups[0].SpacesData, ShouldBeEmpty)
		So(SecGroups[1].Guid, ShouldEqual, "f9ad202b-76dd-44ec-b7c2-fd2417a561e8")
		So(SecGroups[1].Name, ShouldEqual, "secgroup-test2")
		So(SecGroups[1].Running, ShouldEqual, false)
		So(SecGroups[1].Staging, ShouldEqual, false)
		So(SecGroups[1].Rules[0].Protocol, ShouldEqual, "udp")
		So(SecGroups[1].Rules[0].Ports, ShouldEqual, "2222")
		So(SecGroups[1].Rules[0].Destination, ShouldEqual, "2.2.2.2")
		So(SecGroups[1].Rules[1].Protocol, ShouldEqual, "tcp")
		So(SecGroups[1].Rules[1].Ports, ShouldEqual, "443,4443")
		So(SecGroups[1].Rules[1].Destination, ShouldEqual, "4.3.2.1")
		So(SecGroups[1].SpacesData[0].Entity.Guid, ShouldEqual, "e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4")
		So(SecGroups[1].SpacesData[0].Entity.Name, ShouldEqual, "space-test")
		So(SecGroups[1].SpacesData[1].Entity.Guid, ShouldEqual, "a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333")
		So(SecGroups[1].SpacesData[1].Entity.Name, ShouldEqual, "space-test2")
		So(SecGroups[1].SpacesData[2].Entity.Guid, ShouldEqual, "c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1")
		So(SecGroups[1].SpacesData[2].Entity.Name, ShouldEqual, "space-test3")
	})
}
