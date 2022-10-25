package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListServiceInstancesByQuery(t *testing.T) {
	setup(MockRoute{"GET", "/v3/service_instances", []string{listServiceInstancesPayload}, "", http.StatusOK, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	services, err := client.ServiceInstances.List()
	require.NoError(t, err)
	require.Len(t, services, 1)

	require.Equal(t, "my_service_instance", services[0].Name)

	require.Equal(t, "ae0031f9-dd49-461c-a945-df40e77c39cb", services[0].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/spaces/ae0031f9-dd49-461c-a945-df40e77c39cb", services[0].Links["space"].Href)
}
