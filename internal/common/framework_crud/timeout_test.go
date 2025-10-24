package framework_crud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleTimeout(t *testing.T) {
	t.Run("Successful timeout setting", func(t *testing.T) {
		ctx := context.Background()
		defaultTimeout := 5 * time.Second
		var diags diag.Diagnostics

		timeoutFunc := func(ctx context.Context, d time.Duration) (time.Duration, diag.Diagnostics) {
			return 10 * time.Second, diag.Diagnostics{}
		}

		newCtx, cancel := HandleTimeout(ctx, timeoutFunc, defaultTimeout, &diags)
		defer cancel()

		require.NotNil(t, newCtx)
		require.NotNil(t, cancel)
		assert.False(t, diags.HasError())

		deadline, ok := newCtx.Deadline()
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(10*time.Second), deadline, 100*time.Millisecond)
	})

	t.Run("Error in timeout function", func(t *testing.T) {
		ctx := context.Background()
		defaultTimeout := 5 * time.Second
		var diags diag.Diagnostics

		timeoutFunc := func(ctx context.Context, d time.Duration) (time.Duration, diag.Diagnostics) {
			return 0, diag.Diagnostics{diag.NewErrorDiagnostic("Error", "Timeout function failed")}
		}

		newCtx, cancel := HandleTimeout(ctx, timeoutFunc, defaultTimeout, &diags)

		assert.Equal(t, ctx, newCtx)
		assert.Nil(t, cancel)
		assert.True(t, diags.HasError())
		assert.Len(t, diags, 1)
		assert.Equal(t, "Error", diags[0].Summary())
		assert.Equal(t, "Timeout function failed", diags[0].Detail())
	})

	t.Run("Default timeout used", func(t *testing.T) {
		ctx := context.Background()
		defaultTimeout := 5 * time.Second
		var diags diag.Diagnostics

		timeoutFunc := func(ctx context.Context, d time.Duration) (time.Duration, diag.Diagnostics) {
			return d, diag.Diagnostics{}
		}

		newCtx, cancel := HandleTimeout(ctx, timeoutFunc, defaultTimeout, &diags)
		defer cancel()

		require.NotNil(t, newCtx)
		require.NotNil(t, cancel)
		assert.False(t, diags.HasError())

		deadline, ok := newCtx.Deadline()
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(defaultTimeout), deadline, 100*time.Millisecond)
	})

	t.Run("Zero timeout", func(t *testing.T) {
		ctx := context.Background()
		defaultTimeout := 5 * time.Second
		var diags diag.Diagnostics

		timeoutFunc := func(ctx context.Context, d time.Duration) (time.Duration, diag.Diagnostics) {
			return 0, diag.Diagnostics{}
		}

		newCtx, cancel := HandleTimeout(ctx, timeoutFunc, defaultTimeout, &diags)
		defer cancel()

		require.NotNil(t, newCtx)
		require.NotNil(t, cancel)
		assert.False(t, diags.HasError())

		deadline, ok := newCtx.Deadline()
		assert.True(t, ok)
		assert.WithinDuration(t, time.Now(), deadline, 100*time.Millisecond)
	})
}
