package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type DomainClient commonClient

func (c *DomainClient) ListByQuery(query url.Values) ([]resource.Domain, error) {
	var domains []resource.Domain
	requestURL := "/v3/domains"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.client.DoRequest(c.client.NewRequest("GET", requestURL))
		if err != nil {
			return nil, fmt.Errorf("error getting domains: %w", err)
		}
		defer func(b io.ReadCloser) {
			_ = b.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing v3 app domains, response code: %d", resp.StatusCode)
		}

		var data resource.ListDomainsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error parsing JSON from list v3 app domains: %w", err)
		}

		domains = append(domains, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing the next page request url for v3 domains: %w", err)
		}
	}
	return domains, nil
}
