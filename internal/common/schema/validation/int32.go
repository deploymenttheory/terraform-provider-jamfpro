package validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Int32SequenceModel exposes a sequence of int32 values that can be validated.
type Int32SequenceModel interface {
	GetInt32Sequence() []int32
}

// IncrementingInt32SequenceValidator ensures a sequence starts at 0 and increments by 1.
type IncrementingInt32SequenceValidator[T Int32SequenceModel] struct{}

// Description explains what the validator enforces.
func (IncrementingInt32SequenceValidator[T]) Description(ctx context.Context) string {
	return "Ensures the sequence starts at 0 and increments by 1 for each subsequent value"
}

// MarkdownDescription explains what the validator enforces (markdown format).
func (IncrementingInt32SequenceValidator[T]) MarkdownDescription(ctx context.Context) string {
	return "Ensures the sequence starts at 0 and increments by 1 for each subsequent value"
}

// ValidateResource retrieves the model and validates its int32 sequence.
func (IncrementingInt32SequenceValidator[T]) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data T

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateIncrementingInt32Sequence(data.GetInt32Sequence()); err != nil {
		resp.Diagnostics.AddError("Invalid Int32 Sequence", err.Error())
	}
}

// ValidateIncrementingInt32Sequence ensures the slice starts at 0 and increments by 1.
func ValidateIncrementingInt32Sequence(values []int32) error {
	if len(values) <= 1 {
		return nil
	}

	for idx, value := range values {
		expected := int32(idx)
		if value != expected {
			if idx == 0 {
				return fmt.Errorf("the first value must be 0, got %d", value)
			}

			return fmt.Errorf("value at index %d has an invalid value %d, expected %d", idx, value, expected)
		}
	}

	return nil
}
