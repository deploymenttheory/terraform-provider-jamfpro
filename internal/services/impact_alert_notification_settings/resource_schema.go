package impact_alert_notification_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceImpactAlertNotificationSettings defines the schema and CRUD operations for managing Jamf Pro Impact Alert Notification Settings in Terraform.
func ResourceImpactAlertNotificationSettings() *schema.Resource {
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
			"deployable_objects_alert_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Display deployment impact alert on Save for deployable objects. Jamf Pro users will be prompted with a deployment summary if edits are made to a deployable object.",
			},
			"deployable_objects_confirmation_code_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Require Jamf Pro users to acknowledge edits for deployable objects. Jamf Pro users will be prompted to type COMMIT in the criteria summary before saving.",
			},
			"scopeable_objects_alert_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Display criteria impact alert on Save for scopeable object edits. Jamf Pro users will be prompted with a deployment summary if edits are made to a scopeable object.",
			},
			"scopeable_objects_confirmation_code_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Require Jamf Pro users to acknowledge edits for scopeable objects. Jamf Pro users will be prompted to type COMMIT in the criteria summary before saving.",
			},
		},
	}
}
