// apiroles_resource.go
package apiroles

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAPIRoles defines the schema for managing Jamf Pro API Roles in Terraform.
func ResourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Second),
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
				Description: "The unique identifier of the Jamf API Role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Jamf API Role.",
			},
			"privileges": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "List of api role privileges associated with the Jamf API Role. These are compared against the Jamf Pro" +
					"server for validation. You must supply the exact privilege names as they appear in the Jamf Pro server. They are case-sensitive" +
					"and can and do change between Jamf Pro versions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
