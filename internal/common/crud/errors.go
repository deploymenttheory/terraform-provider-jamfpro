package crud

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ErrorInfo represents error information extracted from diagnostics
type ErrorInfo struct {
	StatusCode int
	ErrorCode  string
}

// extractErrorFromDiagnostics analyzes Terraform diagnostics to extract HTTP error information
// for intelligent retry decisions
func extractErrorFromDiagnostics(diagnostics diag.Diagnostics) ErrorInfo {
	// Iterate through diagnostics to find HTTP error information
	for _, d := range diagnostics.Errors() {
		errorTextLower := strings.ToLower(fmt.Sprintf("%s - %s", d.Summary(), d.Detail()))

		// Try to extract status code patterns from error messages
		if strings.Contains(errorTextLower, "404") || strings.Contains(errorTextLower, "not found") {
			return ErrorInfo{StatusCode: 404, ErrorCode: "not found"}
		}

		if strings.Contains(errorTextLower, "400") || strings.Contains(errorTextLower, "bad request") {
			return ErrorInfo{StatusCode: 400, ErrorCode: "bad request"}
		}

		if strings.Contains(errorTextLower, "401") || strings.Contains(errorTextLower, "unauthorized") {
			return ErrorInfo{StatusCode: 401, ErrorCode: "unauthorized"}
		}

		if strings.Contains(errorTextLower, "403") || strings.Contains(errorTextLower, "forbidden") {
			return ErrorInfo{StatusCode: 403, ErrorCode: "forbidden"}
		}

		if strings.Contains(errorTextLower, "409") || strings.Contains(errorTextLower, "conflict") {
			return ErrorInfo{StatusCode: 409, ErrorCode: "conflict"}
		}

		if strings.Contains(errorTextLower, "423") || strings.Contains(errorTextLower, "locked") {
			return ErrorInfo{StatusCode: 423, ErrorCode: "Locked"}
		}

		if strings.Contains(errorTextLower, "429") || strings.Contains(errorTextLower, "too many requests") || strings.Contains(errorTextLower, "throttl") {
			return ErrorInfo{StatusCode: 429, ErrorCode: "TooManyRequests"}
		}

		if strings.Contains(errorTextLower, "500") || strings.Contains(errorTextLower, "internal server error") {
			return ErrorInfo{StatusCode: 500, ErrorCode: "InternalServerError"}
		}

		if strings.Contains(errorTextLower, "502") || strings.Contains(errorTextLower, "bad gateway") {
			return ErrorInfo{StatusCode: 502, ErrorCode: "BadGateway"}
		}

		if strings.Contains(errorTextLower, "503") || strings.Contains(errorTextLower, "service unavailable") {
			return ErrorInfo{StatusCode: 503, ErrorCode: "ServiceUnavailable"}
		}

		if strings.Contains(errorTextLower, "504") || strings.Contains(errorTextLower, "gateway timeout") {
			return ErrorInfo{StatusCode: 504, ErrorCode: "GatewayTimeout"}
		}

		// Look for Jamf Pro specific error patterns
		if strings.Contains(errorTextLower, "service unavailable") {
			return ErrorInfo{StatusCode: 503, ErrorCode: "ServiceUnavailable"}
		}

		if strings.Contains(errorTextLower, "request throttled") {
			return ErrorInfo{StatusCode: 429, ErrorCode: "RequestThrottled"}
		}

		if strings.Contains(errorTextLower, "resource not found") {
			return ErrorInfo{StatusCode: 404, ErrorCode: "ResourceNotFound"}
		}

		if strings.Contains(errorTextLower, "network error") {
			return ErrorInfo{StatusCode: 500, ErrorCode: "NetworkError"}
		}

		if strings.Contains(errorTextLower, "timeout") {
			return ErrorInfo{StatusCode: 504, ErrorCode: "RequestTimeout"}
		}
	}

	// Return empty error info if no patterns match
	return ErrorInfo{}
}

// Retryable/non-retryable codes based on knowledge of the JP's API behaviors.
func isRetryableReadError(statusCode int) bool {
	if statusCode == 0 {
		return true
	}

	retryableCodes := []int{
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusLocked,
		http.StatusTooManyRequests,
		http.StatusRequestTimeout,
		http.StatusFailedDependency,
		http.StatusTooEarly,
		http.StatusGatewayTimeout,
	}

	return slices.Contains(retryableCodes, statusCode)
}
