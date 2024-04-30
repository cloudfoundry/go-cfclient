package http_test

import (
	"bytes"
	"context"
	"golang.org/x/oauth2"
	gohttp "net/http"
	"testing"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedOAuthTokenSourceCreator struct {
	mock.Mock
}

func (tsc *MockedOAuthTokenSourceCreator) CreateOAuth2TokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	args := tsc.Called(ctx)
	return args.Get(0).(oauth2.TokenSource), args.Error(1)
}

type MockedOAuthTokenSource struct {
	mock.Mock
}

func (ts *MockedOAuthTokenSource) Token() (*oauth2.Token, error) {
	args := ts.Called()
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func TestOAuthSessionManager(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	serverURL := testutil.SetupMultiple([]testutil.MockRoute{
		{
			Method:    "POST",
			Endpoint:  "/v3/organizations",
			Output:    []string{"auth error", g.Organization().JSON},
			Statuses:  []int{401, 201},
			UserAgent: "Go-http-client/1.1",
		},
		{
			Method:    "GET",
			Endpoint:  "/v3/spaces",
			Output:    []string{"auth error", "spaces[]"},
			Statuses:  []int{401, 200},
			UserAgent: "Go-http-client/1.1",
		},
	}, t)
	defer testutil.Teardown()

	token := &oauth2.Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		Expiry:       time.Now().Add(time.Minute),
	}

	tokenSrc := &MockedOAuthTokenSource{}
	tokenSrc.On("Token").Return(token, nil)

	tokenSrcCreator := &MockedOAuthTokenSourceCreator{}
	tokenSrcCreator.On("CreateOAuth2TokenSource", context.Background()).Return(tokenSrc, nil)

	client, err := http.NewAuthenticatedClient(context.Background(), gohttp.DefaultClient, tokenSrcCreator)
	require.NoError(t, err)

	resp, err := client.Post(serverURL+"/v3/organizations", "application/json", bytes.NewBufferString(g.Organization().JSON))
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	// to the caller the retry is transparent on 401
	resp, err = client.Get(serverURL + "/v3/spaces")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	tokenSrcCreator.AssertNumberOfCalls(t, "CreateOAuth2TokenSource", 3)
}
