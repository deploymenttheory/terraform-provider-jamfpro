package smart_computer_group_v3

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
)

var errUnstructuredTestError = errors.New("400 INVALID_FIELD The criterion Test is not valid")

func TestModernWriteRefused(t *testing.T) {
	for _, tc := range []struct {
		name, code, description string
		status                  int
		want                    bool
	}{
		{"patch reporting defect", "INVALID_FIELD", "The criterion Patch Reporting: Test is not valid", 400, true},
		{"extension attribute defect", "INVALID_FIELD", "The operator has is not valid for extension attribute Notes", 400, true},
		{"lowercase error code", "invalid_field", "The criterion Test is not valid", 400, true},
		{"duplicate name", "DUPLICATE_FIELD", "The criterion Test is not valid", 422, false},
		{"forbidden", "INVALID_FIELD", "The criterion Test is not valid", 403, false},
		{"wrong priority", "INVALID_REQUEST_PARAMETER_VALUE", "Invalid priority order", 400, false},
		{"bad connector", "INVALID_FIELD", "Invalid andOr value: xor. Must be AND or OR", 400, false},
		{"invalid site", "INVALID_FIELD", "The site is not valid", 400, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			apiErr := &apiresponse.APIError{StatusCode: tc.status, Errors: []apiresponse.Errors{{
				Code: tc.code, Description: tc.description,
			}}}
			err := fmt.Errorf("SDK request failed: %w", apiErr)
			if got := modernWriteRefused(err); got != tc.want {
				t.Fatalf("fallback = %v, want %v", got, tc.want)
			}
		})
	}
	if modernWriteRefused(nil) || modernWriteRefused(errUnstructuredTestError) {
		t.Fatal("fallback requires a structured API error")
	}
}

func TestIsNotFoundUsesStatus(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusGone, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			err := fmt.Errorf("group ID 404: %w", &apiresponse.APIError{StatusCode: status})
			if isNotFound(err) != (status == 404 || status == 410) {
				t.Fatal("resource IDs must not be mistaken for HTTP status codes")
			}
		})
	}
}
