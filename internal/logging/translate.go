package logging

// TranslateLogLevel maps Terraform configuration string values to logger.LogLevel.
func TranslateLogLevel(logLevelStr string) string {
	switch logLevelStr {
	case "debug":
		return "LogLevelDebug"
	case "info":
		return "LogLevelInfo"
	case "warning":
		return "LogLevelWarn"
	case "none":
		return "LogLevelNone" // This handles the case where logging is to be disabled.
	default:
		return "LogLevelNone" // Defaults to warning if no match is found.
	}
}
