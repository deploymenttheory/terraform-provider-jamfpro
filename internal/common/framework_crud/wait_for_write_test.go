package framework_crud

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
	"github.com/stretchr/testify/require"
)

type writeReadBackValue struct {
	version int
}

var errTestWriteReadBackTransport = errors.New("connection closed")

func TestWaitForWriteRetriesStaleAndTransientReads(t *testing.T) {
	for _, status := range []int{404, 410, 429, 500, 502, 503, 504} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stale, desired := &writeReadBackValue{version: 1}, &writeReadBackValue{version: 2}
			reads := 0
			got, err := WaitForWrite(ctx, func(readCtx context.Context) (*writeReadBackValue, error) {
				_, deadline := readCtx.Deadline()
				require.True(t, deadline)
				reads++
				switch reads {
				case 1:
					return stale, nil
				case 2:
					return nil, fmt.Errorf("read group: %w", &apiresponse.APIError{StatusCode: status})
				case 3:
					return nil, nil //nolint:nilnil // Exercise an empty successful SDK response before the write becomes visible.
				default:
					return desired, nil
				}
			}, func(value *writeReadBackValue) bool { return value.version == 2 }, time.Nanosecond)
			require.NoError(t, err)
			require.Same(t, desired, got)
			require.Equal(t, 4, reads)
		})
	}
}

func TestWaitForWriteRetainsObservedValueOnTerminalError(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
	}{
		{name: "permission denied", err: &apiresponse.APIError{StatusCode: 403}},
		{name: "invalid request", err: &apiresponse.APIError{StatusCode: 400}},
		{name: "transport failure", err: errTestWriteReadBackTransport},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stale := &writeReadBackValue{version: 1}
			reads := 0
			got, err := WaitForWrite(ctx, func(context.Context) (*writeReadBackValue, error) {
				reads++
				if reads == 1 {
					return stale, nil
				}
				return nil, tt.err
			}, func(*writeReadBackValue) bool { return false }, time.Nanosecond)
			require.ErrorIs(t, err, tt.err)
			require.Same(t, stale, got)
			require.Equal(t, 2, reads)
		})
	}
}

func TestWaitForWriteDeadlineRetainsObservedValue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	stale := &writeReadBackValue{version: 1}
	got, err := WaitForWrite(ctx, func(context.Context) (*writeReadBackValue, error) {
		return stale, nil
	}, func(*writeReadBackValue) bool { return false }, time.Hour)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Same(t, stale, got)
}

func TestWaitForWriteRequiresOperationDeadline(t *testing.T) {
	got, err := WaitForWrite(context.Background(), func(context.Context) (*writeReadBackValue, error) {
		t.Fatal("must not start an unbounded read-back")
		return &writeReadBackValue{}, nil
	}, func(*writeReadBackValue) bool { return false }, time.Nanosecond)
	require.ErrorIs(t, err, errWriteReadBackNeedsDeadline)
	require.Nil(t, got)
}
