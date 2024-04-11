// mobiledeviceconfigurationprofiles_resource.go
package mobiledeviceconfigurationprofiles

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
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
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
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
			"distribution_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution method for the mobile device configuration profile, can be either 'Install Automatically' or 'Make Available in Self Service'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "Install Automatically" || v == "Make Available in Self Service" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Install Automatically' or 'Make Available in Self Service', got: %s", key, v))
					return warns, errs
				},
			},
			"user_removable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the mobile device configuration profile is removable by the user.",
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
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The universally unique identifier for the profile.",
			},
			"redeploy_on_update": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines if the profile should be redeployed when an update occurs.",
			},
			"payloads": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The payloads included in the configuration profile.",
			},
			// Scope
			"scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The scope in which the mobile device configuration profile is applied.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
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
					"mobile_devices": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "The list of specific mobile devices to which the profile is applied.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The unique identifier of the mobile device.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The name of the mobile device.",
							},
							"udid": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The UDID of the mobile device.",
							},
							"wifi_mac_address": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The WiFi MAC address of the mobile device.",
							},
						}},
					},
					"mobile_device_groups": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "The list of mobile device groups to which the profile is applied.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The unique identifier of the mobile device group.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The name of the mobile device group.",
							},
						}},
					},
					"buildings": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "The list of buildings to which the profile is applied.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The unique identifier of the building.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The name of the building.",
							},
						}},
					},
					"departments": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "The list of departments to which the profile is applied.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The unique identifier of the department.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The name of the department.",
							},
						}},
					},
					// Scope limitations and exclusions
					"limitations": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Restrictions on where or how the profile is applied within the scope.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"network_segments": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of network segments to which limitations apply.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the network segment.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the network segment.",
									},
								}},
							},
							"users": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of users to which limitations apply.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the user.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user.",
									},
								}},
							},
							"user_groups": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of user groups to which limitations apply.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the user group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user group.",
									},
								}},
							},
							"ibeacons": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of iBeacons to which limitations apply.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the iBeacon.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the iBeacon.",
									},
								}},
							},
						}},
					},
					"exclusions": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Items that are excluded from the scope.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"mobile_devices": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of mobile devices excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the mobile device.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the mobile device.",
									},
								}},
							},
							"mobile_device_groups": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of mobile device groups excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the mobile device group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the mobile device group.",
									},
								}},
							},
							"users": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of users excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the user.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user.",
									},
								}},
							},
							"user_groups": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of user groups excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the user group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user group.",
									},
								}},
							},
							"buildings": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of buildings excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the building.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the building.",
									},
								}},
							},
							"departments": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of departments excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the department.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the department.",
									},
								}},
							},
							"network_segments": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of network segments excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the network segment.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the network segment.",
									},
								}},
							},
							"jss_users": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of JSS users excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the JSS user.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the JSS user.",
									},
								}},
							},
							"jss_user_groups": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of JSS user groups excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the JSS user group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the JSS user group.",
									},
								}},
							},
							"ibeacons": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: "The list of iBeacons excluded from the scope.",
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique identifier of the iBeacon.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the iBeacon.",
									},
								}},
							},
						}},
					},
				}},
			},
			// Self Service settings
			"self_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The self-service settings for the configuration profile.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"install_button_text": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Custom text that appears on the installation button in Self Service.",
					},
					"self_service_description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The description of the configuration profile as it appears in Self Service.",
					},
					"force_users_to_view_description": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether to force the user to view the description in Self Service before installing the profile.",
					},
					"self_service_icon": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "The icon displayed for the configuration profile in Self Service.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The ID of the icon resource.",
							},
							"filename": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The filename of the icon resource.",
							},
							"uri": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The URI location of the icon resource.",
							},
						}},
					},
					"feature_on_main_page": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether to feature the profile on the main page of Self Service.",
					},
					"self_service_categories": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "The categories within Self Service in which the profile is displayed.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "The ID of the Self Service category.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The name of the Self Service category.",
							},
							"display_in": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Whether the profile is displayed in this Self Service category.",
							},
							"feature_in": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Whether the profile is featured in this Self Service category.",
							},
						}},
					},
					"notification": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The notification type for the configuration profile.",
					},
					"notification_subject": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The subject of the notification sent when the configuration profile is installed.",
					},
					"notification_message": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The message of the notification sent when the configuration profile is installed.",
					},
				}},
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
		return common.HandleResourceNotFoundError(err, d)
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
