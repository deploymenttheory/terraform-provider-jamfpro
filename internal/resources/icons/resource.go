// Icons_resource.go
package icons

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProIcons defines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProIcons() *schema.Resource {
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "The file path to the icon file (PNG) to be uploaded.",
				ValidateFunc: validateIconFilePath(),
			},
			"icon_file_web_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "The web location of the icon file, can be a http(s) URL",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`), "Must be a valid URL."),
			},
		},
	}
}
