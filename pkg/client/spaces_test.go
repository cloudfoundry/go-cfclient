package client

import (
	"net/http"
	"testing"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateSpace(t *testing.T) {
	Convey("Create  Space", t, func() {
		expectedBody := `{"name":"my-space","relationships":{"organization":{"data":{"guid":"org-guid"}}}}`
		setup(MockRoute{"POST", "/v3/spaces", []string{createSpacePayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.CreateSpace(v3.CreateSpaceRequest{
			Name:    "my-space",
			OrgGUID: "org-guid",
		})
		So(err, ShouldBeNil)
		So(space, ShouldNotBeNil)

		So(space.GUID, ShouldEqual, "space-guid")
		So(space.Relationships["organization"].Data.GUID, ShouldEqual, "org-guid")
		So(space.Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid")
		So(space.Metadata.Annotations, ShouldHaveLength, 0)
		So(space.Metadata.Labels, ShouldContainKey, "SPACE_KEY")
		So(space.Metadata.Labels["SPACE_KEY"], ShouldEqual, "space_value")
	})
}

func TestGetSpace(t *testing.T) {
	Convey("Get  Space", t, func() {
		setup(MockRoute{"GET", "/v3/spaces/space-guid", []string{getSpacePayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.GetSpaceByGUID("space-guid")
		So(err, ShouldBeNil)
		So(space, ShouldNotBeNil)

		So(space.GUID, ShouldEqual, "space-guid")
		So(space.Relationships["organization"].Data.GUID, ShouldEqual, "org-guid")
		So(space.Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid")
		So(space.Metadata.Annotations, ShouldHaveLength, 0)
		So(space.Metadata.Labels, ShouldContainKey, "SPACE_KEY")
		So(space.Metadata.Labels["SPACE_KEY"], ShouldEqual, "space_value")
	})
}

func TestDeleteSpace(t *testing.T) {
	Convey("Delete  Space", t, func() {
		setup(MockRoute{"DELETE", "/v3/spaces/space-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteSpace("space-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateSpace(t *testing.T) {
	Convey("Update  Space", t, func() {
		setup(MockRoute{"PATCH", "/v3/spaces/space-guid", []string{updateSpacePayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.UpdateSpace("space-guid", v3.UpdateSpaceRequest{
			Name: "my-space",
		})
		So(err, ShouldBeNil)
		So(space, ShouldNotBeNil)

		So(space.Name, ShouldEqual, "my-space")
		So(space.GUID, ShouldEqual, "space-guid")
		So(space.Relationships["organization"].Data.GUID, ShouldEqual, "org-guid")
		So(space.Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid")
		So(space.Metadata.Annotations, ShouldHaveLength, 0)
		So(space.Metadata.Labels, ShouldContainKey, "SPACE_KEY")
		So(space.Metadata.Labels["SPACE_KEY"], ShouldEqual, "space_value")
	})
}

func TestListSpacesByQuery(t *testing.T) {
	Convey("List  Spaces", t, func() {
		setup(MockRoute{"GET", "/v3/spaces", []string{listSpacesPayload, listSpacesPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.ListSpacesByQuery(nil)
		So(err, ShouldBeNil)
		So(spaces, ShouldHaveLength, 2)

		So(spaces[0].Name, ShouldEqual, "my-space-1")
		So(spaces[1].Name, ShouldEqual, "my-space-2")

		So(spaces[0].Relationships["organization"].Data.GUID, ShouldEqual, "org-guid")
		So(spaces[0].Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid")
		So(spaces[1].Relationships["organization"].Data.GUID, ShouldEqual, "org-guid")
		So(spaces[1].Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/org-guid")
	})
}

func TestListSpaceUsersByQuery(t *testing.T) {
	Convey("List  Space Users", t, func() {
		setup(MockRoute{"GET", "/v3/spaces/space-guid/users", []string{listSpaceUsersPayload, listSpaceUsersPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListSpaceUsers("space-guid")
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 2)

		So(users[0].Username, ShouldEqual, "some-name-1")
		So(users[1].Username, ShouldEqual, "some-name-2")

		So(users[0].PresentationName, ShouldEqual, "some-name-1")
		So(users[0].Origin, ShouldEqual, "uaa")
		So(users[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/users/10a93b89-3f89-4f05-7238-8a2b123c79l9")
		So(users[1].PresentationName, ShouldEqual, "some-name-2")
		So(users[1].Origin, ShouldEqual, "ldap")
		So(users[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/users/9da93b89-3f89-4f05-7238-8a2b123c79l9")
	})
}
