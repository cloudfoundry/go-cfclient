package client

import (
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateOrganization(t *testing.T) {
	Convey("Create  Organization", t, func() {
		expectedBody := `{"name":"my-org"}`
		setup(MockRoute{"POST", "/v3/organizations", []string{createOrganizationPayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.CreateOrganization(resource.CreateOrganizationRequest{
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

func TestGetOrganization(t *testing.T) {
	Convey("Get  Organization", t, func() {
		setup(MockRoute{"GET", "/v3/organizations/org-guid", []string{getOrganizationPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.GetOrganizationByGUID("org-guid")
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

func TestDeleteOrganization(t *testing.T) {
	Convey("Delete  Organization", t, func() {
		setup(MockRoute{"DELETE", "/v3/organizations/org-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteOrganization("org-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateOrganization(t *testing.T) {
	Convey("Update  Organization", t, func() {
		setup(MockRoute{"PATCH", "/v3/organizations/org-guid", []string{updateOrganizationPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organization, err := client.UpdateOrganization("org-guid", resource.UpdateOrganizationRequest{
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

func TestListOrganizationsByQuery(t *testing.T) {
	Convey("List  Organizations", t, func() {
		setup(MockRoute{"GET", "/v3/organizations", []string{listOrganizationsPayload, listOrganizationsPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		organizations, err := client.ListOrganizationsByQuery(nil)
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
