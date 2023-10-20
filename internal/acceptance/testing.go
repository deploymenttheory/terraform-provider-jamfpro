package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func PreCheck(t *testing.T) {
	variables := []string{
		"JAMFPRO_INSTANCE",
		"JAMFPRO_CLIENT_ID",
		"JAMFPRO_CLIENT_SECRET",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}

	instanceName := os.Getenv("JAMFPRO_INSTANCE")
	if instanceName == "" {
		t.Fatal("instance_name must be provided either as an environment variable (JAMFPRO_INSTANCE) or in the Terraform configuration")
	}

	clientID := os.Getenv("JAMFPRO_CLIENT_ID")
	if clientID == "" {
		t.Fatal("client_id must be provided either as an environment variable (JAMFPRO_CLIENT_ID) or in the Terraform configuration")
	}

	clientSecret := os.Getenv("JAMFPRO_CLIENT_SECRET")
	if clientSecret == "" {
		t.Fatal("client_secret must be provided either as an environment variable (JAMFPRO_CLIENT_SECRET) or in the Terraform configuration")
	}
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
