// Icons_resource.go
package icon

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProIcons defines the schema and RU operations for managing Jamf Pro icons in Terraform.
func ResourceJamfProIcons() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the icon. Returned by the Jamf Pro API.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the icon. Returned by the Jamf Pro API.",
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"icon_file_path": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				Description:      "The file path to the icon file (PNG) to be uploaded.",
				ValidateDiagFunc: validateIconFilePath(),
			},
			"icon_file_web_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "The web location of the icon file, can be a http(s) URL",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`), "Must be a valid URL."),
			},
			"icon_file_base64": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Base64 encoded string of the icon image file (PNG format). Must be a valid base64 encoded image.",
			},
		},
	}
}
