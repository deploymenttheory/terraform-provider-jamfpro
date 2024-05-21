// mobiledeviceconfigurationprofiles_resource.go
package mobiledeviceconfigurationprofiles

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMobileDeviceConfigurationProfile defines the schema for mobile device configuration profiles in Terraform.
func ResourceJamfProMobileDeviceConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMobileDeviceConfigurationProfileCreate,
		ReadContext:   ResourceJamfProMobileDeviceConfigurationProfileRead,
		UpdateContext: ResourceJamfProMobileDeviceConfigurationProfileUpdate,
		DeleteContext: ResourceJamfProMobileDeviceConfigurationProfileDelete,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The iOS / iPadOS / tvOS configuration profile payload. Can be a file path to a .mobileconfig or a string with an embedded mobileconfig plist.",
			},
			// Scope
			"scope": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope in which the mobile device configuration profile is applied.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_mobile_devices": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, the profile is applied to all mobile devices.",
						},
						"all_jss_users": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, the profile is applied to all JSS users.",
						},
						"mobile_device_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of mobile device IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"mobile_device_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of mobile device group IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"building_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of building IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"department_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of department IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"jss_user_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of JSS user IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"jss_user_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of JSS user group IDs associated with the profile.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Limitations for the profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_segment_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of network segment IDs for limitations.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of iBeacon IDs for limitations.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"directory_service_or_local_usernames": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of directory service / local usernames for scoping limitations.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"user_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of user group IDs for limitations.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"exclusions": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Exclusions for the profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mobile_device_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of mobile device IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"mobile_device_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of mobile device group IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"building_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of building IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"department_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of department IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"network_segment_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of network segment IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"jss_user_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of user names for exclusions.",
									},
									"jss_user_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of JSS user group IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of iBeacon IDs for exclusions.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
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

// ResourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
func ResourceJamfProMobileDeviceConfigurationProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProMobileDeviceConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateMobileDeviceConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Mobile Device Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetMobileDeviceConfigurationProfileByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Mobile Device Configuration Profile", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)

	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMobileDeviceConfigurationProfileRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProMobileDeviceConfigurationProfileRead is responsible for reading the current state of a Jamf Pro Mobile Device Configuration Profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProMobileDeviceConfigurationProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := conn.GetMobileDeviceConfigurationProfileByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Define fields to remove from the plist data
	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization"}

	// Process the payload from the server
	_, serverPayloadHash, err := configurationprofiles.ProcessConfigurationProfile(resource.General.Payloads, fieldsToRemove)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to process configuration profile: %v", err))
	}

	// Retrieve the stored payload hash from the state
	storedPayloadHash := d.Get("payloads").(string)

	// Log the two hashes being compared
	log.Printf("[DEBUG] Comparing payload hashes: stored hash = %s, server hash = %s", storedPayloadHash, serverPayloadHash)

	// Compare the hash with the stored hash in the state
	if storedPayloadHash != serverPayloadHash {
		// If hashes do not match, update the payloads field in the state
		if err := d.Set("payloads", serverPayloadHash); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set payload hash in state: %v", err))...)
		}
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProMobileDeviceConfigurationProfileUpdate is responsible for updating an existing Jamf Pro Mobile Device Configuration Profile on the remote system.
func ResourceJamfProMobileDeviceConfigurationProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProMobileDeviceConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateMobileDeviceConfigurationProfileByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", resource.General.Name, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMobileDeviceConfigurationProfileRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProMobileDeviceConfigurationProfileDelete is responsible for deleting a Jamf Pro Mobile Device Configuration Profile.
func ResourceJamfProMobileDeviceConfigurationProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteMobileDeviceConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteMobileDeviceConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
