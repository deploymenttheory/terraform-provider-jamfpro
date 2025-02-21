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
			"text_pane": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the text pane.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The display name of the text pane.",
						},
						"rank": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The rank/order of the text pane in the enrollment process.",
						},
						"title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The title of the text pane.",
						},
						"body": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The main content text of the pane.",
						},
						"subtext": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The subtext content of the pane.",
						},
						"back_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text for the back button.",
						},
						"continue_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text for the continue button.",
						},
					},
				},
			},
			"ldap_pane": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the LDAP pane.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The display name of the LDAP pane.",
						},
						"rank": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The rank/order of the LDAP pane in the enrollment process.",
						},
						"title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The title of the LDAP pane.",
						},
						"username_label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The label for the username field.",
						},
						"password_label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The label for the password field.",
						},
						"back_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text for the back button.",
						},
						"continue_button_text": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text for the continue button.",
						},
						"ldap_group_access": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the LDAP group.",
									},
									"ldap_server_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The ID of the LDAP server.",
									},
								},
							},
						},
					},
				},
			},
			"sso_pane": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the SSO pane.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The display name of the SSO pane.",
						},
						"rank": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The rank/order of the SSO pane in the enrollment process.",
						},
						"is_group_enrollment_access_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether group enrollment access is enabled.",
						},
						"group_enrollment_access_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The name of the group for enrollment access.",
						},
						"is_use_jamf_connect": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to use Jamf Connect.",
						},
						"short_name_attribute": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The attribute to use for short name.",
						},
						"long_name_attribute": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The attribute to use for long name.",
						},
					},
				},
			},
		},
	}
}
