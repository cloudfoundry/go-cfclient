package http

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"

	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	t.Run("Test DecodeJobID", func(t *testing.T) {
		// Test with nil response
		require.Equal(t, "", DecodeJobID(nil))

		// Test with empty response
		resp := &http.Response{}
		require.Equal(t, "", DecodeJobID(resp))

		// Test with valid Location header
		resp.Header = http.Header{"Location": []string{"https://test.cf.com/jobs/jobGUID"}}
		require.Equal(t, "jobGUID", DecodeJobID(resp))
	})

	t.Run("Test DecodeBody", func(t *testing.T) {
		// Test with nil parameters
		require.Nil(t, DecodeBody(nil, nil))

		var body struct {
			Data string `json:"data,omitempty"`
		}

		// Test with no content status
		resp := &http.Response{StatusCode: http.StatusNoContent}
		require.Nil(t, DecodeBody(resp, &body))

		// Test with OK status and invalid JSON
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("testData")),
		}
		require.NotNil(t, DecodeBody(resp, &body))
		require.Equal(t, "", body.Data)

		// Test with OK status and valid JSON
		resp.Body = io.NopCloser(strings.NewReader(`{"data": "test"}`))
		require.Nil(t, DecodeBody(resp, &body))
		require.Equal(t, "test", body.Data)
	})

	t.Run("Test DecodeJobIDOrBody and DecodeJobIDAndBody", func(t *testing.T) {
		resp := &http.Response{
			Header: http.Header{"Location": []string{"https://test.cf.com/jobs/jobGUID"}},
		}

		var body struct {
			Data string `json:"data,omitempty"`
		}

		// Test DecodeJobIDOrBody
		jobGuid, err := DecodeJobIDOrBody(resp, &body)
		require.Equal(t, "jobGUID", jobGuid)
		require.Nil(t, err)

		// Test DecodeJobIDAndBody with invalid JSON
		resp.Body = io.NopCloser(strings.NewReader("testData"))
		jobGuid, err = DecodeJobIDAndBody(resp, &body)
		require.Equal(t, "jobGUID", jobGuid)
		require.NotNil(t, err)
		require.Equal(t, "", body.Data)

		// Test DecodeJobIDAndBody with valid JSON
		resp.Body = io.NopCloser(strings.NewReader(`{"data": "test"}`))
		jobGuid, err = DecodeJobIDAndBody(resp, &body)
		require.Equal(t, "jobGUID", jobGuid) // Expecting jobGuid to be empty due to updated Location header
		require.Nil(t, err)
		require.Equal(t, "test", body.Data)

	})

	t.Run("Test DecodeError", func(t *testing.T) {
		err := DecodeError(nil)
		require.EqualError(t, err, "response has empty or invalid body")

		resp := &http.Response{Body: io.NopCloser(strings.NewReader(`{"errors":[{"detail":"Unknown request","title":"CF-NotFound","code":10000}]}`))}
		err = DecodeError(resp)
		require.IsType(t, resource.CloudFoundryError{}, err)
		require.EqualError(t, err, "cfclient error (CF-NotFound|10000): Unknown request")

		resp = &http.Response{Body: io.NopCloser(strings.NewReader(`invalid request`)), StatusCode: 404, Status: "Not Found"}
		err = DecodeError(resp)
		require.IsType(t, resource.CloudFoundryHTTPError{}, err)
		require.EqualError(t, err, "cfclient: HTTP error (404): Not Found")
	})
}
