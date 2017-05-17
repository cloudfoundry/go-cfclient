package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServiceKeys(t *testing.T) {
	Convey("List Service Keys", t, func() {
		setup(MockRoute{"GET", "/v2/service_keys", listServiceKeysPayload, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		serviceKeys, err := client.ListServiceKeys()
		So(err, ShouldBeNil)

		So(len(serviceKeys), ShouldEqual, 2)
		So(serviceKeys[0].Guid, ShouldEqual, "3b933598-64ed-4613-a0f5-b7e8c0379368")
		So(serviceKeys[0].Name, ShouldEqual, "RedisMonitoringKey")
		So(serviceKeys[0].ServiceInstanceGuid, ShouldEqual, "ad98f310-a3a0-47aa-9116-f8295d41a9b2")
		So(serviceKeys[0].Credentials, ShouldNotEqual, nil)
		So(serviceKeys[0].ServiceInstanceUrl, ShouldEqual, "/v2/service_instances/ad98f310-a3a0-47aa-9116-f8295d41a9b2")
		So(serviceKeys[1].Guid, ShouldEqual, "8be3911b-c621-4467-8866-f8b924aaee57")
		So(serviceKeys[1].Name, ShouldEqual, "test01_key")
		So(serviceKeys[1].ServiceInstanceGuid, ShouldEqual, "ecf26687-e176-4784-b181-b3c942fecb62")
		So(serviceKeys[1].Credentials, ShouldNotEqual, nil)
		m := serviceKeys[1].Credentials.(map[string]interface{})
		So(m["uri"], ShouldEqual, "nhp://100.100.100.100:9008")
		So(serviceKeys[1].ServiceInstanceUrl, ShouldEqual, "/v2/service_instances/fcf26687-e176-4784-b181-b3c942fecb62")
	})
}
