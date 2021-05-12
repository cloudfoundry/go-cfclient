package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateV3Space(t *testing.T) {
	Convey("Create V3 Space", t, func() {
		expectedBody := `{"name":"my-space","relationships":{"organization":{"data":{"guid":"org-guid"}}}}`
		setup(MockRoute{"POST", "/v3/spaces", []string{createV3SpacePayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.CreateV3Space(CreateV3SpaceRequest{
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

func TestGetV3Space(t *testing.T) {
	Convey("Get V3 Space", t, func() {
		setup(MockRoute{"GET", "/v3/spaces/space-guid", []string{getV3SpacePayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.GetV3SpaceByGUID("space-guid")
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

func TestDeleteV3Space(t *testing.T) {
	Convey("Delete V3 Space", t, func() {
		setup(MockRoute{"DELETE", "/v3/spaces/space-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteV3Space("space-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateV3Space(t *testing.T) {
	Convey("Update V3 Space", t, func() {
		setup(MockRoute{"PATCH", "/v3/spaces/space-guid", []string{updateV3SpacePayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space, err := client.UpdateV3Space("space-guid", UpdateV3SpaceRequest{
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

func TestListV3SpacesByQuery(t *testing.T) {
	Convey("List V3 Spaces", t, func() {
		setup(MockRoute{"GET", "/v3/spaces", []string{listV3SpacesPayload, listV3SpacesPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.ListV3SpacesByQuery(nil)
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
