// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProDiskEncryptionConfigurations defines the schema and CRUD operations for managing Jamf Pro Disk Encryption Configurations in Terraform.
func ResourceJamfProDiskEncryptionConfigurations() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProDiskEncryptionConfigurationsCreate,
		ReadContext:   ResourceJamfProDiskEncryptionConfigurationsRead,
		UpdateContext: ResourceJamfProDiskEncryptionConfigurationsUpdate,
		DeleteContext: ResourceJamfProDiskEncryptionConfigurationsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the disk encryption configuration.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the disk encryption configuration.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the key used in the disk encryption which can be either 'Institutional' or 'Individual and Institutional'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					validValues := []string{"Individual", "Institutional", "Individual and Institutional"}

					found := false
					for _, v := range validValues {
						if value == v {
							found = true
							break
						}
					}

					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of [%s], got '%s'", key, strings.Join(validValues, ", "), value))
					}

					return
				},
			},
			"file_vault_enabled_users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines which user to enable for FileVault 2. Value can be either 'Management Account' or 'Current or Next User'",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					validValues := []string{"Management Account", "Current or Next User"}

					found := false
					for _, v := range validValues {
						if value == v {
							found = true
							break
						}
					}

					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of [%s], got '%s'", key, strings.Join(validValues, ", "), value))
					}

					return
				},
			},
			"institutional_recovery_key": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Details of the institutional recovery key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Description: "The type of certificate used for the institutional recovery key. e.g 'PKCS12' for .p12 certificate types.",
							Optional:    true,
						},
						"password": {
							Type:        schema.TypeString,
							Description: "The password for the institutional recovery key certificate.",
							Optional:    true,
							Sensitive:   true,
						},
						"data": {
							Type:        schema.TypeString,
							Description: "The certificate payload.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// constructDiskEncryptionConfiguration constructs a ResourceDiskEncryptionConfiguration object from the provided schema data.
func constructDiskEncryptionConfiguration(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
	diskEncryptionConfig := &jamfpro.ResourceDiskEncryptionConfiguration{}

	// Utilize type assertion helper functions for direct field extraction
	diskEncryptionConfig.Name = util.GetStringFromInterface(d.Get("name"))
	diskEncryptionConfig.KeyType = util.GetStringFromInterface(d.Get("key_type"))
	diskEncryptionConfig.FileVaultEnabledUsers = util.GetStringFromInterface(d.Get("file_vault_enabled_users"))

	// Handling the institutional_recovery_key which is a list of maps
	if irk, ok := d.Get("institutional_recovery_key").([]interface{}); ok && len(irk) > 0 {
		institutionalRecoveryKeyMap := irk[0].(map[string]interface{})
		// Do not need to base64 as within tf you use the filebase64 method when referencing the certificate.
		certificatePayloadData := util.GetStringFromMap(institutionalRecoveryKeyMap, "data")

		diskEncryptionConfig.InstitutionalRecoveryKey = &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{
			Key:             util.GetStringFromMap(institutionalRecoveryKeyMap, "key"),
			CertificateType: util.GetStringFromMap(institutionalRecoveryKeyMap, "certificate_type"),
			Password:        util.GetStringFromMap(institutionalRecoveryKeyMap, "password"),
			Data:            certificatePayloadData,
		}
	} else {
		// Set InstitutionalRecoveryKey to nil or a default value if it's not provided
		diskEncryptionConfig.InstitutionalRecoveryKey = nil
	}

	// Marshal the search object into XML for logging
	xmlData, err := xml.MarshalIndent(diskEncryptionConfig, "", "  ")
	if err != nil {
		// Handle the error if XML marshaling fails
		log.Printf("[ERROR] Error marshaling DiskEncryptionConfiguration object to XML: %s", err)
		return nil, fmt.Errorf("error marshaling DiskEncryptionConfiguration object to XML: %v", err)
	}

	// Log the XML formatted search object
	tflog.Debug(ctx, fmt.Sprintf("Constructed DiskEncryptionConfiguration Object:\n%s", string(xmlData)))

	log.Printf("[INFO] Successfully constructed DiskEncryptionConfiguration with name: %s", diskEncryptionConfig.Name)

	return diskEncryptionConfig, nil
}

// Helper function to generate diagnostics based on the error type.
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsCreate is responsible for creating a new Jamf Pro Disk Encryption Configuration in the remote system.
// The function:
// 1. Constructs the disk encryption configuration data using the provided Terraform configuration.
// 2. Calls the API to create the disk encryption configuration in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created disk encryption configuration.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDiskEncryptionConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation.
	var createdConfigResponse *jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the disk encryption configuration.
		diskEncryptionConfig, err := constructDiskEncryptionConfiguration(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the disk encryption configuration for terraform create: %w", err))
		}

		// Directly call the API to create the resource.
		createdConfigResponse, err = conn.CreateDiskEncryptionConfiguration(diskEncryptionConfig)
		if err != nil {
			// Check if the error is an APIError.
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable.
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdConfigResponse.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsRead is responsible for reading the current state of a Jamf Pro Disk Encryption Configuration Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diskEncryptionConfig *jamfpro.ResourceDiskEncryptionConfiguration

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Try fetching the disk encryption configuration using the ID
		diskEncryptionConfig, err = conn.GetDiskEncryptionConfigurationByID(configID)
		if err != nil {
			// Handle the APIError
			if apiError, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			configName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			diskEncryptionConfig, err = conn.GetDiskEncryptionConfigurationByName(configName)
			if err != nil {
				// Handle the APIError
				if apiError, ok := err.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(err)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with disk encryption configuration attributes
	// Check if the InstitutionalRecoveryKey is nil or empty
	if diskEncryptionConfig.InstitutionalRecoveryKey == nil ||
		(diskEncryptionConfig.InstitutionalRecoveryKey.Key == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.Password == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.Data == "") {

		// If InstitutionalRecoveryKey is nil or empty, ensure it is not set in the Terraform state
		d.Set("institutional_recovery_key", []interface{}{})
	} else {
		// If InstitutionalRecoveryKey has data, set it in the Terraform state
		irk := make(map[string]interface{})
		// Removing as this causes state management false positives.
		//irk["key"] = diskEncryptionConfig.InstitutionalRecoveryKey.Key
		irk["certificate_type"] = diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType
		irk["password"] = diskEncryptionConfig.InstitutionalRecoveryKey.Password
		irk["data"] = diskEncryptionConfig.InstitutionalRecoveryKey.Data

		d.Set("institutional_recovery_key", []interface{}{irk})
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsUpdate is responsible for updating an existing Jamf Pro Disk Encryption Configuration on the remote system.
func ResourceJamfProDiskEncryptionConfigurationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}

	// The apiclient, which is of type *client.APIClient, holds a reference to the Jamf Pro client in its Conn field.
	// By assigning apiclient.Conn to jamfProClient, we are obtaining the actual Jamf Pro client (*jamfpro.Client)
	// that will be used for making API calls to the Jamf Pro server.
	jamfProClient := apiclient.Conn

	// Use the retry function for the update operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated disk encryption configuration
		diskEncryptionConfig, err := constructDiskEncryptionConfiguration(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the disk encryption configuration for terraform update: %w", err))
		}

		// Obtain the ID from the Terraform state to be used for the API request
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Directly call the API to update the resource
		_, apiErr := jamfProClient.UpdateDiskEncryptionConfigurationByID(configID, diskEncryptionConfig)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			configName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string in update"))
			}

			_, apiErr = jamfProClient.UpdateDiskEncryptionConfigurationByName(configName, diskEncryptionConfig)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func ResourceJamfProDiskEncryptionConfigurationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the DELETE operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		diskEncryptionConfigurationID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to DELETE the resource
		apiErr := conn.DeleteDiskEncryptionConfigurationByID(diskEncryptionConfigurationID)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			diskEncryptionConfigurationName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeleteDiskEncryptionConfigurationByName(diskEncryptionConfigurationName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
