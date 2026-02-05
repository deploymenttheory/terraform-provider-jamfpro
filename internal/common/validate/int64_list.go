package validate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.List = int64ListSumEqualsValidator{}

// int64ListSumEqualsValidator validates that the sum of int64 values in a list equals a specific value.
type int64ListSumEqualsValidator struct {
	expectedSum int64
}

// Description provides a human-readable description of the validator's purpose.
func (v int64ListSumEqualsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v int64ListSumEqualsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("sum of values must equal %d", v.expectedSum)
}

// ValidateList performs the validation logic on the list attribute.
func (v int64ListSumEqualsValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {

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

// Int64ListSumEquals returns a validator that ensures the sum of int64 values in a list equals the expected value.
func Int64ListSumEquals(expectedSum int64) validator.List {
	return int64ListSumEqualsValidator{
		expectedSum: expectedSum,
	}
}
