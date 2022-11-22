package http

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	reader := strings.NewReader("blah")
	object := time.Minute
	ctx := context.Background()

	r := NewRequest(ctx, "GET", "/v3/apps")
	require.Nil(t, r.contentLength, "content length should default to nil")
	require.Empty(t, r.contentType, "content type defaults to empty")
	require.NotNil(t, r.headers, "headers should default to empty, not nil")

	r = NewRequest(ctx, "GET", "/v3/apps").WithObject(object)
	require.Equal(t, "application/json", r.contentType, "with object, content type should be json")
	require.Equal(t, object, r.object)

	r = NewRequest(ctx, "GET", "/v3/apps").WithBody(reader)
	require.Equal(t, "application/json", r.contentType, "with body, content type should be json")
	require.Equal(t, reader, r.body)

	r = NewRequest(ctx, "GET", "/v3/apps").WithContentType("application/yaml")
	require.Equal(t, "application/yaml", r.contentType)
	r.WithBody(reader)
	require.Equal(t, "application/yaml", r.contentType,
		"with content type previously set, setting a body shouldn't override it")
	require.Equal(t, reader, r.body)
}
