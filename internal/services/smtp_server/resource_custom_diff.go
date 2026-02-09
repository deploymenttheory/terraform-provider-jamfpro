package smtp_server

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func customDiff(ctx context.Context, d *schema.ResourceDiff, m any) error {
	authType := d.Get("authentication_type").(string)
	var errorMessages []string

	// Common validations for all types except NONE
	if authType != "NONE" {
		if _, ok := d.GetOk("sender_settings"); !ok {
			errorMessages = append(errorMessages, "sender_settings is required for all authentication types")
		}
	}

	switch authType {
	case "BASIC":
		// Basic auth requires connection settings and basic auth credentials
		if _, ok := d.GetOk("connection_settings"); !ok {
			errorMessages = append(errorMessages, "connection_settings is required when authentication_type is BASIC")
		}
		if _, ok := d.GetOk("basic_auth_credentials"); !ok {
			errorMessages = append(errorMessages, "basic_auth_credentials is required when authentication_type is BASIC")
		}
		// Ensure other credential types are not set
		if _, ok := d.GetOk("graph_api_credentials"); ok {
			errorMessages = append(errorMessages, "graph_api_credentials cannot be set when authentication_type is BASIC")
		}
		if _, ok := d.GetOk("google_mail_credentials"); ok {
			errorMessages = append(errorMessages, "google_mail_credentials cannot be set when authentication_type is BASIC")
		}

	case "GRAPH_API":
		// Graph API requires sender settings and graph api credentials
		if _, ok := d.GetOk("graph_api_credentials"); !ok {
			errorMessages = append(errorMessages, "graph_api_credentials is required when authentication_type is GRAPH_API")
		}
		// Ensure other credential types are not set
		if _, ok := d.GetOk("basic_auth_credentials"); ok {
			errorMessages = append(errorMessages, "basic_auth_credentials cannot be set when authentication_type is GRAPH_API")
		}
		if _, ok := d.GetOk("google_mail_credentials"); ok {
			errorMessages = append(errorMessages, "google_mail_credentials cannot be set when authentication_type is GRAPH_API")
		}
		if _, ok := d.GetOk("connection_settings"); ok {
			errorMessages = append(errorMessages, "connection_settings cannot be set when authentication_type is GRAPH_API")
		}

	case "GOOGLE_MAIL":
		// Google Mail requires sender settings and google mail credentials
		if _, ok := d.GetOk("google_mail_credentials"); !ok {
			errorMessages = append(errorMessages, "google_mail_credentials is required when authentication_type is GOOGLE_MAIL")
		}
		// Ensure other credential types are not set
		if _, ok := d.GetOk("basic_auth_credentials"); ok {
			errorMessages = append(errorMessages, "basic_auth_credentials cannot be set when authentication_type is GOOGLE_MAIL")
		}
		if _, ok := d.GetOk("graph_api_credentials"); ok {
			errorMessages = append(errorMessages, "graph_api_credentials cannot be set when authentication_type is GOOGLE_MAIL")
		}
		if _, ok := d.GetOk("connection_settings"); ok {
			errorMessages = append(errorMessages, "connection_settings cannot be set when authentication_type is GOOGLE_MAIL")
		}

	case "NONE":
		// When type is NONE, ensure no credential types are set
		if _, ok := d.GetOk("basic_auth_credentials"); ok {
			errorMessages = append(errorMessages, "basic_auth_credentials cannot be set when authentication_type is NONE")
		}
		if _, ok := d.GetOk("graph_api_credentials"); ok {
			errorMessages = append(errorMessages, "graph_api_credentials cannot be set when authentication_type is NONE")
		}
		if _, ok := d.GetOk("google_mail_credentials"); ok {
			errorMessages = append(errorMessages, "google_mail_credentials cannot be set when authentication_type is NONE")
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

// validateGUID validates that a string is a properly formatted GUID/UUID
func validateGUID() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringMatch(
		regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
		"must be a valid GUID/UUID in the format 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'",
	))
}

// validateClientID validates that a string is either a GUID/UUID or a Google OAuth client ID
func validateClientID() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i any, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		// Check if it matches GUID/UUID format
		guidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		if guidPattern.MatchString(v) {
			return nil, nil
		}

		// Check if it matches Google OAuth client ID format
		googlePattern := regexp.MustCompile(`^[0-9]{12}-[0-9a-z]{32}\.apps\.googleusercontent\.com$`)
		if googlePattern.MatchString(v) {
			return nil, nil
		}

		return nil, []error{fmt.Errorf("%s must be either a valid GUID/UUID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx) or a Google OAuth client ID (012345678901-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com)", k)}
	})
}
