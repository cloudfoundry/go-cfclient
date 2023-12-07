package http_test

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
	"time"

	gohttp "net/http"
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
	serverURL := testutil.SetupMultiple([]testutil.MockRoute{
		{
			Method:    "GET",
			Endpoint:  "/v3/organizations",
			Output:    []string{"organizations[]"},
			Statuses:  []int{200},
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

	resp, err := client.Get(serverURL + "/v3/organizations")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	// to the caller the retry is transparent on 401
	resp, err = client.Get(serverURL + "/v3/spaces")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	tokenSrcCreator.AssertNumberOfCalls(t, "CreateOAuth2TokenSource", 2)
}
