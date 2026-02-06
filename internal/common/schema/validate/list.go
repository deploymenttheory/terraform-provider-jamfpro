package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// =============================================================================
// List: Requirement Validators
// =============================================================================

var _ validator.List = listRequiredWhenEqualsValidator{}

// listRequiredWhenEqualsValidator validates that a list field is required when another field equals a specific value.
type listRequiredWhenEqualsValidator struct {
	dependentField string
	requiredValue  string
}

// Description returns a plain-text description of the validator's behavior.
func (v listRequiredWhenEqualsValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("this attribute is required when %s is set to %s", v.dependentField, v.requiredValue)
}

// MarkdownDescription returns a markdown-formatted description of the validator's behavior.
func (v listRequiredWhenEqualsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateList performs the validation.
func (v listRequiredWhenEqualsValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {

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

	// If the dependent attribute isn't set to the required value, the validation passes.
	if dependentValue.IsUnknown() || dependentValue.IsNull() || dependentValue.ValueString() != v.requiredValue {
		return
	}

	if req.ConfigValue.IsNull() || len(req.ConfigValue.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute Required",
			fmt.Sprintf("Attribute %q is required because attribute %q is set to %q.", req.Path, dependentPath, v.requiredValue),
		)
	}
}

// ListRequiredWhenEquals returns a validator that ensures the list attribute is not null or empty when another field
// has a specific string value.
func ListRequiredWhenEquals(dependentField string, requiredValue string) validator.List {
	return listRequiredWhenEqualsValidator{
		dependentField: dependentField,
		requiredValue:  requiredValue,
	}
}

// =============================================================================
// List: Value Validators
// =============================================================================

var _ validator.List = listInt64SumEqualsValidator{}

// listInt64SumEqualsValidator validates that the sum of int64 values in a list equals a specific value.
type listInt64SumEqualsValidator struct {
	expectedSum int64
}

// Description provides a human-readable description of the validator's purpose.
func (v listInt64SumEqualsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v listInt64SumEqualsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("sum of values must equal %d", v.expectedSum)
}

// ValidateList performs the validation logic on the list attribute.
func (v listInt64SumEqualsValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	var sum int64
	for _, elem := range request.ConfigValue.Elements() {
		int64Elem, ok := elem.(types.Int64)
		if !ok {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid List Element Type",
				"Expected all elements to be int64 values.",
			)
			return
		}

		if int64Elem.IsNull() || int64Elem.IsUnknown() {
			return
		}

		sum += int64Elem.ValueInt64()
	}

	if sum != v.expectedSum {
		message := fmt.Sprintf("Sum of values is %d but expected %d. All percentages must sum to exactly %d.", sum, v.expectedSum, v.expectedSum)
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid List Sum",
			message,
		)
	}
}

// ListInt64SumEquals returns a validator that ensures the sum of int64 values in a list equals the expected value.
func ListInt64SumEquals(expectedSum int64) validator.List {
	return listInt64SumEqualsValidator{
		expectedSum: expectedSum,
	}
}

//--------------------------------------------------------

// Insert next validator here
