package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type UserClient commonClient

// ListByQuery by query
func (c *UserClient) ListByQuery(query url.Values) ([]resource.User, error) {
	var users []resource.User
	requestURL, err := url.Parse("/v3/users")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.client.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.client.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error requesting users: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing users, response code: %d", resp.StatusCode)
		}

		var data resource.ListUsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list users: %w", err)
		}

		users = append(users, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, fmt.Errorf("error parsing next page URL: %w", err)
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}
