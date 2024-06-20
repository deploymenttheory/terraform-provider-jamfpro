// computerprestages_resource.go
package computerprestageenrollments

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProComputerPrestageEnrollmentEnrollment defines the schema for managing Jamf Pro Computer Prestages in Terraform.
func ResourceJamfProComputerPrestageEnrollmentEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProComputerPrestageEnrollmentCreate,
		ReadContext:   resourceJamfProComputerPrestageEnrollmentRead,
		UpdateContext: resourceJamfProComputerPrestageEnrollmentUpdate,
		DeleteContext: resourceJamfProComputerPrestageEnrollmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Second),
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
				Description: "The unique identifier of the computer prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the computer prestage.",
			},
			"mandatory": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether the computer prestage is mandatory.",
			},
			"mdm_removable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the MDM profile is removable.",
			},
			"support_phone_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Support phone number for the organization.",
			},
			"support_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Support email address for the organization.",
			},
			"department": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The department the computer prestage is assigned to.",
			},
			"default_prestage": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if this is the default computer prestage enrollment configuration. If yes then new devices will be automatically assigned to this PreStage enrollment",
			},
			"enrollment_site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The jamf pro Site ID that computers will be added to during enrollment. Default is -1, aka not used.",
			},
			"keep_existing_site_membership": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if existing device site membership should be retained.",
			},
			"keep_existing_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if existing device location information should be retained.",
			},
			"require_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is required to provide username and password on computers with macOS 10.10 or later.",
			},
			"authentication_prompt": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The authentication prompt message displayed to the user during enrollment.",
			},
			"prevent_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if activation lock should be prevented.",
			},
			"enable_device_based_activation_lock": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if device-based activation lock should be enabled.",
			},
			"device_enrollment_program_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The device enrollment program instance ID.",
			},
			"skip_setup_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Selected items are not displayed in the Setup Assistant during macOS device setup within Apple Device Enrollment (ADE).",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biometric": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip biometric setup.",
						},
						"terms_of_address": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip terms of address setup.",
						},
						"file_vault": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip FileVault setup.",
						},
						"icloud_diagnostics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip iCloud diagnostics setup.",
						},
						"diagnostics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip diagnostics setup.",
						},
						"accessibility": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip accessibility setup.",
						},
						"apple_id": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Apple ID setup.",
						},
						"screen_time": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Screen Time setup.",
						},
						"siri": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Siri setup.",
						},
						"display_tone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Display Tone setup.",
						},
						"restore": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Restore setup.",
						},
						"appearance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Appearance setup.",
						},
						"privacy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Privacy setup.",
						},
						"payment": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Payment setup.",
						},
						"registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Registration setup.",
						},
						"tos": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Terms of Service setup.",
						},
						"icloud_storage": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip iCloud Storage setup.",
						},
						"location": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Skip Location setup.",
						},
					},
				},
			},
			"location_information": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Location information associated with the Jamf Pro computer prestage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The username for the location information.",
						},
						"realname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The real name associated with this location.",
						},
						"phone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The phone number associated with this location.",
						},
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The email address associated with this location.",
						},
						"room": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The room associated with this location.",
						},
						"position": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The position associated with this location.",
						},
						"department_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The computerPrestage ID associated with this location.",
							Default:     "-1",
						},
						"building_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The building ID associated with this location.",
							Default:     "-1",
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the location information.",
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version lock of the location information.",
						},
					},
				},
			},
			"purchasing_information": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Purchasing information associated with the computer prestage.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the purchasing information.",
						},
						"leased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is leased.",
						},
						"purchased": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates if the item is purchased.",
						},
						"apple_care_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AppleCare ID.",
						},
						"po_number": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase order number.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The vendor name.",
						},
						"purchase_price": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase price.",
						},
						"life_expectancy": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The life expectancy in years.",
						},
						"purchasing_account": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing account.",
						},
						"purchasing_contact": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchasing contact.",
						},
						"lease_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The lease date.",
						},
						"po_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The purchase order date.",
						},
						"warranty_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The warranty date.",
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version lock.",
						},
					},
				},
			},
			"anchor_certificates": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Base64 encoded PEM Certificates.",
			},
			"enrollment_customization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The enrollment customization ID.",
				Default:     "0",
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The language setting.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region setting.",
			},
			"auto_advance_setup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if setup should auto-advance.",
			},
			"install_profiles_during_setup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if profiles should be installed during setup.",
			},
			"prestage_installed_profile_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "IDs of profiles installed during prestage.",
			},
			"custom_package_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom package IDs.",
			},
			"custom_package_distribution_point_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom package distribution point ID.",
			},
			"enable_recovery_lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if recovery lock should be enabled.",
			},
			"recovery_lock_password_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The recovery lock password type.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					validTypes := map[string]bool{
						"MANUAL": true,
						"RANDOM": true,
					}
					if _, valid := validTypes[v]; !valid {
						errs = append(errs, fmt.Errorf("%q must be one of 'MANUAL', 'RANDOM', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"recovery_lock_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The recovery lock password.",
			},
			"rotate_recovery_lock_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the recovery lock password should be rotated.",
			},
			"profile_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile UUID.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The site ID.",
				Default:     "-1",
			},
			"version_lock": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The version lock.",
			},
			"account_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of Account Settings.",
						},
						"payload_configured": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the payload is configured.",
						},
						"local_admin_account_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the local admin account is enabled.",
						},
						"admin_username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The admin username.",
						},
						"admin_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The admin password.",
						},
						"hidden_admin_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the admin account is hidden.",
						},
						"local_user_managed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the local user is managed.",
						},
						"user_account_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of user account (ADMINISTRATOR, STANDARD, SKIP).",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := util.GetString(val)
								validTypes := map[string]bool{
									"ADMINISTRATOR": true,
									"STANDARD":      true,
									"SKIP":          true,
								}
								if _, valid := validTypes[v]; !valid {
									errs = append(errs, fmt.Errorf("%q must be one of 'ADMINISTRATOR', 'STANDARD', 'SKIP', got: %s", key, v))
								}
								return warns, errs
							},
						},
						"version_lock": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The version lock for account settings.",
						},
						"prefill_primary_account_info_feature_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if prefilling primary account info feature is enabled.",
						},
						"prefill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of prefill (CUSTOM, DEVICE_OWNER).",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := util.GetString(val)
								validTypes := map[string]bool{
									"CUSTOM":       true,
									"DEVICE_OWNER": true,
								}
								if _, valid := validTypes[v]; !valid {
									errs = append(errs, fmt.Errorf("%q must be one of 'CUSTOM', 'DEVICE_OWNER', got: %s", key, v))
								}
								return warns, errs
							},
						},
						"prefill_account_full_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Full name for the account to prefill.",
						},
						"prefill_account_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Username for the account to prefill.",
						},
						"prevent_prefill_info_from_modification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if prefill info is prevented from modification.",
						},
					},
				},
			},
		},
	}
}

