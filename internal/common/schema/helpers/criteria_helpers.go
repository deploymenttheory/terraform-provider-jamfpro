package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AttrTypes = map[string]attr.Type{
	"name":          types.StringType,
	"priority":      types.Int32Type,
	"and_or":        types.StringType,
	"search_type":   types.StringType,
	"value":         types.StringType,
	"opening_paren": types.BoolType,
	"closing_paren": types.BoolType,
}

var ObjectType = types.ObjectType{AttrTypes: AttrTypes}

// Expand converts a Terraform list of criteria into a slice of the specified type.
func Expand[T any](ctx context.Context, list types.List) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var result []T
	diags.Append(list.ElementsAs(ctx, &result, false)...) // populate when known

	return result, diags
}

// Flatten converts a slice of criteria of the specified type into a Terraform list.
func Flatten[T any](ctx context.Context, criteria []T) (types.List, diag.Diagnostics) {
	if criteria == nil {
		criteria = []T{}
	}

	return types.ListValueFrom(ctx, ObjectType, criteria)
}
