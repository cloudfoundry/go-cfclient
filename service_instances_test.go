package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceInstanceByGuid(t *testing.T) {
	Convey("Service instance by Guid", t, func() {
		setup(MockRoute{"GET", "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc", serviceInstancePayload, ""}, t)
		defer teardown()

		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		service, err := client.ServiceInstanceByGuid("8423ca96-90ad-411f-b77a-0907844949fc")
		So(err, ShouldBeNil)

		expected := ServiceInstance{
			Guid:               "8423ca96-90ad-411f-b77a-0907844949fc",
			Name:               "fortunes-db",
			ServicePlanGuid:    "f48419f7-4717-4706-86e4-a24973848a77",
			SpaceGuid:          "21e5fdc7-5131-4743-8447-6373cf336a77",
			DashboardUrl:       "https://p-mysql.system.example.com/manage/instances/8423ca96-90ad-411f-b77a-0907844949fc",
			Type:               "managed_service_instance",
			SpaceUrl:           "/v2/spaces/21e5fdc7-5131-4743-8447-6373cf336a77",
			ServicePlanUrl:     "/v2/service_plans/f48419f7-4717-4706-86e4-a24973848a77",
			ServiceBindingsUrl: "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_bindings",
			ServiceKeysUrl:     "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_keys",
			RoutesUrl:          "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/routes",
			c:                  client,
		}
		So(service, ShouldResemble, expected)
	})
}
