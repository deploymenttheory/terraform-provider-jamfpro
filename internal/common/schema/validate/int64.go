package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// =============================================================================
// Int64: Requirement Validators
// =============================================================================

var _ validator.Int64 = int64RequiredWhenOneOfValidator{}

// int64RequiredWhenOneOfValidator validates that an int64 field is required when another field matches any of the specified values.
type int64RequiredWhenOneOfValidator struct {
	dependentField string
	requiredValues []string
}

// Description returns a plain-text description of the validator's behavior.
func (v int64RequiredWhenOneOfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("this attribute is required when %s is one of: %v", v.dependentField, v.requiredValues)
}

// MarkdownDescription returns a markdown-formatted description of the validator's behavior.
func (v int64RequiredWhenOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateInt64 performs the validation.
func (v int64RequiredWhenOneOfValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	// If the attribute being validated is not configured, we don't need to check the dependency.
	if req.ConfigValue.IsUnknown() {
		return
	}

	dependentPath := req.Path.ParentPath().AtName(v.dependentField)

	var dependentValue types.String
	diags := req.Config.GetAttribute(ctx, dependentPath, &dependentValue)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If the dependent attribute isn't set to one of the required values, the validation passes.
	if dependentValue.IsUnknown() || dependentValue.IsNull() {
		return
	}

	valueMatches := false
	matchedValue := ""
	for _, reqVal := range v.requiredValues {
		if dependentValue.ValueString() == reqVal {
			valueMatches = true
			matchedValue = reqVal
			break
		}
	}

	if !valueMatches {
		return
	}

	// At this point, the dependent field has one of the required values, so this field must be set.
	// Check if the current attribute value is null.
	if req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute Required",
			fmt.Sprintf("Attribute %q is required because attribute %q is set to %q.", req.Path, dependentPath, matchedValue),
		)
	}
}

// Int64RequiredWhenOneOf returns a validator that ensures the int64 attribute is not null when another field
// matches any of the specified values.
func Int64RequiredWhenOneOf(dependentField string, requiredValues ...string) validator.Int64 {
	return int64RequiredWhenOneOfValidator{
		dependentField: dependentField,
		requiredValues: requiredValues,
	}
}

//--------------------------------------------------------

// Insert next validator here
