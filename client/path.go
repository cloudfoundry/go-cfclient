package client

import (
	"fmt"
	"net/url"
	"strings"
)

func extractPathFromURL(requestURL string) (string, error) {
	u, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	result := u.Path
	if q := u.Query().Encode(); q != "" {
		result = result + "?" + q
	}
	return result, nil
}

func path(urlFormat string, params ...any) string {
	// url encode any querystring params
	p := make([]any, len(params))
	for i, u := range params {
		switch v := u.(type) {
		case url.Values:
			p[i] = v.Encode()
		default:
			p[i] = u
		}
	}

	s := fmt.Sprintf(urlFormat, p...)
	return strings.TrimSuffix(s, "?")
}
