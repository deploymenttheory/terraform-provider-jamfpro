package framework

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// HandleTimeout is a helper function to manage context timeouts
func HandleTimeout(ctx context.Context, timeoutFunc func(context.Context, time.Duration) (time.Duration, diag.Diagnostics), defaultTimeout time.Duration, diags *diag.Diagnostics) (context.Context, context.CancelFunc) {
	timeout, timeoutDiags := timeoutFunc(ctx, defaultTimeout)
	*diags = append(*diags, timeoutDiags...)
	if diags.HasError() {
		return ctx, nil
	}
	return context.WithTimeout(ctx, timeout)
}
