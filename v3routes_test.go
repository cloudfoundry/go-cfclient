package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3Routes(t *testing.T) {
	Convey("List V3 Routes", t, func() {
		setup(MockRoute{"GET", "/v3/routes", []string{listV3RoutesPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		routes, err := client.ListV3Routes()
		So(err, ShouldBeNil)
		So(routes, ShouldHaveLength, 1)

		So(routes[0].Host, ShouldEqual, "a-hostname")
		So(routes[0].Path, ShouldEqual, "/some_path")
		So(routes[0].Url, ShouldEqual, "a-hostname.a-domain.com/some_path")

		So(routes[0].Relationships["space"].Data.GUID, ShouldEqual, "885a8cb3-c07b-4856-b448-eeb10bf36236")
		So(routes[0].Relationships["domain"].Data.GUID, ShouldEqual, "0b5f3633-194c-42d2-9408-972366617e0e")

		So(routes[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31")
		So(routes[0].Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236")
		So(routes[0].Links["domain"].Href, ShouldEqual, "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e")
		So(routes[0].Links["destinations"].Href, ShouldEqual, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations")
	})
}

func TestCreateV3Routes(t *testing.T) {
	Convey("Create V3 Route", t, func() {
		setup(MockRoute{"POST", "/v3/routes", []string{createV3RoutePayload}, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		route, err := client.CreateV3Route(
			"885a8cb3-c07b-4856-b448-eeb10bf36236",
			"0b5f3633-194c-42d2-9408-972366617e0e",
			nil,
		)
		So(err, ShouldBeNil)
		So(route.Host, ShouldEqual, "a-hostname")
		So(route.Path, ShouldEqual, "/some_path")
		So(route.Relationships["space"].Data.GUID, ShouldEqual, "885a8cb3-c07b-4856-b448-eeb10bf36236")
		So(route.Relationships["domain"].Data.GUID, ShouldEqual, "0b5f3633-194c-42d2-9408-972366617e0e")
		So(route.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31")
		So(route.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236")
		So(route.Links["domain"].Href, ShouldEqual, "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e")
		So(route.Links["destinations"].Href, ShouldEqual, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations")
	})
}
