package enrollmentcustomizations

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProEnrollmentCustomization defines the schema and CRUD operations for managing Jamf Pro Enrollment Customizations in Terraform.
func ResourceJamfProEnrollmentCustomization() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
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
				Description: "The unique identifier of the enrollment customization.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
				Description: "The ID of the site associated with the enrollment customization.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the enrollment customization.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the enrollment customization.",
			},
			"enrollment_customization_image_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The .png image source file for upload to the enrollment customization. Recommended: 180x180 pixels and GIF or PNG format",
			},
			"branding_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"text_color": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The text color in hexadecimal format (6 characters, no # prefix).",
							ValidateFunc: validateHexColor,
						},
						"button_color": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The button color in hexadecimal format (6 characters, no # prefix).",
							ValidateFunc: validateHexColor,
						},
						"button_text_color": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The button text color in hexadecimal format (6 characters, no # prefix).",
							ValidateFunc: validateHexColor,
						},
						"background_color": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The background color in hexadecimal format (6 characters, no # prefix).",
							ValidateFunc: validateHexColor,
						},
						"icon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the icon image. the format must be 'https://your_jamfUrl/api/v2/enrollment-customizations/images/1'",
						},
					},
				},
			},
		},
	}
}
