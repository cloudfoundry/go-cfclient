package client

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListServiceCredentialBindingsByQuery(t *testing.T) {
	setup(MockRoute{"GET", "/v3/service_credential_bindings", []string{listServiceCredentialBindingsPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	serviceCredentialsBindings, err := client.ServiceCredentialBindings.List()
	require.NoError(t, err)
	require.Len(t, serviceCredentialsBindings, 1)

	require.Equal(t, "my_service_key", serviceCredentialsBindings[0].Name)
	require.Equal(t, "key", serviceCredentialsBindings[0].Type)

	require.Equal(t, "85ccdcad-d725-4109-bca4-fd6ba062b5c8", serviceCredentialsBindings[0].Relationships["service_instance"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8", serviceCredentialsBindings[0].Links["service_instance"].Href)
}

func TestGetServiceCredentialBindingsByGUID(t *testing.T) {
	GUID := "d9634934-8e1f-4c2d-bb33-fa5df019cf9d"
	setup(MockRoute{"GET", fmt.Sprintf("/v3/service_credential_bindings/%s", GUID), []string{getServiceCredentialBindingsByGUIDPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	serviceCredentialsBinding, err := client.ServiceCredentialBindings.Get(GUID)
	require.NoError(t, err)

	require.Equal(t, "my_service_key", serviceCredentialsBinding.Name)
	require.Equal(t, "key", serviceCredentialsBinding.Type)

	require.Equal(t, "85ccdcad-d725-4109-bca4-fd6ba062b5c8", serviceCredentialsBinding.Relationships["service_instance"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8", serviceCredentialsBinding.Links["service_instance"].Href)
}
