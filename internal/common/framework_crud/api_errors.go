package framework_crud

import (
	"errors"
	"net/http"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
)

// APIStatusCode returns a wrapped SDK API error's HTTP status, or zero for a transport error.
func APIStatusCode(err error) int {
	var apiErr *apiresponse.APIError
	if errors.As(err, &apiErr) && apiErr != nil {
		return apiErr.StatusCode
	}
	return 0
}

// IsNotFound identifies a missing remote resource without matching text in unrelated errors.
func IsNotFound(err error) bool {
	code := APIStatusCode(err)
	return code == http.StatusNotFound || code == http.StatusGone
}
