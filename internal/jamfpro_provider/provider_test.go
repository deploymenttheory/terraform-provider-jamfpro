package jamfpro_provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jamfpro": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Check if environment variables are set
	if v := os.Getenv("JAMFPRO_INSTANCE_NAME"); v == "" {
		t.Fatal("JAMFPRO_INSTANCE_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("JAMFPRO_CLIENT_ID"); v == "" {
		t.Fatal("JAMFPRO_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("JAMFPRO_CLIENT_SECRET"); v == "" {
		t.Fatal("JAMFPRO_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("JAMFPRO_DEBUG_MODE"); v == "" {
		t.Fatal("JAMFPRO_DEBUG_MODE must be set to true or false for acceptance tests")
	}
	// ... add any other necessary checks
}
