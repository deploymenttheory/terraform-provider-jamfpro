package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = requiredWhenEqualsValidator{}

// requiredWhenEqualsValidator is the implementation of the validator.
type requiredWhenEqualsValidator struct {
	dependentField string
	requiredValue  types.String
}

// Description returns a plain-text description of the validator's behavior.
func (v requiredWhenEqualsValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("this attribute is required when %s is set to %s", v.dependentField, v.requiredValue.ValueString())
}

// MarkdownDescription returns a markdown-formatted description of the validator's behavior.
func (v requiredWhenEqualsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWhenEqualsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the attribute being validated is not configured, we don't need to check the dependency.
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Get the path to the dependent attribute.
	dependentPath := req.Path.ParentPath().AtName(v.dependentField)

	// Get the value of the dependent attribute from the configuration.
	var dependentValue types.String
	diags := req.Config.GetAttribute(ctx, dependentPath, &dependentValue)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If the dependent attribute isn't set to the required value, the validation passes.
	if dependentValue.IsUnknown() || dependentValue.IsNull() || !dependentValue.Equal(v.requiredValue) {
		return
	}

	// At this point, the dependent field has the required value, so this field must be set.
	// Check if the current attribute value is null or an empty string.
	if req.ConfigValue.IsNull() || req.ConfigValue.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute Required",
			fmt.Sprintf("Attribute %q is required because attribute %q is set to %q.", req.Path, dependentPath, v.requiredValue.ValueString()),
		)
	}
}

// RequiredWhenEquals is a factory function that returns a new requiredWhenEqualsValidator.
// It validates that the attribute is not null or empty when another field in the same object
// has a specific string value.
func RequiredWhenEquals(dependentField string, requiredValue types.String) validator.String {
	return requiredWhenEqualsValidator{
		dependentField: dependentField,
		requiredValue:  requiredValue,
	}
}

var _ validator.String = requiredWhenOneOfValidator{}

// requiredWhenOneOfValidator validates that a string field is required when another field matches any of the specified values.
type requiredWhenOneOfValidator struct {
	dependentField string
	requiredValues []string
}

// Description returns a plain-text description of the validator's behavior.
func (v requiredWhenOneOfValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("this attribute is required when %s is one of: %v", v.dependentField, v.requiredValues)
}

// MarkdownDescription returns a markdown-formatted description of the validator's behavior.
func (v requiredWhenOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWhenOneOfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the attribute being validated is not configured, we don't need to check the dependency.
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Get the path to the dependent attribute.
	dependentPath := req.Path.ParentPath().AtName(v.dependentField)

	// Get the value of the dependent attribute from the configuration.
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

	// Check if the dependent value matches any of the required values
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
	// Check if the current attribute value is null or an empty string.
	if req.ConfigValue.IsNull() || req.ConfigValue.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute Required",
			fmt.Sprintf("Attribute %q is required because attribute %q is set to %q.", req.Path, dependentPath, matchedValue),
		)
	}
}

// RequiredWhenOneOf returns a validator that ensures the attribute is not null or empty when another field
// matches any of the specified values.
func RequiredWhenOneOf(dependentField string, requiredValues ...string) validator.String {
	return requiredWhenOneOfValidator{
		dependentField: dependentField,
		requiredValues: requiredValues,
	}
}
