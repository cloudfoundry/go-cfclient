package cfclient

import (
	"net/http"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServiceBindings(t *testing.T) {
	Convey("List Service Bindings", t, func() {
		mocks := []MockRoute{
			{
				Method:   "GET",
				Endpoint: "/v2/service_bindings",
				Output:   listServiceBindingsPayloadPage1,
				Status:   200,
			},
			{
				Method:   "GET",
				Endpoint: "/v2/service_bindings2",
				Output:   listServiceBindingsPayloadPage2,
				Status:   200,
			},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		serviceBindings, err := client.ListServiceBindings()
		So(err, ShouldBeNil)

		So(len(serviceBindings), ShouldEqual, 2)
		So(serviceBindings[0].Guid, ShouldEqual, "aa599bb3-4811-405a-bbe3-a68c7c55afc8")
		So(serviceBindings[0].AppGuid, ShouldEqual, "b26e7e98-f002-41a8-a663-1b60f808a92a")
		So(serviceBindings[0].ServiceInstanceGuid, ShouldEqual, "bde206e0-1ee8-48ad-b794-44c857633d50")
		So(reflect.DeepEqual(
			serviceBindings[0].Credentials,
			map[string]interface{}{"creds-key-66": "creds-val-66"}), ShouldBeTrue)
		So(serviceBindings[0].BindingOptions, ShouldBeEmpty)
		So(serviceBindings[0].GatewayData, ShouldBeNil)
		So(serviceBindings[0].GatewayName, ShouldEqual, "")
		So(serviceBindings[0].SyslogDrainUrl, ShouldEqual, "")
		So(serviceBindings[0].VolumeMounts, ShouldBeEmpty)
		So(serviceBindings[0].AppUrl, ShouldEqual, "/v2/apps/b26e7e98-f002-41a8-a663-1b60f808a92a")
		So(serviceBindings[0].ServiceInstanceUrl, ShouldEqual, "/v2/service_instances/bde206e0-1ee8-48ad-b794-44c857633d50")
		So(serviceBindings[1].Guid, ShouldEqual, "8201b87d-b273-4fdf-8dd4-4b42ce970cc7")
		So(serviceBindings[1].AppGuid, ShouldEqual, "636bbf83-5b54-488d-9528-066f680a99dc")
		So(serviceBindings[1].ServiceInstanceGuid, ShouldEqual, "c3023201-44f5-4dc8-a903-c69c9eba9809")
		So(reflect.DeepEqual(
			serviceBindings[1].Credentials,
			map[string]interface{}{"creds-key-66": "creds-val-66"}), ShouldBeTrue)
		So(serviceBindings[1].BindingOptions, ShouldBeEmpty)
		So(serviceBindings[1].GatewayData, ShouldBeNil)
		So(serviceBindings[1].GatewayName, ShouldEqual, "")
		So(serviceBindings[1].SyslogDrainUrl, ShouldEqual, "")
		So(serviceBindings[1].VolumeMounts, ShouldBeEmpty)
		So(serviceBindings[1].AppUrl, ShouldEqual, "/v2/apps/636bbf83-5b54-488d-9528-066f680a99dc")
		So(serviceBindings[1].ServiceInstanceUrl, ShouldEqual, "/v2/service_instances/c3023201-44f5-4dc8-a903-c69c9eba9809")

	})
}
func TestServiceBindingByGuid(t *testing.T) {
	Convey("Service Binding By Guid", t, func() {
		setup(MockRoute{"GET", "/v2/service_bindings/foo-bar-baz", serviceBindingByGuidPayload, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		serviceBinding, err := client.GetServiceBindingByGuid("foo-bar-baz")
		So(err, ShouldBeNil)

		So(serviceBinding.Guid, ShouldEqual, "foo-bar-baz")
		So(serviceBinding.AppGuid, ShouldEqual, "app-bar-baz")
	})
}

func TestDeleteServiceBinding(t *testing.T) {
	Convey("Delete service binding", t, func() {
		setup(MockRoute{"DELETE", "/v2/service_bindings/guid", "", "", http.StatusNoContent, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteServiceBinding("guid")
		So(err, ShouldBeNil)
	})
}

func TestCreateServiceBinding(t *testing.T) {
	Convey("Create service binding", t, func() {
		body := `{"app_guid":"app-guid","service_instance_guid":"service-instance-guid"}`
		setup(MockRoute{"POST", "/v2/service_bindings", postServiceBindingPayload, "", http.StatusCreated, "", &body}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		binding, err := client.CreateServiceBinding("app-guid", "service-instance-guid")
		So(err, ShouldBeNil)
		So(binding.Guid, ShouldEqual, "4e690cd4-66ef-4052-a23d-0d748316f18c")
		So(binding.AppGuid, ShouldEqual, "081d55a0-1bfa-4e51-8d08-273f764988db")
		So(binding.ServiceInstanceGuid, ShouldEqual, "a0029c76-7017-4a74-94b0-54a04ad94b80")
	})
}

func TestCreateRouteServiceBinding(t *testing.T) {
	Convey("Create route service binding", t, func() {
		setup(MockRoute{"PUT", "/v2/user_provided_service_instances/5badd282-6e07-4fc6-a8c4-78be99040774/routes/237d9236-7997-4b1a-be8d-2aaf2d85421a", "", "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.CreateRouteServiceBinding("237d9236-7997-4b1a-be8d-2aaf2d85421a", "5badd282-6e07-4fc6-a8c4-78be99040774")
		So(err, ShouldBeNil)
	})
}

func TestDeleteRouteServiceBinding(t *testing.T) {
	Convey("Delete route service binding", t, func() {
		setup(MockRoute{"DELETE", "/v2/service_instances/5badd282-6e07-4fc6-a8c4-78be99040774/routes/237d9236-7997-4b1a-be8d-2aaf2d85421a", "", "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteRouteServiceBinding("237d9236-7997-4b1a-be8d-2aaf2d85421a", "5badd282-6e07-4fc6-a8c4-78be99040774")
		So(err, ShouldBeNil)
	})
}
