package cloud_distribution_point

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func boolPointerFromValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueBool()
	return &v
}

func intPointerFromValue(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := int(value.ValueInt64())
	return &v
}
