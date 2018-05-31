package cfclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestMappingAppAndRoute(t *testing.T) {
	Convey("Map app and route", t, func() {
		setup(MockRoute{"POST", "/v2/route_mappings", postRouteMappingsPayload, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		mappingRequest := MappingRequest{AppGUID: "fa23ddfc-b635-4205-8283-844c53122888", RouteGUID: "e00fb1e1-f7d4-4e36-9912-f76a587e9858", AppPort: 8888}

		mapping, err := client.MappingAppAndRoute(mappingRequest)
		So(err, ShouldBeNil)
		So(mapping.Guid, ShouldEqual, "f869fa46-22b1-40ee-b491-58e321345528")
		So(mapping.AppGUID, ShouldEqual, "fa23ddfc-b635-4205-8283-844c53122888")
		So(mapping.RouteGUID, ShouldEqual, "e00fb1e1-f7d4-4e36-9912-f76a587e9858")
	})
}
