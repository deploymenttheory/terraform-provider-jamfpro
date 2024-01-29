package logging

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// DiagnosticLogger is an interface for logging messages that can create diagnostics for errors and warnings.
type DiagnosticLogger interface {
	// Error logs an error message and creates a diagnostic with Error severity.
	Error(summary string, detail string, attributePath ...interface{})

	// Errorf logs a formatted error message and creates a diagnostic with Error severity.
	Errorf(format string, args ...interface{})

	// Warn logs a warning message and creates a diagnostic with Warning severity.
	Warn(summary string, detail string, attributePath ...interface{})

	// Warnf logs a formatted warning message and creates a diagnostic with Warning severity.
	Warnf(format string, args ...interface{})
}

// Ensure ConsoleLogger implements DiagnosticLogger interface.
var _ DiagnosticLogger = &ConsoleLogger{}

// ConsoleLogger provides an implementation of the DiagnosticLogger interface.
// It logs messages to StdOut and creates diagnostics for errors and warnings.
type ConsoleLogger struct {
	Diagnostics diag.Diagnostics
}

// Error logs an error message and creates a diagnostic with Error severity.
func (l *ConsoleLogger) Error(summary string, detail string, attributePath ...interface{}) {
	log.Printf("[ERROR] %s: %s", summary, detail)
	l.Diagnostics = append(l.Diagnostics, diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       summary,
		Detail:        detail,
		AttributePath: l.constructAttributePath(attributePath...),
	})
}

// Errorf logs a formatted error message and creates a diagnostic with Error severity.
func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Error(message, "", nil)
}

// Warn logs a warning message and creates a diagnostic with Warning severity.
func (l *ConsoleLogger) Warn(summary string, detail string, attributePath ...interface{}) {
	log.Printf("[WARN] %s", summary)
	l.Diagnostics = append(l.Diagnostics, diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       summary,
		Detail:        summary, // Detail is the same as summary for warnings.
		AttributePath: l.constructAttributePath(attributePath...),
	})
}

// Warnf logs a formatted warning message and creates a diagnostic with Warning severity.
func (l *ConsoleLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Warn(message, message, nil) // Passing message as both summary and detail.
}

// Info logs an informational message without creating a diagnostic.
func (l *ConsoleLogger) Info(message string) {
	log.Printf("[INFO] %s", message)
}

// Infof logs a formatted informational message without creating a diagnostic.
func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// constructAttributePath converts a slice of interfaces to a cty.Path for setting the AttributePath in a Diagnostic.
func (l *ConsoleLogger) constructAttributePath(attributePath ...interface{}) cty.Path {
	var path cty.Path
	for _, step := range attributePath {
		switch s := step.(type) {
		case string:
			path = append(path, cty.GetAttrStep{Name: s})
		case int:
			path = append(path, cty.IndexStep{Key: cty.NumberIntVal(int64(s))})
			// Additional cases for other types as necessary.
		}
	}
	return path
}

// NullDiagnosticLogger is an implementation of the DiagnosticLogger interface that disregards log output.
type NullDiagnosticLogger struct{}

// Ensure NullDiagnosticLogger implements DiagnosticLogger interface.
var _ DiagnosticLogger = &NullDiagnosticLogger{}

// Error disregards the error log output.
func (NullDiagnosticLogger) Error(_ string, _ string, _ ...interface{}) {
	// No operation
}

// Errorf disregards the formatted error log output.
func (NullDiagnosticLogger) Errorf(_ string, _ ...interface{}) {
	// No operation
}

// Warn disregards the warning log output.
func (NullDiagnosticLogger) Warn(_ string, _ string, _ ...interface{}) {
	// No operation
}

// Warnf disregards the formatted warning log output.
func (NullDiagnosticLogger) Warnf(_ string, _ ...interface{}) {
	// No operation
}

// Info disregards the informational log output.
func (NullDiagnosticLogger) Info(_ string) {
	// No operation
}

// Infof disregards the formatted informational log output.
func (NullDiagnosticLogger) Infof(_ string, _ ...interface{}) {
	// No operation
}
