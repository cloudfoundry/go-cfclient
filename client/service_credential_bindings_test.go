package client

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServiceCredentialBindingsByQuery(t *testing.T) {
	Convey("List  Service Credential Bindings", t, func() {
		setup(MockRoute{"GET", "/v3/service_credential_bindings", []string{listServiceCredentialBindingsPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c, _ := NewTokenConfig(server.URL, "foobar")
		client, err := New(c)
		So(err, ShouldBeNil)

		serviceCredentialsBindings, err := client.ListServiceCredentialBindings()
		So(err, ShouldBeNil)
		So(serviceCredentialsBindings, ShouldHaveLength, 1)

		So(serviceCredentialsBindings[0].Name, ShouldEqual, "my_service_key")
		So(serviceCredentialsBindings[0].Type, ShouldEqual, "key")

		So(serviceCredentialsBindings[0].Relationships["service_instance"].Data.GUID, ShouldEqual, "85ccdcad-d725-4109-bca4-fd6ba062b5c8")
		So(serviceCredentialsBindings[0].Links["service_instance"].Href, ShouldEqual, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8")
	})
}

func TestGetServiceCredentialBindingsByGUID(t *testing.T) {
	Convey("Get  Service Credential Binding by GUID", t, func() {
		GUID := "d9634934-8e1f-4c2d-bb33-fa5df019cf9d"
		setup(MockRoute{"GET", fmt.Sprintf("/v3/service_credential_bindings/%s", GUID), []string{getServiceCredentialBindingsByGUIDPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c, _ := NewTokenConfig(server.URL, "foobar")
		client, err := New(c)
		So(err, ShouldBeNil)

		serviceCredentialsBinding, err := client.GetServiceCredentialBindingsByGUID(GUID)
		So(err, ShouldBeNil)

		So(serviceCredentialsBinding.Name, ShouldEqual, "my_service_key")
		So(serviceCredentialsBinding.Type, ShouldEqual, "key")

		So(serviceCredentialsBinding.Relationships["service_instance"].Data.GUID, ShouldEqual, "85ccdcad-d725-4109-bca4-fd6ba062b5c8")
		So(serviceCredentialsBinding.Links["service_instance"].Href, ShouldEqual, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8")
	})
}
