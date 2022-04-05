package cfclient

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3ServiceCredentialBindingsByQuery(t *testing.T) {
	Convey("List V3 Service Credential Bindings", t, func() {
		setup(MockRoute{"GET", "/v3/service_credential_bindings", []string{listV3ServiceCredentialBindingsPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		serviceCredentialsBindings, err := client.ListV3ServiceCredentialBindings()
		So(err, ShouldBeNil)
		So(serviceCredentialsBindings, ShouldHaveLength, 1)

		So(serviceCredentialsBindings[0].Name, ShouldEqual, "my_service_key")
		So(serviceCredentialsBindings[0].Type, ShouldEqual, "key")

		So(serviceCredentialsBindings[0].Relationships["service_instance"].Data.GUID, ShouldEqual, "85ccdcad-d725-4109-bca4-fd6ba062b5c8")
		So(serviceCredentialsBindings[0].Links["service_instance"].Href, ShouldEqual, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8")
	})
}

func TestGetV3ServiceCredentialBindingsByGUID(t *testing.T) {
	Convey("Get V3 Service Credential Binding by GUID", t, func() {
		GUID := "d9634934-8e1f-4c2d-bb33-fa5df019cf9d"
		setup(MockRoute{"GET", fmt.Sprintf("/v3/service_credential_bindings/%s", GUID), []string{GetV3ServiceCredentialBindingsByGUIDPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		serviceCredentialsBinding, err := client.GetV3ServiceCredentialBindingsByGUID(GUID)
		So(err, ShouldBeNil)

		So(serviceCredentialsBinding.Name, ShouldEqual, "my_service_key")
		So(serviceCredentialsBinding.Type, ShouldEqual, "key")

		So(serviceCredentialsBinding.Relationships["service_instance"].Data.GUID, ShouldEqual, "85ccdcad-d725-4109-bca4-fd6ba062b5c8")
		So(serviceCredentialsBinding.Links["service_instance"].Href, ShouldEqual, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8")
	})
}
