// accounts_resource.go
package accounts

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/jamfprivileges"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAccount defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccounts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProAccountCreate,
		ReadContext:   resourceJamfProAccountReadWithCleanup,
		UpdateContext: resourceJamfProAccountUpdate,
		DeleteContext: resourceJamfProAccountDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the jamf pro account.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the jamf pro account.",
			},
			"directory_user": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is a directory user.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The full name of the account user.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the account user.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of the account user.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access status of the account (“enabled” or “disabled”).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "Enabled" || v == "Disabled" {
						return
					}

					return warns, append(errs, fmt.Errorf("%q must be either 'Enabled' or 'Disabled', got: %s", key, v))
				},
			},
			"identity_server_id": {
				Type:        schema.TypeInt,
				Description: "The Id of the identity server",
				Optional:    true,
			},
			"force_password_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is forced to change password on next login.",
			},
			"access_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The access level of the account. This can be either Full Access, scoped to a jamf pro site with Site Access, or scoped to a jamf pro account group with Group Access",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "Full Access" || v == "Site Access" || v == "Group Access" {
						return
					}
					return warns, append(errs, fmt.Errorf("%q must be either 'Full Access' or 'Site Access' or 'Group Access', got: %s", key, v))
				},
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the account.",
				Sensitive:   true,
			},
			"privilege_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The privilege set assigned to the account.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validPrivileges := []string{"Administrator", "Auditor", "Enrollment Only", "Custom"}
					for _, validPriv := range validPrivileges {
						if v == validPriv {
							return
						}
					}
					return warns, append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validPrivileges, v))
				},
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"jss_objects_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Objects.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: jamfprivileges.ValidateJSSObjectsPrivileges,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Settings.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: jamfprivileges.ValidateJSSSettingsPrivileges,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Actions.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: jamfprivileges.ValidateJSSActionsPrivileges,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Casper Admin.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: jamfprivileges.ValidateCasperAdminPrivileges,
				},
			},
			"casper_remote_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Casper Remote.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_imaging_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Casper Imaging.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recon_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Recon.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
