package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	MaxRedirects    = 10
	ErrMaxRedirects = "stopped after maximum allowed redirects"
)

// contextKey is a private static type to avoid potential collisions.
type contextKey struct{}

// ignoreRedirectKey is a context key used for storing the redirect ignore flag.
var ignoreRedirectKey = contextKey{}

// IgnoreRedirect sets a flag in the request's context to indicate that redirects should be ignored.
func IgnoreRedirect(req *http.Request) *http.Request {
	if req == nil {
		return nil
	}
	return req.WithContext(context.WithValue(req.Context(), ignoreRedirectKey, true))
}

// IsIgnoredRedirect checks if the 'ignore redirect' flag is set in the request's context.
func IsIgnoredRedirect(req *http.Request) bool {
	if req == nil {
		return false
	}
	v, ok := req.Context().Value(ignoreRedirectKey).(bool)
	return ok && v
}

// CheckRedirect checks the redirect policy for the HTTP client.
func CheckRedirect(req *http.Request, via []*http.Request) error {
	if IsIgnoredRedirect(req) {
		return http.ErrUseLastResponse
	}
	if len(via) >= MaxRedirects {
		return errors.New(ErrMaxRedirects)
	}
	return nil
}

func EncodeBody(obj any) (io.Reader, error) {
	if obj == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(obj); err != nil {
		return nil, fmt.Errorf("error encoding object to JSON: %w", err)
	}
	return buf, nil
}
