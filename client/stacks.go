package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type StackClient commonClient

// ListByQuery retrieves stacks based on query
func (c *StackClient) ListByQuery(query url.Values) ([]resource.Stack, error) {
	var stacks []resource.Stack
	requestURL, err := url.Parse("/v3/stacks")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.client.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting stacks: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing stacks, response code: %d", resp.StatusCode)
		}

		var data resource.ListStacksResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list stacks: %w", err)
		}

		stacks = append(stacks, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, fmt.Errorf("error parsing next page URL: %w", err)
		}
		if requestURL.String() == "" {
			break
		}
	}

	return stacks, nil
}
