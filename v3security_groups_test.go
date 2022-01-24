package cfclient

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3SecurityGroupsByQuery(t *testing.T) {
	Convey("List All V3 Security Groups", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listV3SecurityGroupsPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.ListV3SecGroupsByQuery(nil)
		So(err, ShouldBeNil)
		So(securityGroups, ShouldHaveLength, 2)

		So(securityGroups[0].Name, ShouldEqual, "my-group1")
		So(securityGroups[0].GUID, ShouldEqual, "guid-1")
		So(securityGroups[1].Name, ShouldEqual, "my-group2")
		So(securityGroups[1].GUID, ShouldEqual, "guid-2")

		So(securityGroups[0].GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups[0].Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups[0].Rules[0].Destination, ShouldEqual, "1.2.3.4/10")
		So(securityGroups[0].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroups[0].Rules[1].Type, ShouldEqual, 8)
		So(securityGroups[0].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[0].Rules[1].Description, ShouldEqual, "test-desc-1")
		So(securityGroups[0].Relationships["staging_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
		So(securityGroups[1].GloballyEnabled.Staging, ShouldEqual, true)
		So(securityGroups[1].Rules[1].Protocol, ShouldEqual, "icmp")
		So(securityGroups[1].Rules[1].Destination, ShouldEqual, "1.2.3.4/16")
		So(securityGroups[1].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroups[1].Rules[1].Type, ShouldEqual, 5)
		So(securityGroups[1].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[1].Rules[1].Description, ShouldEqual, "test-desc-2")
		So(securityGroups[1].Relationships["running_spaces"].Data[0].GUID, ShouldEqual, "space-guid-5")
		So(securityGroups[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-2")
	})

	Convey("List V3 Security Groups By GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listV3SecurityGroupsByGuidPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["guids"] = []string{"guid-1"}

		securityGroups, err := client.ListV3SecGroupsByQuery(query)
		So(err, ShouldBeNil)
		So(securityGroups, ShouldHaveLength, 1)

		So(securityGroups[0].Name, ShouldEqual, "my-group1")
		So(securityGroups[0].GUID, ShouldEqual, "guid-1")

		So(securityGroups[0].GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups[0].Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups[0].Rules[0].Destination, ShouldEqual, "1.2.3.4/10")
		So(securityGroups[0].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroups[0].Rules[1].Type, ShouldEqual, 8)
		So(securityGroups[0].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[0].Rules[1].Description, ShouldEqual, "test-desc-1")
		So(securityGroups[0].Relationships["staging_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
	})
}
