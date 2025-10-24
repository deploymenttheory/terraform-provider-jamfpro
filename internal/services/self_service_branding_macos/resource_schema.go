package self_service_branding_macos

// self_service_branding_macos_resource.go
import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProSelfServiceBrandingMacOS defines the schema and CRUD operations for self-service branding (macOS).
func ResourceJamfProSelfServiceBrandingMacOS() *schema.Resource {
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
				Description: "The unique identifier of the Self Service branding configuration.",
			},
			"application_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application header for the branding configuration.",
			},
			"sidebar_heading": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The sidebar heading for the branding configuration.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"sidebar_subheading": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sidebar subheading for the branding configuration.",
			},
			"icon_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Icon ID to use for the branding.",
			},
			"home_page_banner_image_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Home page banner image ID used in the branding.",
			},
			"home_page_heading": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Home screen heading text for the self service branding.",
			},
			"home_page_subheading": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Home screen subheading text for the self service branding.",
			},
		},
	}
}
