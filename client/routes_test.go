package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListRoutes(t *testing.T) {
	setup(MockRoute{"GET", "/v3/routes", []string{listRoutesPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	routes, err := client.Routes.List()
	require.NoError(t, err)
	require.Len(t, routes, 1)

	require.Equal(t, "a-hostname", routes[0].Host)
	require.Equal(t, "/some_path", routes[0].Path)
	require.Equal(t, "a-hostname.a-domain.com/some_path", routes[0].Url)

	require.Equal(t, "885a8cb3-c07b-4856-b448-eeb10bf36236", routes[0].Relationships["space"].Data.GUID)
	require.Equal(t, "0b5f3633-194c-42d2-9408-972366617e0e", routes[0].Relationships["domain"].Data.GUID)

	require.Equal(t, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31", routes[0].Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236", routes[0].Links["space"].Href)
	require.Equal(t, "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e", routes[0].Links["domain"].Href)
	require.Equal(t, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations", routes[0].Links["destinations"].Href)
}

func TestCreateRoutes(t *testing.T) {
	setup(MockRoute{"POST", "/v3/routes", []string{createRoutePayload}, "", http.StatusCreated, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	route, err := client.Routes.Create(
		"885a8cb3-c07b-4856-b448-eeb10bf36236",
		"0b5f3633-194c-42d2-9408-972366617e0e",
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "a-hostname", route.Host)
	require.Equal(t, "/some_path", route.Path)
	require.Equal(t, "885a8cb3-c07b-4856-b448-eeb10bf36236", route.Relationships["space"].Data.GUID)
	require.Equal(t, "0b5f3633-194c-42d2-9408-972366617e0e", route.Relationships["domain"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31", route.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236", route.Links["space"].Href)
	require.Equal(t, "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e", route.Links["domain"].Href)
	require.Equal(t, "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations", route.Links["destinations"].Href)
}
