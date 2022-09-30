package client

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServiceInstancesByQuery(t *testing.T) {
	Convey("List  Service Instances", t, func() {
		setup(MockRoute{"GET", "/v3/service_instances", []string{listServiceInstancesPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		services, err := client.ListServiceInstances()
		So(err, ShouldBeNil)
		So(services, ShouldHaveLength, 1)

		So(services[0].Name, ShouldEqual, "my_service_instance")

		So(services[0].Relationships["space"].Data.GUID, ShouldEqual, "ae0031f9-dd49-461c-a945-df40e77c39cb")
		So(services[0].Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/ae0031f9-dd49-461c-a945-df40e77c39cb")
	})
}
