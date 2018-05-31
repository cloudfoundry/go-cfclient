package cfclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestMappingAppAndRoute(t *testing.T) {
	Convey("Mapping app and route", t, func() {
		setup(MockRoute{"POST", "/v2/route_mappings", postRouteMappingsPayload, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		mappingRequest := RouteMappingRequest{AppGUID: "fa23ddfc-b635-4205-8283-844c53122888", RouteGUID: "e00fb1e1-f7d4-4e36-9912-f76a587e9858", AppPort: 8888}

		mapping, err := client.MappingAppAndRoute(mappingRequest)
		So(err, ShouldBeNil)
		So(mapping.Guid, ShouldEqual, "f869fa46-22b1-40ee-b491-58e321345528")
		So(mapping.AppGUID, ShouldEqual, "fa23ddfc-b635-4205-8283-844c53122888")
		So(mapping.RouteGUID, ShouldEqual, "e00fb1e1-f7d4-4e36-9912-f76a587e9858")
	})
}

func TestListRouteMappings(t *testing.T) {
	Convey("List Route Mappings", t, func() {
		setup(MockRoute{"GET", "/v2/route_mappings", listRouteMappingsPayload, "", http.StatusOK, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token: "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		routeMappings, err := client.ListRouteMappings()
		So(err, ShouldBeNil)

		So(len(routeMappings), ShouldEqual, 2)
		So(routeMappings[0].Guid, ShouldEqual, "63603ed7-bd4a-4475-a371-5b34381e0cf7")
		So(routeMappings[1].Guid, ShouldEqual, "63603ed7-bd4a-4475-a371-5b34381e0cf8")
		So(routeMappings[0].AppGUID, ShouldEqual, "ee8b175a-2228-4931-be8a-1f6445bd63bc")
		So(routeMappings[1].AppGUID, ShouldEqual, "ee8b175a-2228-4931-be8a-1f6445bd63bd")
		So(routeMappings[0].RouteGUID, ShouldEqual, "eb1c4fcd-7d6d-41d2-bd2f-5811f53b6677")
		So(routeMappings[0].RouteGUID, ShouldEqual, "eb1c4fcd-7d6d-41d2-bd2f-5811f53b6678")
	})
}
