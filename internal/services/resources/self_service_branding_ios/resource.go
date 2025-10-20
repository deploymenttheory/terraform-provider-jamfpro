package self_service_branding_ios

// self_service_branding_ios_resource.go
import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProSelfServiceBrandingIOS defines the schema and CRUD operations for self-service branding (iOS).
func ResourceJamfProSelfServiceBrandingIOS() *schema.Resource {
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
			"main_header": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The main header for the branding configuration.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"icon_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Icon ID to use for the branding.",
			},
			"header_background_color_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Header background color code (RRGGBB, without leading '#').",
			},
			"menu_icon_color_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Menu icon color code (RRGGBB, without leading '#').",
			},
			"branding_name_color_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Branding name color code (RRGGBB, without leading '#').",
			},
			"status_bar_text_color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status bar text color; typically 'light' or 'dark'.",
			},
		},
	}
}
