package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"io"
	web "net/http"
)

// RootClient queries the global API root /
type RootClient struct {
	httpExecutor *http.Executor
}

// NewRootClient creates an initialized root client
func NewRootClient(httpExecutor *http.Executor) *RootClient {
	return &RootClient{
		httpExecutor: httpExecutor,
	}
}

// Get queries the global API root /
//
// These endpoints link to other resources, endpoints, and external services that are relevant to
// authenticated API clients.
func (c *RootClient) Get() (*resource.Root, error) {
	req := http.NewRequest("GET", "/")
	res, err := c.httpExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(res.Body)
	if res.StatusCode != web.StatusOK {
		return nil, fmt.Errorf("error getting global API root, got status code %d", res.StatusCode)
	}

	//var sb strings.Builder
	//_, _ = io.Copy(&sb, res.Body)
	//s := sb.String()
	//fmt.Println(s)

	var root resource.Root
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&root)
	return &root, err
}
