package logging

import (
	"context"
	"fmt"
)

const (
	// TF Resource Construction
	MsgTFConstructResourceFailure = "Failed to construct %s resource"
	MsgTFConstructResourceSuccess = "%s resource constructed successfully"

	// TF State
	MsgTFStateSyncFailure    = "Failed to synchronize Terraform state for %s"
	MsgTFStateSyncSuccess    = "Terraform state synchronized successfully for %s"
	MsgTFStateRemovalWarning = "Removing %s with ID: %s from Terraform state"

	// Jamf Pro
	MsgTFResourceDuplicateName = "A %s with the name '%s' already exists. %s names must be unique."

	// Create
	MsgAPICreateFailure          = "API error occurred during %s creation"
	MsgAPICreateFailedAfterRetry = "Final attempt to create %s failed"
	MsgAPICreateSuccess          = "%s created successfully though the API"

	// Read
	MsgAPIReadFailureByID      = "API error occurred while reading %s with ID: %s"
	MsgAPIReadFailureByName    = "API error occurred while reading %s with Name: %s"
	MsgAPIReadFailedAfterRetry = "Final attempt to read  %s failed"
	MsgAPIReadSuccess          = "%s with ID: %s successfully read from API"

	// Update
	MsgAPIUpdateFailureByID      = "API error occurred while updating %s with ID: %s"
	MsgAPIUpdateFailureByName    = "API error occurred while updating %s with Name: %s"
	MsgAPIUpdateFailedAfterRetry = "Final attempt to update %s failed"
	MsgAPIUpdateSuccess          = "%s successfully updated in Terraform state"

	// Delete
	MsgAPIDeleteFailureByID      = "API error occurred while deleting %s by ID"
	MsgAPIDeleteFailureByName    = "API error occurred while deleting %s by name"
	MsgAPIDeleteFailedAfterRetry = "Final attempt to delete %s failed"
	MsgAPIDeleteSuccess          = "%s successfully removed from Terraform state"
)

// TF Resource Construction

// LogTFConstructResourceFailure provides structured logging for errors during object construction
func LogTFConstructResourceFailure(ctx context.Context, resourceType, errorMsg string) {
	logMessage := fmt.Sprintf(MsgTFConstructResourceFailure, resourceType)

	Error(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
	})
}

// LogTFConstructResourceSuccess provides structured logging for successful object construction
func LogTFConstructResourceSuccess(ctx context.Context, resourceType string) {
	logMessage := fmt.Sprintf(MsgTFConstructResourceSuccess, resourceType)

	Info(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
	})
}

// LogTFConstructResource provides structured logging for successful construction and serialization of the resource object to XML
func LogTFConstructedXMLResource(ctx context.Context, resourceType, xmlData string) {
	logMessage := fmt.Sprintf("%s resource constructed and serialized to XML successfully", resourceType)

	Debug(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"xml":           xmlData,
	})
}

// LogTFConstructResourceXMLMarshalFailure provides structured logging for errors during marshaling of the resource object to XML
func LogTFConstructResourceXMLMarshalFailure(ctx context.Context, resourceType, errorMsg string) {
	logMessage := fmt.Sprintf("Failed to marshal %s object to XML", resourceType)

	Error(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
	})
}

// LogTFConstructedResource provides structured logging for successful object construction
func LogTFConstructedResource(ctx context.Context, resourceType string) {
	logMessage := fmt.Sprintf(MsgTFConstructResourceSuccess, resourceType)

	Info(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
	})
}

// LogTFConstructedJSONResource provides structured logging for successful construction and serialization of the resource object to JSON
func LogTFConstructedJSONResource(ctx context.Context, resourceType, jsonData string) {
	logMessage := fmt.Sprintf("%s resource constructed and serialized to JSON successfully", resourceType)

	Debug(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"json":          jsonData,
	})
}

// LogTFConstructResourceJSONMarshalFailure provides structured logging for errors during marshaling of the resource object to JSON
func LogTFConstructResourceJSONMarshalFailure(ctx context.Context, resourceType, errorMsg string) {
	logMessage := fmt.Sprintf("Failed to marshal %s object to JSON", resourceType)

	Error(ctx, SubsystemConstruct, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
	})
}

// TF State

// LogWarnRemoveFromState logs a warning when a resource is being removed from the Terraform state
func LogTFStateRemovalWarning(ctx context.Context, resourceType, resourceID string) {
	logMessage := fmt.Sprintf(MsgTFStateRemovalWarning, resourceType, resourceID)

	Warn(ctx, SubsystemRead, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
	})
}

// LogTFStateSyncFailure provides structured logging for errors during Terraform state synchronization
func LogTFStateSyncFailure(ctx context.Context, resourceType, errorMsg string) {
	logMessage := fmt.Sprintf(MsgTFStateSyncFailure, resourceType)

	Error(ctx, SubsystemSync, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
	})
}

// LogTFStateSyncSuccess provides structured logging for successful Terraform state synchronization
func LogTFStateSyncSuccess(ctx context.Context, resourceType, resourceID string) {
	logMessage := fmt.Sprintf(MsgTFStateSyncSuccess, resourceType)

	Info(ctx, SubsystemSync, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
	})
}

// LogTFStateSyncFailedAfterRetry provides structured logging for scenarios where state synchronization fails after retries
func LogTFStateSyncFailedAfterRetry(ctx context.Context, resourceType, resourceID, errorMsg string) {
	logMessage := fmt.Sprintf("Final attempt to synchronize Terraform state for %s with ID: %s failed", resourceType, resourceID)

	Error(ctx, SubsystemSync, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"error":         errorMsg,
	})
}

