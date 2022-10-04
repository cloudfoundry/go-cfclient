package client

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListDomains(t *testing.T) {
	Convey("List  Domains ", t, func() {
		setup(MockRoute{"GET", "/v3/domains", []string{listDomainsPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.ListDomains(nil)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp, ShouldHaveLength, 1)
		So(resp[0].Name, ShouldEqual, "test-domain.com")
		So(resp[0].GUID, ShouldEqual, "3a5d3d89-3f89-4f05-8188-8a2b298c79d5")
		So(resp[0].Internal, ShouldEqual, false)
		So(resp[0].Relationships.Organization.Data.GUID, ShouldEqual, "3a3f3d89-3f89-4f05-8188-751b298c79d5")
		So(resp[0].Relationships.SharedOrganizations.Data[0].GUID, ShouldEqual, "404f3d89-3f89-6z72-8188-751b298d88d5")
		So(resp[0].Relationships.SharedOrganizations.Data[1].GUID, ShouldEqual, "416d3d89-3f89-8h67-2189-123b298d3592")
		So(resp[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5")
		So(resp[0].Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/3a3f3d89-3f89-4f05-8188-751b298c79d5")
		So(resp[0].Links["route_reservations"].Href, ShouldEqual, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/route_reservations")
		So(resp[0].Links["shared_organizations"].Href, ShouldEqual, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/relationships/shared_organizations")
	})
}
