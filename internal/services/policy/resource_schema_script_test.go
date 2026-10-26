package policy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Script parameters can carry credentials, so Terraform must mask them in plans.
func TestPolicyScriptParametersAreSensitive(t *testing.T) {
	scripts := ResourceJamfProPolicies().CoreConfigSchema().BlockTypes["payloads"].BlockTypes["scripts"]
	for i := 4; i <= 11; i++ {
		name := fmt.Sprintf("parameter%d", i)
		require.Contains(t, scripts.Attributes, name)
		require.True(t, scripts.Attributes[name].Sensitive, "%s must be sensitive", name)
	}
	require.False(t, scripts.Attributes["id"].Sensitive)
	require.False(t, scripts.Attributes["priority"].Sensitive)
}
