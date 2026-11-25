package sdk

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func TestNewClientReturnsSDKClient(t *testing.T) {
	// Services assert meta.(*jamfpro.Client) until they move to the interface.
	got := NewClient(nil)
	if _, ok := got.(*jamfpro.Client); !ok {
		t.Fatalf("NewClient(nil) = %T, want *jamfpro.Client", got)
	}
}
