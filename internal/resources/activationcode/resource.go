// activationcode_resource.go
package activationcode

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProActivationCode defines the schema and CRUD operations for managing Jamf Pro activation code configuration in Terraform.
func ResourceJamfProActivationCode() *schema.Resource {
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
			"organization_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization associated with the activation code.",
			},
			"code": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The activation code.",
			},
			thisisntwihsw0pgnwpidfgnwnbiofipqwe
		},
	}
}
