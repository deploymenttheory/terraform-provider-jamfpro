package utilities

import "net/http"

// ResponseWasNotFound checks if an HTTP status code indicates that a resource was not found (HTTP status code 404).
func ResponseWasNotFound(statusCode int) bool {
	return statusCode == http.StatusNotFound
}

// ResponseWasBadRequest checks if an HTTP status code indicates a bad request (HTTP status code 400).
func ResponseWasBadRequest(statusCode int) bool {
	return statusCode == http.StatusBadRequest
}

// ResponseWasForbidden checks if an HTTP status code indicates a forbidden request (HTTP status code 403).
func ResponseWasForbidden(statusCode int) bool {
	return statusCode == http.StatusForbidden
}

// ResponseWasConflict checks if an HTTP status code indicates a conflict (HTTP status code 409).
func ResponseWasConflict(statusCode int) bool {
	return statusCode == http.StatusConflict
}
