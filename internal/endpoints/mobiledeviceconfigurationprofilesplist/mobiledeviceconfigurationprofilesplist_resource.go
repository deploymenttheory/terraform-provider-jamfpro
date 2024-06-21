// mobiledeviceconfigurationprofilesplist_resource.go
package mobiledeviceconfigurationprofilesplist

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProMobileDeviceConfigurationProfilesPlist defines the schema for mobile device configuration profiles in Terraform.
func ResourceJamfProMobileDeviceConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProMobileDeviceConfigurationProfilePlistCreate,
		ReadContext:   resourceJamfProMobileDeviceConfigurationProfilePlistRead,
		UpdateContext: resourceJamfProMobileDeviceConfigurationProfilePlistUpdate,
		DeleteContext: resourceJamfProMobileDeviceConfigurationProfilePlistDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
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
				Description: "The unique identifier for the mobile device configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the mobile device configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the mobile device configuration profile.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level at which the mobile device configuration profile is applied, can be either 'Device Level' or 'User Level'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "Device Level" || v == "User Level" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Device Level' or 'User Level', got: %s", key, v))
					return warns, errs
				},
			},
			"site": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The site information associated with the mobile device configuration profile.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				}},
			},
			"category": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The jamf pro category information for the mobile device configuration profile.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The unique identifier for the Jamf Pro category.",
					},
				}},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The universally unique identifier for the profile.",
			},
			"deployment_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deployment method for the mobile device configuration profile, can be either 'Install Automatically' or 'Make Available in Self Service'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "Install Automatically" || v == "Make Available in Self Service" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Install Automatically' or 'Make Available in Self Service', got: %s", key, v))
					return warns, errs
				},
			},
			"redeploy_on_update": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Newly Assigned", // This is always "Newly Assigned" on existing profile objects, but may be set "All" on profile update requests and in TF state.
				Description: "Defines the redeployment behaviour when a mobile device config profile update occurs.This is always 'Newly Assigned' on new profile objects, but may be set 'All' on profile update requests and in TF state",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "All" || v == "Newly Assigned" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'All' or 'Newly Assigned', got: %s", key, v))
					return warns, errs
				},
			},
			"redeploy_days_before_cert_expires": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of days before certificate expiration when the profile should be redeployed.",
			},
			"payloads": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        plist.NormalizePayloadState,
				ValidateFunc:     plist.ValidatePayload,
				DiffSuppressFunc: DiffSuppressPayloads,
				Description:      "The iOS / iPadOS / tvOS configuration profile payload. Can be a file path to a .mobileconfig or a string with an embedded mobileconfig plist.",
			},
			// Scope
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the configuration profile.",
				Required:    true,
				Elem:        sharedschemas.GetSharedMobileDeviceSchemaScope(),
			},
		},
	}
}
