package http_test

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"io"
	http2 "net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteRequest(t *testing.T) {
	serverURL := testutil.Setup(testutil.MockRoute{Method: "GET", Endpoint: "/v3/organizations", Output: []string{"fake payload"}, Status: 200}, t)
	defer testutil.Teardown()

	httpClient := &http2.Client{
		Transport: http2.DefaultTransport,
	}
	clientProvider := http.NewUnauthenticatedClientProvider(httpClient)
	e := http.NewExecutor(clientProvider, serverURL, config.UserAgent)
	req := http.NewRequest(context.Background(), "GET", "/v3/organizations")
	r, err := e.ExecuteRequest(req)
	require.NoError(t, err)
	require.Equal(t, 200, r.StatusCode)
}

func TestExecuteRequestWithAuthFailure(t *testing.T) {
	// use a custom httptest server that we can call the same endpoint but get different status codes each time
	callCount := 0
	server := httptest.NewServer(http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		if callCount == 0 {
			w.WriteHeader(401)
		} else {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("success"))
		}
		callCount++

	}))
	defer server.Close()

	httpClient := &http2.Client{
		Transport: http2.DefaultTransport,
	}
	clientProvider := http.NewUnauthenticatedClientProvider(httpClient)
	e := http.NewExecutor(clientProvider, server.URL, config.UserAgent)
	req := http.NewRequest(context.Background(), "GET", "/does_not_matter")
	r, err := e.ExecuteRequest(req)
	require.NoError(t, err)
	require.Equal(t, 200, r.StatusCode)
	require.Equal(t, 2, callCount)
}

func TestExecuteRequestCreatesHTTPRequest(t *testing.T) {
	type testObj struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	obj := &testObj{
		ID:   12,
		Name: "Shenanigans",
	}
	expectedJSON := `{"id":12,"name":"Shenanigans"}
`

	// use a custom httptest server so that we can validate the server was passed all the right info via the request
	server := httptest.NewServer(http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, "bar", r.Header.Get("foo"))
		require.NotNil(t, r.Body)
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, expectedJSON, string(b))

		w.WriteHeader(200)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	httpClient := &http2.Client{
		Transport: http2.DefaultTransport,
	}
	clientProvider := http.NewUnauthenticatedClientProvider(httpClient)
	e := http.NewExecutor(clientProvider, server.URL, config.UserAgent)
	req := http.NewRequest(context.Background(), "POST", "/does_not_matter").WithObject(obj).WithHeader("foo", "bar")
	r, err := e.ExecuteRequest(req)
	require.NoError(t, err)
	require.Equal(t, 200, r.StatusCode)
}
