package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

// ListUsersByQuery by query
func (c *Client) ListUsersByQuery(query url.Values) ([]resource.User, error) {
	var users []resource.User
	requestURL, err := url.Parse("/v3/users")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting  users")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing  users, response code: %d", resp.StatusCode)
		}

		var data resource.ListUsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list  users")
		}

		users = append(users, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}