// resourceJamfProComputerPrestageEnrollmentCreate is responsible for creating a new computer prestage in Jamf Pro with terraform.
// The function:
// 1. Constructs the computer prestage data using the provided Terraform configuration.
// 2. Calls the API to create the computer prestage in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created computer prestage.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProComputerPrestageEnrollmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProComputerPrestageEnrollment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Prestage Enrollment: %v", err))
	}

	var creationResponse *jamfpro.ResponseComputerPrestageCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateComputerPrestage(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Computer Prestage Enrollment '%s' after retries: %v", resource.DisplayName, err))
	}

	d.SetId(creationResponse.ID)

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	return client.GetComputerPrestageByID(id.(string))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Computer Prestage Enrollment", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, resourceJamfProComputerPrestageEnrollmentRead(ctx, d, meta)...)
}

// resourceJamfProComputerPrestageEnrollmentRead is responsible for reading the current state of a Building Resource from the remote system.
// The function:
// 1. Fetches the building's current state using its ID. If it fails, then obtain the building's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the building being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProComputerPrestageEnrollmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := client.GetComputerPrestageByID(resourceID)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// resourceJamfProComputerPrestageEnrollmentUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func resourceJamfProComputerPrestageEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProComputerPrestageEnrollment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Disk Computer Prestage for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateComputerPrestageByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Computer Prestage '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	return append(diags, resourceJamfProComputerPrestageEnrollmentRead(ctx, d, meta)...)
}

// resourceJamfProComputerPrestageEnrollmentDelete is responsible for deleting a Jamf Pro Computer Prestage.
func resourceJamfProComputerPrestageEnrollmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteComputerPrestageByID(resourceID)
		if apiErr != nil {
			resourceDisplayName := d.Get("display_name").(string)
			apiErrByDisplayName := client.DeleteComputerPrestageByName(resourceDisplayName)
			if apiErrByDisplayName != nil {
				return retry.RetryableError(apiErrByDisplayName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Computer Prestage '%s' (ID: %s) after retries: %v", d.Get("display_name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
