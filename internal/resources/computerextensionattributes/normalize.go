package computerextensionattributes

import "strings"

func normalizeScript(script string) string {
	return strings.TrimRight(script, "\n")
}
