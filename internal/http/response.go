package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloudfoundry/go-cfclient/v3/internal/ios"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

// DecodeJobIDAndBody returns the jobGUID if specified in the Location response header and
// unmarshalls the JSON response body to result if available
func DecodeJobIDAndBody(resp *http.Response, result any) (string, error) {
	return DecodeJobID(resp), DecodeBody(resp, result)
}

// DecodeJobIDOrBody returns the jobGUID if specified in the Location response header or
// unmarshalls the JSON response body if no job ID and result is non nil
func DecodeJobIDOrBody(resp *http.Response, result any) (string, error) {
	if jobGUID := DecodeJobID(resp); jobGUID != "" {
		return jobGUID, nil
	}
	return "", DecodeBody(resp, result)
}

// DecodeJobID returns the jobGUID if specified in the Location response header
func DecodeJobID(resp *http.Response) string {
	// Return empty if the response is nil
	if resp == nil {
		return ""
	}

	// Extract the Location header
	location, err := resp.Location()
	if err != nil {
		// Return empty if there's an error (e.g., no Location header)
		return ""
	}

	// Split the path in the URL and check for the 'jobs' segment
	parts := strings.Split(location.Path, "/")
	numParts := len(parts)
	// Ensure 'jobs' is the second last element and return the last element as job ID
	if numParts >= 2 && parts[numParts-2] == "jobs" {
		return parts[numParts-1]
	}

	// Return empty if the URL doesn't follow the expected pattern
	return ""
}

// DecodeBody unmarshalls the JSON response body if the result is non nil
func DecodeBody(resp *http.Response, result any) error {
	if result == nil || resp == nil || resp.Body == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil && err != io.EOF {
		return fmt.Errorf("error decoding response JSON: %w", err)
	}
	return nil
}

func DecodeError(resp *http.Response) error {
	if resp == nil || resp.Body == nil {
		return errors.New("response has empty or invalid body")
	}

	defer ios.Close(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var cfErrs resource.CloudFoundryErrors
		if err = json.Unmarshal(body, &cfErrs); err == nil && len(cfErrs.Errors) > 0 {
			return cfErrs.Errors[0]
		}
	}
	return resource.CloudFoundryHTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}
}
