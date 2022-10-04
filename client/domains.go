package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/pkg/errors"
)

func (c *Client) ListDomains(query url.Values) ([]resource.Domain, error) {
	var domains []resource.Domain
	requestURL := "/v3/domains"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.DoRequest(c.NewRequest("GET", requestURL))
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting domains")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 app domains, response code: %d", resp.StatusCode)
		}

		var data resource.ListDomainsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 app domains")
		}

		domains = append(domains, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 domains")
		}
	}
	return domains, nil
}
