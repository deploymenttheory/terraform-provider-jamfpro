package framework_crud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var errWriteReadBackNeedsDeadline = errors.New("write read-back requires a context deadline")

// WaitForWrite polls SDK reads until matches confirms the written definition is visible.
// Unlike ReadWithRetry, an error-free read is not sufficient: APIs can return stale
// data after a successful write. The caller supplies the operation deadline and
// owns state updates; the last observed value is returned even on failure.
func WaitForWrite[T any](
	ctx context.Context,
	read func(context.Context) (*T, error),
	matches func(*T) bool,
	interval time.Duration,
) (*T, error) {
	if _, ok := ctx.Deadline(); !ok {
		return nil, errWriteReadBackNeedsDeadline
	}
	if interval <= 0 {
		interval = 250 * time.Millisecond
	}
	const maxInterval = 5 * time.Second
	interval = min(interval, maxInterval)
	var last *T
	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return last, fmt.Errorf("resource write did not converge: %w", err)
		}
		got, err := read(ctx)
		if err != nil {
			code := APIStatusCode(err)
			if !IsNotFound(err) && code != http.StatusTooManyRequests && code < http.StatusInternalServerError {
				return last, fmt.Errorf("read resource after write: %w", err)
			}
		} else if got != nil {
			last = got
			if matches(got) {
				return got, nil
			}
		}
		tflog.Debug(ctx, "Waiting for resource write to become visible", map[string]any{
			"attempt": attempt, "status_code": APIStatusCode(err), "next_delay": interval,
		})
		select {
		case <-ctx.Done():
			return last, fmt.Errorf("resource write did not converge: %w", ctx.Err())
		case <-time.After(interval):
		}
		interval = min(interval*2, maxInterval)
	}
}
