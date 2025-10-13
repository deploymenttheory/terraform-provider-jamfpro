package crud

import (
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
	if !diagnostics.HasError() {
		return ErrorInfo{}
	}

	// Iterate through diagnostics to find HTTP error information
	for _, d := range diagnostics.Errors() {
		summary := d.Summary()
		detail := d.Detail()

		// Combine summary and detail for analysis
		errorText := summary + " " + detail
		errorTextLower := strings.ToLower(errorText)

		// Try to extract status code patterns from error messages
		if strings.Contains(errorTextLower, "404") || strings.Contains(errorTextLower, "not found") {
			return ErrorInfo{StatusCode: 404, ErrorCode: "NotFound"}
		}

		if strings.Contains(errorTextLower, "400") || strings.Contains(errorTextLower, "bad request") {
			return ErrorInfo{StatusCode: 400, ErrorCode: "BadRequest"}
		}

		if strings.Contains(errorTextLower, "401") || strings.Contains(errorTextLower, "unauthorized") {
			return ErrorInfo{StatusCode: 401, ErrorCode: "Unauthorized"}
		}

		if strings.Contains(errorTextLower, "403") || strings.Contains(errorTextLower, "forbidden") {
			return ErrorInfo{StatusCode: 403, ErrorCode: "Forbidden"}
		}

		if strings.Contains(errorTextLower, "409") || strings.Contains(errorTextLower, "conflict") {
			return ErrorInfo{StatusCode: 409, ErrorCode: "Conflict"}
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

// isRetryableReadError determines if an error should trigger a retry for read operations
func isRetryableReadError(errorInfo *ErrorInfo) bool {
	if errorInfo == nil {
		return true // Unknown errors are retryable for safety
	}

	switch errorInfo.StatusCode {
	case 404, 409, 423, 429: // Not found (propagation), conflict, locked, rate limited
		return true
	case 500, 502, 503, 504: // Server errors
		return true
	default:
		// Check specific error codes that might be retryable
		retryableErrorCodes := map[string]bool{
			"ServiceUnavailable":  true,
			"RequestThrottled":    true,
			"RequestTimeout":      true,
			"InternalServerError": true,
			"BadGateway":          true,
			"GatewayTimeout":      true,
			"NotFound":            true, // Resource propagation
			"ResourceNotFound":    true, // Resource propagation
			"NetworkError":        true, // Network connectivity issues
		}
		return retryableErrorCodes[errorInfo.ErrorCode]
	}
}

// isNonRetryableReadError determines if an error should NOT trigger a retry for read operations
func isNonRetryableReadError(errorInfo *ErrorInfo) bool {
	if errorInfo == nil {
		return false
	}

	switch errorInfo.StatusCode {
	case 200, 204: // Success cases
		return true
	case 400, 401, 403, 405, 406, 410, 422: // Client errors that won't change on retry
		return true
	// Note: 404 is NOT here because it's retryable for reads (propagation)
	// Note: 409 Conflict is retryable for reads (removed from here)
	default:
		// Check specific error codes that are permanent failures
		nonRetryableErrorCodes := map[string]bool{
			"BadRequest":          true,
			"Unauthorized":        true,
			"Forbidden":           true,
			"Gone":                true,
			"UnprocessableEntity": true,
			"ValidationError":     true,
			// Note: "NotFound" is NOT here because it's retryable for reads (propagation)
		}
		return nonRetryableErrorCodes[errorInfo.ErrorCode]
	}
}
