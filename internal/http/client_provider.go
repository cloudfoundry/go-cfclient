package http

import "net/http"

type ClientProvider interface {
	// Client returns a *http.Client
	Client() (*http.Client, error)

	// ReAuthenticate tells the provider to re-initialize the auth context
	ReAuthenticate() error
}

type UnauthenticatedClientProvider struct {
	httpClient *http.Client
}

func (c *UnauthenticatedClientProvider) Client() (*http.Client, error) {
	return c.httpClient, nil
}

func (c *UnauthenticatedClientProvider) ReAuthenticate() error {
	return nil
}

func NewUnauthenticatedClientProvider(httpClient *http.Client) *UnauthenticatedClientProvider {
	return &UnauthenticatedClientProvider{
		httpClient: httpClient,
	}
}
