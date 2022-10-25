package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestGetDeployment(t *testing.T) {
	setup(MockRoute{"GET", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4", []string{getDeploymentPayload}, "", http.StatusOK, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Deployments.Get("59c3d133-2b83-46f3-960e-7765a129aea4")
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "59c3d133-2b83-46f3-960e-7765a129aea4", resp.GUID)
	require.Equal(t, "DEPLOYING", resp.Status.Reason)
	require.Equal(t, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2", resp.Relationships["app"].Data.GUID)
}

func TestCreateDeployment(t *testing.T) {
	body := `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
	setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", body}, t)

	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Deployments.Create("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "59c3d133-2b83-46f3-960e-7765a129aea4", resp.GUID)
	require.Equal(t, "DEPLOYING", resp.Status.Reason)
	require.Equal(t, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2", resp.Relationships["app"].Data.GUID)
}

func TestCreateDeploymentWithDroplet(t *testing.T) {
	body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
	setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", body}, t)

	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Deployments.Create("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &resource.CreateDeploymentOptionalParameters{
		Droplet: &resource.Relationship{
			GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "59c3d133-2b83-46f3-960e-7765a129aea4", resp.GUID)
	require.Equal(t, "DEPLOYING", resp.Status.Reason)
	require.Equal(t, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2", resp.Relationships["app"].Data.GUID)
}
func TestCreateDeploymentWithRevision(t *testing.T) {
	body := `{"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
	setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", body}, t)

	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Deployments.Create("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &resource.CreateDeploymentOptionalParameters{
		Revision: &resource.DeploymentRevision{
			GUID:    "56126cba-656a-4eba-a81e-7e9951b2df57",
			Version: 1,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "59c3d133-2b83-46f3-960e-7765a129aea4", resp.GUID)
	require.Equal(t, "DEPLOYING", resp.Status.Reason)
	require.Equal(t, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2", resp.Relationships["app"].Data.GUID)
}

func TestCreateDeploymentWithRevisionAndDroplet(t *testing.T) {
	body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
	setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", body}, t)

	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Deployments.Create("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &resource.CreateDeploymentOptionalParameters{
		Droplet: &resource.Relationship{
			GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
		},
		Revision: &resource.DeploymentRevision{
			GUID:    "56126cba-656a-4eba-a81e-7e9951b2df57",
			Version: 1,
		},
	})
	require.NotNil(t, err)
	require.Nil(t, resp)
}

func TestCancelDeployment(t *testing.T) {
	setup(MockRoute{"POST", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4/actions/cancel", []string{""}, "", http.StatusOK, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Deployments.Cancel("59c3d133-2b83-46f3-960e-7765a129aea4")
	require.NoError(t, err)
}
