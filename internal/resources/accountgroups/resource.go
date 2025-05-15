// accountgroups_resource.go
package accountgroups

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAccountGroup defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: customDiffAccountGroups,
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the account group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account group.",
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
					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					return warns, append(errs, fmt.Errorf("%q must be either 'Full Access' or 'Site Access' or 'Group Access', got: %s", key, v))
				},
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
					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					return warns, append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validPrivileges, v))
				},
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"jss_objects_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Objects.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Settings.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Actions.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Casper Admin.(DEPRECATED)",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_ids": {
				Type:        schema.TypeList,
				Description: "Accounts which should be a member of this group by ID",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"identity_server_id": {
				Type:        schema.TypeInt,
				Description: "The Id of the identity server",
				Optional:    true,
			},
		},
	}
}
