package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

type Executor struct {
	userAgent      string
	apiAddress     string
	clientProvider ClientProvider
}

func NewExecutor(clientProvider ClientProvider, apiAddress, userAgent string) *Executor {
	return &Executor{
		userAgent:      userAgent,
		apiAddress:     apiAddress,
		clientProvider: clientProvider,
	}
}

func (c *Executor) ExecuteRequest(request *Request) (*http.Response, error) {
	reqBody := request.body
	if request.object != nil {
		b, err := encodeBody(request.object)
		if err != nil {
			return nil, fmt.Errorf("error executing request, failed to encode the request object to JSON: %w", err)
		}
		reqBody = b
	}
	u := path.Join(c.apiAddress, request.pathAndQuery)

	req, err := http.NewRequest(request.method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed to create a new underlying HTTP request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	if request.contentType != "" {
		req.Header.Set("Content-type", request.contentType)
	}
	if request.contentLength != nil {
		req.ContentLength = *request.contentLength
	}

	client, err := c.clientProvider.Client()
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed to get the underlying HTTP client: %w", err)
	}
	return client.Do(req)
}

// encodeBody is used to encode a request body
func encodeBody(obj any) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}
