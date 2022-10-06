package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSetCurrentDropletForApp(t *testing.T) {
	body := `{"data":{"guid":"3fc0916f-2cea-4f3a-ae53-048388baa6bd"}}`
	setup(MockRoute{"PATCH", "/v3/apps/9d8e007c-ce52-4ea7-8a57-f2825d2c6b39/relationships/current_droplet", []string{currentDropletV3AppPayload}, "", http.StatusOK, "", &body}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Droplets.SetCurrentForApp("9d8e007c-ce52-4ea7-8a57-f2825d2c6b39", "3fc0916f-2cea-4f3a-ae53-048388baa6bd")
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "9d8e007c-ce52-4ea7-8a57-f2825d2c6b39", resp.Data.GUID)
	require.Equal(t, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/relationships/current_droplet", resp.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/droplets/current", resp.Links["related"].Href)
}

func TestGetCurrentDropletForApp(t *testing.T) {
	setup(MockRoute{"GET", "/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396/droplets/current", []string{getCurrentAppDropletPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Droplets.GetCurrentForApp("7b34f1cf-7e73-428a-bb5a-8a17a8058396")
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "585bc3c1-3743-497d-88b0-403ad6b56d16", resp.GUID)
	require.Equal(t, "https://api.example.org/v3/droplets/585bc3c1-3743-497d-88b0-403ad6b56d16", resp.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396/relationships/current_droplet", resp.Links["assign_current_droplet"].Href)
	require.Equal(t, "PATCH", resp.Links["assign_current_droplet"].Method)
}

func TestDeleteDroplet(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Droplets.Delete("59c3d133-2b83-46f3-960e-7765a129aea4")
	require.NoError(t, err)
}
