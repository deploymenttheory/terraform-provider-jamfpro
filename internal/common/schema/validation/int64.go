package validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Int64SequenceModel exposes a sequence of int64 values that can be validated.
type Int64SequenceModel interface {
	GetInt64Sequence() []int64
}

// IncrementingInt64SequenceValidator ensures a sequence starts at 0 and increments by 1.
type IncrementingInt64SequenceValidator[T Int64SequenceModel] struct{}

// Description explains what the validator enforces.
func (IncrementingInt64SequenceValidator[T]) Description(ctx context.Context) string {
	return "Ensures the sequence starts at 0 and increments by 1 for each subsequent value"
}

// MarkdownDescription explains what the validator enforces (markdown format).
func (IncrementingInt64SequenceValidator[T]) MarkdownDescription(ctx context.Context) string {
	return "Ensures the sequence starts at 0 and increments by 1 for each subsequent value"
}

// ValidateResource retrieves the model and validates its int64 sequence.
func (IncrementingInt64SequenceValidator[T]) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data T

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateIncrementingInt64Sequence(data.GetInt64Sequence()); err != nil {
		resp.Diagnostics.AddError("Invalid Int64 Sequence", err.Error())
	}
}

// ValidateIncrementingInt64Sequence ensures the slice starts at 0 and increments by 1.
func ValidateIncrementingInt64Sequence(values []int64) error {
	if len(values) <= 1 {
		return nil
	}

	for idx, value := range values {
		expected := int64(idx)
		if value != expected {
			if idx == 0 {
				return fmt.Errorf("the first value must be 0, got %d", value)
			}

			return fmt.Errorf("value at index %d has an invalid value %d, expected %d", idx, value, expected)
		}
	}

	return nil
}
