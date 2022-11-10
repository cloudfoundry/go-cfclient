package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	defer Teardown()
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	c.SkipSSLValidation(true)
	client, err := New(c)
	require.NoError(t, err)
	req := client.NewRequest("GET", "/v2/organizations")
	resp, err := client.DoRequest(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestMakeRequestFailure(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	defer Teardown()
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	c.SkipSSLValidation(true)
	client, err := New(c)
	require.NoError(t, err)
	req := client.NewRequest("GET", "/v2/organizations")
	req.url = "%gh&%ij"
	resp, err := client.DoRequest(req)
	require.Nil(t, resp)
	require.NotNil(t, err)
}

func TestMakeRequestWithTimeout(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	defer Teardown()
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	c.SkipSSLValidation(true)
	c.HTTPClient(&http.Client{Timeout: 10 * time.Nanosecond})
	client, err := New(c)
	require.NotNil(t, err)
	require.Nil(t, client)
}

func TestHTTPErrorHandling(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"502 Bad Gateway"}, Status: 502}, t)
	defer Teardown()
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	c.SkipSSLValidation(true)
	client, err := New(c)
	require.NoError(t, err)
	req := client.NewRequest("GET", "/v2/organizations")
	resp, err := client.DoRequest(req)
	require.NotNil(t, err)
	require.NotNil(t, resp)

	httpErr := err.(CloudFoundryHTTPError)
	require.Equal(t, 502, httpErr.StatusCode)
	require.Equal(t, "502 Bad Gateway", httpErr.Status)
	require.Equal(t, "502 Bad Gateway", string(httpErr.Body))
}

func TestTokenRefresh(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	fakeUAAServer = FakeUAAServer(1)
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	client, err := New(c)
	require.NoError(t, err)

	token, err := client.GetToken()
	require.NoError(t, err)
	require.Equal(t, "bearer foobar2", token)

	for i := 0; i < 5; i++ {
		token, _ = client.GetToken()
		if token == "bearer foobar3" {
			break
		}
		time.Sleep(time.Second)
	}
	require.Equal(t, "bearer foobar3", token)
}

func TestEndpointRefresh(t *testing.T) {
	Setup(MockRoute{Method: "GET", Endpoint: "/v2/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	fakeUAAServer = FakeUAAServer(0)
	c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
	client, err := New(c)
	require.NoError(t, err)

	//lastTokenSource := client.Config.TokenSource
	for i := 1; i < 5; i++ {
		_, err := client.GetToken()
		require.NoError(t, err)
		//So(client.Config.TokenSource, ShouldNotEqual, lastTokenSource)
		//lastTokenSource = client.Config.TokenSource
	}
}
