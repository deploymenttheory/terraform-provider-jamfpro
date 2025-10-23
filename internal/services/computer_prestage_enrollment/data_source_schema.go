package computer_prestage_enrollment

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerPrestageEnrollment provides information about a specific computer prestage enrollment in Jamf Pro.
func DataSourceJamfProComputerPrestageEnrollment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the computer prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the computer prestage.",
			},
		},
	}
}
