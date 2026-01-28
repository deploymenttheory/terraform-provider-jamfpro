package computer_extension_attribute

import "strings"

// normalizeScript normalizes a script by replacing all CRLF with LF and trimming trailing newlines
func normalizeScript(script string) string {
	normalized := strings.ReplaceAll(script, "\r\n", "\n")

	return strings.TrimRight(normalized, "\n")
}
