package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/ios"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

// OAuthTokenSourceCreator implementations create OAuth2 TokenSources
type OAuthTokenSourceCreator interface {
	// CreateOAuth2TokenSource creates a new OAuth2 TokenSource when called
	CreateOAuth2TokenSource(ctx context.Context) (oauth2.TokenSource, error)
}

// retryableAuthTransport wraps a http.RoundTripper and combines it with an OAuthTokenSourceCreator so
// that any 401s cause a re-authentication and request retry
type retryableAuthTransport struct {
	transport          http.RoundTripper
	tokenSourceCreator OAuthTokenSourceCreator
}

// NewAuthenticatedClient creates a new http.Client with a retryableAuthTransport that supports re-authentication
// and request retry should a request cause a 401.
func NewAuthenticatedClient(ctx context.Context, baseClient *http.Client, tokenSourceCreator OAuthTokenSourceCreator) (*http.Client, error) {
	src, err := tokenSourceCreator.CreateOAuth2TokenSource(ctx)
	if err != nil {
		return nil, err
	}

	transport := &retryableAuthTransport{
		tokenSourceCreator: tokenSourceCreator,
		transport: &oauth2.Transport{
			Base:   baseClient.Transport,
			Source: src,
		},
	}

	// oauth2.NewClient only copies the transport, so explicitly create our own http.Client
	// https://github.com/golang/oauth2/issues/368
	return &http.Client{
		Transport:     transport,
		Timeout:       baseClient.Timeout,
		CheckRedirect: CheckRedirect,
		Jar:           baseClient.Jar,
	}, nil
}

func (t *retryableAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Send the request
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Retry logic
	for shouldRetryAuth(resp) {
		// We're going to retry, consume any response to reuse the connection.
		drainBody(resp)

		// Clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Recreate the token source
		src, tsErr := t.tokenSourceCreator.CreateOAuth2TokenSource(req.Context())
		if tsErr != nil {
			return nil, fmt.Errorf("error re-authenticating with the OAuth2 token source: %w", tsErr)
		}
		t.transport.(*oauth2.Transport).Source = src

		// Retry the request
		resp, err = t.transport.RoundTrip(req)
	}

	// Return the response
	return resp, err
}

func shouldRetryAuth(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusUnauthorized
}

func drainBody(resp *http.Response) {
	if resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		ios.Close(resp.Body)
	}
}
