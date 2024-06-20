// activationcode_resource.go
package activationcode

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProActivationCode defines the schema and CRUD operations for managing Jamf Pro activation code configuration in Terraform.
func ResourceJamfProActivationCode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProActivationCodeCreate,
		ReadContext:   resourceJamfProActivationCodeRead,
		UpdateContext: resourceJamfProActivationCodeUpdate,
		DeleteContext: resourceJamfProActivationCodeDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
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
				Description: "The activation code.",
			},
		},
	}
}
