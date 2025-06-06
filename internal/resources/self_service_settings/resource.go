package self_service_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProSelfServiceSettings defines the schema and CRUD operations for managing Jamf Pro Self Service settings in Terraform.
func ResourceJamfProSelfServiceSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"install_automatically": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to install Self Service automatically",
			},
			"install_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location where Self Service should be installed",
			},
			"user_login_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The login level required for users. Set to 'NotRequired', 'Anonymous', or 'Required'",
			},
			"allow_remember_me": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to allow the Remember Me option",
			},
			"use_fido2": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable FIDO2 authentication",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of authentication to use. Set to 'Basic' or 'Saml'",
			},
			"notifications_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable notifications",
			},
			"alert_user_approved_mdm": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to alert users about approved MDM status",
			},
			"default_landing_page": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default landing page in Self Service. Set to 'HOME', 'BROWSE', 'HISTORY', or 'NOTIFICATIONS'",
			},
			"default_home_category_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the default home category. Set to -1 for no default category",
			},
			"bookmarks_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name to use for bookmarks. Set to 'Bookmarks' by default",
			},
		},
	}
}
