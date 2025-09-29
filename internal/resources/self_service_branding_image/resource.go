// self_service_branding_image_resource.go
package self_service_branding_image

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProSelfServiceBrandingImage defines the schema and RU operations for managing Jamf Pro self service branding images in Terraform.
func ResourceJamfProSelfServiceBrandingImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Self Service branding image. Derived from the URL.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the uploaded Self Service branding image.",
			},
			"self_service_branding_image_file_path": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ForceNew:         true,
				Description:      "The file path to the Self Service branding image file (PNG) to be uploaded.",
				ValidateDiagFunc: validateImageFilePath(),
			},
			"self_service_branding_image_file_web_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				Description:  "The web location of the Self Service branding image file, can be a http(s) URL",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`), "Must be a valid URL."),
			},
		},
	}
}
