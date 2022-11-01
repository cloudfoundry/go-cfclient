package client

import (
	"fmt"
	"net/url"
	"strings"
)

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
