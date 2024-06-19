// mobiledeviceconfigurationprofilesplist_resource.go
package mobiledeviceconfigurationprofilesplist

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMobileDeviceConfigurationProfilesPlist defines the schema for mobile device configuration profiles in Terraform.
func ResourceJamfProMobileDeviceConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMobileDeviceConfigurationProfilePlistCreate,
		ReadContext:   ResourceJamfProMobileDeviceConfigurationProfilePlistRead,
		UpdateContext: ResourceJamfProMobileDeviceConfigurationProfilePlistUpdate,
		DeleteContext: ResourceJamfProMobileDeviceConfigurationProfilePlistDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
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

// ResourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProMobileDeviceConfigurationProfilePlistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProMobileDeviceConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile: %v", err))
	}

	var creationResponse *jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateMobileDeviceConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Mobile Device Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetMobileDeviceConfigurationProfileByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Mobile Device Configuration Profile", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx, d, meta)...)
}

// ResourceJamfProMobileDeviceConfigurationProfilePlistRead is responsible for reading the current state of a Jamf Pro Mobile Device Configuration Profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetMobileDeviceConfigurationProfileByID(resourceIDInt)
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProMobileDeviceConfigurationProfilePlistUpdate is responsible for updating an existing Jamf Pro Mobile Device Configuration Profile on the remote system.
func ResourceJamfProMobileDeviceConfigurationProfilePlistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProMobileDeviceConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateMobileDeviceConfigurationProfileByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", resource.General.Name, resourceID, err))
	}

	return append(diags, ResourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx, d, meta)...)
}

// ResourceJamfProMobileDeviceConfigurationProfilePlistDelete is responsible for deleting a Jamf Pro Mobile Device Configuration Profile.
func ResourceJamfProMobileDeviceConfigurationProfilePlistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteMobileDeviceConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteMobileDeviceConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
