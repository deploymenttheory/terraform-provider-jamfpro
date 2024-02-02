package logging

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// LogSubsystem identifies different subsystems in your Terraform provider.
type LogSubsystem string

const (
	SubsystemDelete     LogSubsystem = "delete"
	SubsystemCreate     LogSubsystem = "create"
	SubsystemRead       LogSubsystem = "read"
	SubsystemUpdate     LogSubsystem = "update"
	SubsystemSync       LogSubsystem = "sync"       // For synchronization operations
	SubsystemAPI        LogSubsystem = "api"        // For direct API interaction logs
	SubsystemRetry      LogSubsystem = "retry"      // For retry logic
	SubsystemValidation LogSubsystem = "validation" // For input validation
	SubsystemConfig     LogSubsystem = "config"     // For configuration-related logs
	SubsystemInit       LogSubsystem = "init"       // For provider initialization
	SubsystemCleanup    LogSubsystem = "cleanup"    // For cleanup operations
	SubsystemConstruct  LogSubsystem = "construct"  // For resource construction operations
	SubsystemGeneral    LogSubsystem = "general"
	// Add more subsystems as needed
)

// NewSubsystemLogger creates a new logger for a specific subsystem with the provided log level.
func NewSubsystemLogger(ctx context.Context, subsystem LogSubsystem, level hclog.Level) context.Context {
	// Assuming `WithLevel` is a function that returns `logging.Option` and is accessible here
	levelOption := tflog.WithLevel(level) // Replace with the actual function that returns a `logging.Option` for setting the log level

	// Create a new subsystem logger with the specified level option
	subCtx := tflog.NewSubsystem(ctx, string(subsystem), levelOption)
	return subCtx
}

// mergeFields merges subsystem info with additional fields.
func MergeFields(subsystem LogSubsystem, additionalFields map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"subsystem": string(subsystem),
	}
	for k, v := range additionalFields {
		fields[k] = v
	}
	return fields
}

// Debug logs a debug message with additional fields for a specific subsystem.
func Debug(ctx context.Context, subsystem LogSubsystem, msg string, additionalFields map[string]interface{}) {
	tflog.SubsystemDebug(ctx, string(subsystem), msg, additionalFields)
}

// Info logs an info message with additional fields for a specific subsystem.
func Info(ctx context.Context, subsystem LogSubsystem, msg string, additionalFields map[string]interface{}) {
	tflog.SubsystemInfo(ctx, string(subsystem), msg, additionalFields)
}

// Warn logs a warning message with additional fields for a specific subsystem.
func Warn(ctx context.Context, subsystem LogSubsystem, msg string, additionalFields map[string]interface{}) {
	tflog.SubsystemWarn(ctx, string(subsystem), msg, additionalFields)
}

// Error logs an error message with additional fields for a specific subsystem.
func Error(ctx context.Context, subsystem LogSubsystem, msg string, additionalFields map[string]interface{}) {
	tflog.SubsystemError(ctx, string(subsystem), msg, additionalFields)
}

// MaskSensitiveData masks sensitive data in logs based on regex patterns or strings for a specific subsystem.
func MaskSensitiveData(ctx context.Context, subsystem LogSubsystem, expressions []*regexp.Regexp, strings []string) context.Context {
	ctx = tflog.SubsystemMaskAllFieldValuesRegexes(ctx, string(subsystem), expressions...)
	ctx = tflog.SubsystemMaskAllFieldValuesStrings(ctx, string(subsystem), strings...)
	ctx = tflog.SubsystemMaskMessageRegexes(ctx, string(subsystem), expressions...)
	ctx = tflog.SubsystemMaskMessageStrings(ctx, string(subsystem), strings...)
	return ctx
}

// OmitLogs omits logs containing specified keys or matching certain message patterns for a specific subsystem.
func OmitLogs(ctx context.Context, subsystem LogSubsystem, keys []string, expressions []*regexp.Regexp, strings []string) context.Context {
	ctx = tflog.SubsystemOmitLogWithFieldKeys(ctx, string(subsystem), keys...)
	ctx = tflog.SubsystemOmitLogWithMessageRegexes(ctx, string(subsystem), expressions...)
	ctx = tflog.SubsystemOmitLogWithMessageStrings(ctx, string(subsystem), strings...)
	return ctx
}

// SetLogField adds a field to all logs emitted from the provided context for a specific subsystem.
func SetLogField(ctx context.Context, subsystem LogSubsystem, key string, value interface{}) context.Context {
	return tflog.SubsystemSetField(ctx, string(subsystem), key, value)
}
