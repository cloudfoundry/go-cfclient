//go:generate go run ../tools/gen_error.go

package resource

import (
	"fmt"
	"strings"
)

type CloudFoundryHTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e CloudFoundryHTTPError) Error() string {
	return fmt.Sprintf("cfclient: HTTP error (%d): %s", e.StatusCode, e.Status)
}

type CloudFoundryErrors struct {
	Errors []CloudFoundryError `json:"errors"`
}

func (e CloudFoundryErrors) Error() string {
	var sb strings.Builder
	for _, err := range e.Errors {
		sb.WriteString(err.Error())
	}
	return sb.String()
}

type CloudFoundryError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e CloudFoundryError) Error() string {
	return fmt.Sprintf("cfclient error (%s|%d): %s", e.Title, e.Code, e.Detail)
}
