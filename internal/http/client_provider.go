package http

import "net/http"

type ClientProvider interface {
	Client() (*http.Client, error)
}

type UnauthenticatedClientProvider struct {
	httpClient *http.Client
}

func (c *UnauthenticatedClientProvider) Client() (*http.Client, error) {
	return c.httpClient, nil
}

func NewUnauthenticatedClientProvider(httpClient *http.Client) *UnauthenticatedClientProvider {
	return &UnauthenticatedClientProvider{
		httpClient: httpClient,
	}
}