// Jamf Pro

// LogTFResourceDuplicateName provides structured logging for duplicate resource names
func LogTFResourceDuplicateName(ctx context.Context, resourceType, resourceName string) {
	logMessage := fmt.Sprintf(MsgTFResourceDuplicateName, resourceType, resourceName, resourceType)

	Error(ctx, SubsystemCreate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"name":          resourceName,
	})
}

// Create

// LogAPICreateFailure provides structured logging for API errors during creation
func LogAPICreateFailure(ctx context.Context, resourceType, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf(MsgAPICreateFailure, resourceType)

	Error(ctx, SubsystemCreate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPICreateSuccess provides structured logging for successful resource creation
func LogAPICreateSuccess(ctx context.Context, resourceType, resourceID string) {
	logMessage := fmt.Sprintf(MsgAPICreateSuccess, resourceType)

	Info(ctx, SubsystemCreate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
	})
}

// LogAPICreateFailedAfterRetryprovides structured logging for scenarios where update operation fails after retries
func LogAPICreateFailedAfterRetry(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string) {
	logMessage := fmt.Sprintf("retry attempt to create %s with ID: %s, Name: %q failed", resourceType, resourceID, resourceName)

	Error(ctx, SubsystemCreate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"name":          resourceName,
		"error":         errorMsg,
	})
}

// Read

// LogFailedReadByID provides structured logging for errors when failing to read by ID
func LogFailedReadByID(ctx context.Context, resourceType, resourceID, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf(MsgAPIReadFailureByID, resourceType, resourceID)

	Error(ctx, SubsystemRead, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogFailedReadByName provides structured logging for errors when failing to read by Name
func LogFailedReadByName(ctx context.Context, resourceType, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf(MsgAPIReadFailureByName, resourceType, resourceName)

	Error(ctx, SubsystemRead, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPIReadSuccess provides structured logging for successful read operations
func LogAPIReadSuccess(ctx context.Context, resourceType, resourceID string) {
	logMessage := fmt.Sprintf(MsgAPIReadSuccess, resourceType, resourceID)

	Info(ctx, SubsystemRead, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
	})
}

// LogAPIReadFailedAfterRetry provides structured logging for scenarios where update operation fails after retries
func LogAPIReadFailedAfterRetry(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf("Retry attempt to read %s with ID: %s, Name: %q failed", resourceType, resourceID, resourceName)

	Error(ctx, SubsystemRead, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// Update

// LogAPIUpdateFailureByID provides structured logging for errors when failing to update by ID
func LogAPIUpdateFailureByID(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf(MsgAPIUpdateFailureByID, resourceType, resourceID)

	Error(ctx, SubsystemUpdate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
		"error_code":    errorCode,
		"id":            resourceID,
		"name":          resourceName,
	})
}

// LogAPIUpdateFailureByName provides structured logging for errors when failing to update by Name
func LogAPIUpdateFailureByName(ctx context.Context, resourceType, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf(MsgAPIUpdateFailureByName, resourceType, resourceName)

	Error(ctx, SubsystemUpdate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"error":         errorMsg,
		"error_code":    errorCode,
		"name":          resourceName,
	})
}

// LogAPIUpdateFailedAfterRetry provides structured logging for scenarios where update operation fails after retries
func LogAPIUpdateFailedAfterRetry(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf("Final attempt to update %s with ID: %s, Name: %q failed", resourceType, resourceID, resourceName)

	Error(ctx, SubsystemUpdate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPIUpdateSuccess provides structured logging for successful update operations
func LogAPIUpdateSuccess(ctx context.Context, resourceType, resourceID, resourceName string) {
	logMessage := fmt.Sprintf("%s with (ID: %s, Name: %q) successfully updated in Terraform state", resourceType, resourceID, resourceName)

	Info(ctx, SubsystemUpdate, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
	})
}

// Delete

// LogAPIDeleteFailureByID provides structured logging for errors when failing to delete by ID
func LogAPIDeleteFailureByID(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf("Failed to delete %s with ID: %s, Name: %q", resourceType, resourceID, resourceName)

	Error(ctx, SubsystemDelete, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPIDeleteFailureByName provides structured logging for errors when failing to delete by Name
func LogAPIDeleteFailureByName(ctx context.Context, resourceType, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf("API error occurred while deleting %s with Name: %q", resourceType, resourceName)

	Error(ctx, SubsystemDelete, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPIDeleteFailedAfterRetry provides structured logging for scenarios where delete operation fails after retries
func LogAPIDeleteFailedAfterRetry(ctx context.Context, resourceType, resourceID, resourceName, errorMsg string, errorCode int) {
	logMessage := fmt.Sprintf("Final attempt to delete %s with ID: %s, Name: %q failed", resourceType, resourceID, resourceName)

	Error(ctx, SubsystemDelete, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
		"error":         errorMsg,
		"error_code":    errorCode,
	})
}

// LogAPIDeleteSuccess provides structured logging for successful delete operations
func LogAPIDeleteSuccess(ctx context.Context, resourceType, resourceID, resourceName string) {

	logMessage := fmt.Sprintf("%s with (ID: %s, Name: %q) successfully removed from Terraform state", resourceType, resourceID, resourceName)

	Info(ctx, SubsystemDelete, logMessage, map[string]interface{}{
		"resource_type": resourceType,
		"id":            resourceID,
		"name":          resourceName,
	})
}
