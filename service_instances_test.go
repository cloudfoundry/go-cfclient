package cfclient

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServicesInstances(t *testing.T) {
	Convey("List Service Instances", t, func() {
		setup(MockRoute{"GET", "/v2/service_instances", []string{listServiceInstancePayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		instances, err := client.ListServiceInstances()
		So(err, ShouldBeNil)

		So(len(instances), ShouldEqual, 2)
		So(instances[0].Guid, ShouldEqual, "8423ca96-90ad-411f-b77a-0907844949fc")
		So(instances[0].Name, ShouldEqual, "fortunes-db")
	})
}

func TestGetServiceInstanceParams(t *testing.T) {
	Convey("Service instance parameters", t, func() {
		setup(MockRoute{"GET", "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/parameters", []string{serviceInstanceParamsPayload}, "", 200, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		params, err := client.GetServiceInstanceParams("8423ca96-90ad-411f-b77a-0907844949fc")
		So(err, ShouldBeNil)

		So(params, ShouldContainKey, "foo")
		So(params, ShouldContainKey, "baz")
		So(params["foo"], ShouldEqual, "bar")
		So(params["baz"], ShouldEqual, 42)
	})
}

func TestServiceInstanceByGuid(t *testing.T) {
	Convey("Service instance by Guid", t, func() {
		setup(MockRoute{"GET", "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc", []string{serviceInstancePayload}, "", 200, "", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		service, err := client.GetServiceInstanceByGuid("8423ca96-90ad-411f-b77a-0907844949fc")
		So(err, ShouldBeNil)

		expected := ServiceInstance{
			Guid:        "8423ca96-90ad-411f-b77a-0907844949fc",
			CreatedAt:   "2016-10-21T18:22:56Z",
			UpdatedAt:   "2016-10-21T18:22:56Z",
			Credentials: map[string]interface{}{},
			Name:        "fortunes-db",
			LastOperation: LastOperation{
				Type:        "create",
				State:       "succeeded",
				Description: "",
				UpdatedAt:   "",
				CreatedAt:   "2016-10-21T18:22:56Z",
			},
			Tags:               []string{},
			ServiceGuid:        "440ce9d9-b108-4bbe-80b4-08338f3cc25b",
			ServicePlanGuid:    "f48419f7-4717-4706-86e4-a24973848a77",
			SpaceGuid:          "21e5fdc7-5131-4743-8447-6373cf336a77",
			DashboardUrl:       "https://p-mysql.system.example.com/manage/instances/8423ca96-90ad-411f-b77a-0907844949fc",
			Type:               "managed_service_instance",
			SpaceUrl:           "/v2/spaces/21e5fdc7-5131-4743-8447-6373cf336a77",
			ServicePlanUrl:     "/v2/service_plans/f48419f7-4717-4706-86e4-a24973848a77",
			ServiceBindingsUrl: "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_bindings",
			ServiceKeysUrl:     "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_keys",
			RoutesUrl:          "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/routes",
			ServiceUrl:         "/v2/services/440ce9d9-b108-4bbe-80b4-08338f3cc25b",
			c:                  client,
		}
		So(service, ShouldResemble, expected)
	})
}

func TestCreateServiceInstance(t *testing.T) {
	Convey("Create service instance", t, func() {
		setup(MockRoute{"POST", "/v2/service_instances", []string{serviceInstancePayload}, "", 202, "accepts_incomplete=true", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		req := ServiceInstanceRequest{
			Name:            "test-service",
			ServicePlanGuid: "f48419f7-4717-4706-86e4-a24973848a77",
			SpaceGuid:       "21e5fdc7-5131-4743-8447-6373cf336a77",
		}

		service, err := client.CreateServiceInstance(req)
		So(err, ShouldBeNil)

		expected := ServiceInstance{
			Guid:        "8423ca96-90ad-411f-b77a-0907844949fc",
			CreatedAt:   "2016-10-21T18:22:56Z",
			UpdatedAt:   "2016-10-21T18:22:56Z",
			Credentials: map[string]interface{}{},
			Name:        "fortunes-db",
			LastOperation: LastOperation{
				Type:        "create",
				State:       "succeeded",
				Description: "",
				UpdatedAt:   "",
				CreatedAt:   "2016-10-21T18:22:56Z",
			},
			Tags:               []string{},
			ServiceGuid:        "440ce9d9-b108-4bbe-80b4-08338f3cc25b",
			ServicePlanGuid:    "f48419f7-4717-4706-86e4-a24973848a77",
			SpaceGuid:          "21e5fdc7-5131-4743-8447-6373cf336a77",
			DashboardUrl:       "https://p-mysql.system.example.com/manage/instances/8423ca96-90ad-411f-b77a-0907844949fc",
			Type:               "managed_service_instance",
			SpaceUrl:           "/v2/spaces/21e5fdc7-5131-4743-8447-6373cf336a77",
			ServicePlanUrl:     "/v2/service_plans/f48419f7-4717-4706-86e4-a24973848a77",
			ServiceBindingsUrl: "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_bindings",
			ServiceKeysUrl:     "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_keys",
			RoutesUrl:          "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/routes",
			ServiceUrl:         "/v2/services/440ce9d9-b108-4bbe-80b4-08338f3cc25b",
			c:                  client,
		}
		So(service, ShouldResemble, expected)
	})
}

func TestUpdateSIt(t *testing.T) {
	Convey("Update SI", t, func() {
		expectedPayload := "{\"name\":\"test-cs\",\"parameters\":{\"git\":{\"label\":\"master\",\"uri\":\"https://github.com/cloudfoundry-community/go-cfclient.git\"}},\"tags\":[\"tag1\",\"tag2\",\"tag3\"]}"
		req := ServiceInstanceUpdateRequest{
			Name: "test-cs",
			Tags: []string{"tag1", "tag2", "tag3"},
			Parameters: map[string]interface{}{
				"git": map[string]interface{}{
					"uri":   "https://github.com/cloudfoundry-community/go-cfclient.git",
					"label": "master",
				},
			},
		}

		setup(MockRoute{"PUT", "/v2/service_instances/guid", []string{""}, "", http.StatusAccepted, "accepts_incomplete=false", &expectedPayload}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.UpdateSI("guid", req, false)
		So(err, ShouldBeNil)
	})
}

func TestUpdateServiceInstance(t *testing.T) {
	Convey("Update service instance", t, func() {
		updateBody := "myUpdate"

		setup(MockRoute{"PUT", "/v2/service_instances/guid", []string{""}, "", http.StatusAccepted, "accepts_incomplete=false", &updateBody}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.UpdateServiceInstance("guid", strings.NewReader(updateBody), false)
		So(err, ShouldBeNil)
	})
}

func TestDeleteServiceInstance(t *testing.T) {
	Convey("Delete service instance", t, func() {
		setup(MockRoute{"DELETE", "/v2/service_instances/guid", []string{""}, "", http.StatusAccepted, "recursive=true&accepts_incomplete=false&async=false", nil}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteServiceInstance("guid", true, false)
		So(err, ShouldBeNil)
	})
}
