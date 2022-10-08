package client

import (
	"net/url"
	"strings"
)

func joinPathAndQS(qs url.Values, pathParts ...string) string {
	u := url.URL{
		Path: joinPath(pathParts...),
	}
	u.RawQuery = qs.Encode()
	return u.String()
}

func joinPath(pathParts ...string) string {
	return strings.Join(pathParts, "/")
}

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
