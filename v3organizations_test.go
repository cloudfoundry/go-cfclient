package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateV3Organization(t *testing.T) {
	Convey("Create V3 Organization", t, func() {
		expectedBody := `{"name":"my-org"}`
		setup(MockRoute{"POST", "/v3/organizations", []string{createV3OrganizationPayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.CreateV3Organization(CreateV3OrganizationRequest{
			Name: "my-org",
		})
		So(err, ShouldBeNil)
		So(organization, ShouldNotBeNil)

		So(organization.GUID, ShouldEqual, "org-guid")
		So(organization.Relationships["quota"].Data.GUID, ShouldEqual, "quota-guid")
		So(organization.Links["domains"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid/domains")
		So(organization.Metadata.Annotations, ShouldHaveLength, 0)
		So(organization.Metadata.Labels, ShouldContainKey, "ORG_KEY")
		So(organization.Metadata.Labels["ORG_KEY"], ShouldEqual, "org_value")
	})
}

func TestGetV3Organization(t *testing.T) {
	Convey("Get V3 Organization", t, func() {
		setup(MockRoute{"GET", "/v3/organizations/org-guid", []string{getV3OrganizationPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.GetV3OrganizationByGUID("org-guid")
		So(err, ShouldBeNil)
		So(organization, ShouldNotBeNil)

		So(organization.GUID, ShouldEqual, "org-guid")
		So(organization.Relationships["quota"].Data.GUID, ShouldEqual, "quota-guid")
		So(organization.Links["domains"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid/domains")
		So(organization.Metadata.Annotations, ShouldHaveLength, 0)
		So(organization.Metadata.Labels, ShouldContainKey, "ORG_KEY")
		So(organization.Metadata.Labels["ORG_KEY"], ShouldEqual, "org_value")
	})
}

func TestDeleteV3Organization(t *testing.T) {
	Convey("Delete V3 Organization", t, func() {
		setup(MockRoute{"DELETE", "/v3/organizations/org-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteV3Organization("org-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateV3Organization(t *testing.T) {
	Convey("Update V3 Organization", t, func() {
		setup(MockRoute{"PATCH", "/v3/organizations/org-guid", []string{updateV3OrganizationPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.UpdateV3Organization("org-guid", UpdateV3OrganizationRequest{
			Name: "my-org",
		})
		So(err, ShouldBeNil)
		So(organization, ShouldNotBeNil)

		So(organization.Name, ShouldEqual, "my-org")
		So(organization.GUID, ShouldEqual, "org-guid")
		So(organization.Relationships["quota"].Data.GUID, ShouldEqual, "")
		So(organization.Links["domains"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid/domains")
		So(organization.Metadata.Annotations, ShouldHaveLength, 0)
		So(organization.Metadata.Labels, ShouldContainKey, "ORG_KEY")
		So(organization.Metadata.Labels["ORG_KEY"], ShouldEqual, "org_value")
	})
}

func TestListV3OrganizationsByQuery(t *testing.T) {
	Convey("List V3 Organizations", t, func() {
		setup(MockRoute{"GET", "/v3/organizations", []string{listV3OrganizationsPayload, listV3OrganizationsPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organizations, err := client.ListV3OrganizationsByQuery(nil)
		So(err, ShouldBeNil)
		So(organizations, ShouldHaveLength, 2)

		So(organizations[0].Name, ShouldEqual, "my-org-1")
		So(organizations[1].Name, ShouldEqual, "my-org-2")

		So(organizations[0].Relationships["quota"].Data.GUID, ShouldEqual, "quota-guid")
		So(organizations[0].Links["domains"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid/domains")
		So(organizations[1].Relationships["quota"].Data.GUID, ShouldEqual, "")
		So(organizations[1].Links["domains"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid-2/domains")
	})
}
