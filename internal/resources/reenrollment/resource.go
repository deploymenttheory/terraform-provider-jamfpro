package reenrollment

import (
	// "context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceReenrollmentSettings() *schema.Resource {
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
			"flush_location_information": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Clears computer and mobile device information from the User and Location category on the Inventory tab in inventory information during re-enrollment.",
			},
			"flush_location_information_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Clears computer and mobile device information from the User and Location History category on the History tab in inventory information during re-enrollment.",
			},
			"flush_policy_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Clears the logs for policies that ran on the computer and clears computer information from the Policy Logs category on the History tab in inventory information during re-enrollment.",
			},
			"flush_extension_attributes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Clears all values for extension attributes from computer and mobile device inventory information during re-enrollment. This does not apply to extension attributes populated by scripts or Directory Service Attribute Mapping.",
			},
			"flush_software_update_plans": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Clears all values for software update plans from computer and mobile device inventory information during re-enrollment.",
			},
			"flush_mdm_queue": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Clears computer and mobile device information from the Management History category on the History tab in inventory information during re-enrollment. Valid values are DELETE_NOTHING, DELETE_ERRORS, DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED, or DELETE_EVERYTHING.",
			},
		},
	}
}
