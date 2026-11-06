package framework_crud

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
)

var errTestUnrelatedGroup = errors.New("permission denied for group 404")

func TestIsNotFound(t *testing.T) {
	for _, code := range []int{404, 410, 403, 500} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			err := fmt.Errorf("read group: %w", &apiresponse.APIError{StatusCode: code, Message: "group 404"})
			if IsNotFound(err) != (code == 404 || code == 410) {
				t.Fatalf("incorrect missing-resource classification for status %d", code)
			}
		})
	}
	if IsNotFound(errTestUnrelatedGroup) || IsNotFound(nil) {
		t.Fatal("untyped errors must not remove resources from state")
	}
}
