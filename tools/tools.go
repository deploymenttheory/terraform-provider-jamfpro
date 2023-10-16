//go:build tools
// +build tools

package tools

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlintx"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/katbyte/terrafmt"
	_ "github.com/pavius/impi/cmd/impi"
	_ "github.com/terraform-linters/tflint"
)
