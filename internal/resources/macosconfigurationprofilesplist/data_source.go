package macosconfigurationprofilesplist

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSConfigurationProfilesPlist provides information about a specific department in Jamf Pro.
func DataSourceJamfProMacOSConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Jamf UI name for configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the configuration profile.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The universally unique identifier for the profile.",
			},
			"site_id":     sharedschemas.GetSharedSchemaSite(),
			"category_id": sharedschemas.GetSharedSchemaCategory(),
			"distribution_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distribution method for the configuration profile. ['Make Available in Self Service','Install Automatically']",
			},
			"user_removable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the configuration profile is user removeable or not.",
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment level of the configuration profile. Available options are: 'User' or 'System'. Note: 'System' is mapped to 'Computer Level' in the Jamf Pro GUI.",
			},
			"payloads": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The macOS configuration profile payload. Can be a file path to a .mobileconfig or a string with an embedded mobileconfig plist." +
					"Jamf Pro stores configuration profiles as XML property lists (plists).",
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Defines the redeployment behaviour when an update to a macOS config profile" +
					"occurs. This is always 'Newly Assigned' on new profile objects, but may be set to 'All'" +
					"on profile update requests once the configuration profile has been deployed to at least" +
					" one device.",
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "The scope of the configuration profile.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
			"self_service": {
				Type:        schema.TypeList,
				Description: "Self Service Configuration",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self_service_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name for the profile in Self Service (Self Service 10.0.0 or later)",
						},
						"install_button_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name for the button that users click to install the profile",
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description to display for the profile in Self Service",
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Force users to view the description before the profile installs",
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Shows Configuration Profile on Self Service main page",
						},
						"notification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enables Notification for this profile in self service",
						},
						"notification_subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message Subject",
						},
						"notification_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message body",
						},
						"self_service_icon_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Icon for policy to use in self-service. Can be used in conjection with the icons resource",
						},
						"self_service_category": {
							Type:        schema.TypeSet,
							Description: "Self Service category options",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Description: "ID of category. Both ID and Name are required",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of category. Both ID and Name are required",
										Computed:    true,
									},
									"display_in": {
										Type:        schema.TypeBool,
										Description: "Display this profile in this category?",
										Computed:    true,
									},
									"feature_in": {
										Type:        schema.TypeBool,
										Description: "Feature this profile in this category?",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// dataSourceRead fetches the details of a macOS configuration profile.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id != "" && name != "" {

		return diag.FromErr(fmt.Errorf("please provide either 'id' or 'name', not both"))
	}

	var getFunc func() (*jamfpro.ResourceMacOSConfigurationProfile, error)
	var identifier string

	switch {
	case id != "":
		getFunc = func() (*jamfpro.ResourceMacOSConfigurationProfile, error) {
			return client.GetMacOSConfigurationProfileByID(id)
		}
		identifier = id
	case name != "":
		getFunc = func() (*jamfpro.ResourceMacOSConfigurationProfile, error) {
			return client.GetMacOSConfigurationProfileByName(name)
		}
		identifier = name
	default:

		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceMacOSConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Dock Item Configuration resource with identifier '%s' after retries: %v", identifier, err))
	}

	if resource == nil {
		d.SetId("")

		return diag.FromErr(fmt.Errorf("the Jamf Pro Dock Item Configuration resource was not found using identifier '%s'", identifier))
	}

	d.SetId(strconv.Itoa(resource.General.ID))
	return updateState(d, resource)
}
