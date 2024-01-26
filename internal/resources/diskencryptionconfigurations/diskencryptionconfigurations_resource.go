// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	utils "github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"

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
			},
			"file_vault_enabled_users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines which user to enable for FileVault 2. Value can be either 'Management Account' or 'Current or Next User'",
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
		// Assuming the first element of the list is used for configuration
		institutionalRecoveryKeyMap := irk[0].(map[string]interface{})
		data := util.GetStringFromMap(institutionalRecoveryKeyMap, "data")

		// Base64 encode the 'data' value using Base64EncodeCertificate utility function
		encodedBase64CertData, err := utils.Base64EncodeCertificate(data)
		if err != nil {
			// Handle the error if base64 encoding fails
			log.Printf("[ERROR] Error encoding Data to Base64: %s", err)
			return nil, fmt.Errorf("error encoding Data to Base64: %v", err)
		}

		diskEncryptionConfig.InstitutionalRecoveryKey = &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{
			Key:             util.GetStringFromMap(institutionalRecoveryKeyMap, "key"),
			CertificateType: util.GetStringFromMap(institutionalRecoveryKeyMap, "certificate_type"),
			Password:        util.GetStringFromMap(institutionalRecoveryKeyMap, "password"),
			Data:            encodedBase64CertData,
		}
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
	var createdConfigResponse *jamfpro.ResponseDiskEncryptionConfiguration
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
// 1. Fetches the disk encryption configuration's current state using its ID.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the disk encryption configuration being deleted outside of Terraform, to keep the Terraform state synchronized.
// ResourceJamfProDiskEncryptionConfigurationsRead is responsible for reading the current state of a Jamf Pro Disk Encryption Configuration Resource from the remote system.
func ResourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var responseConfig *jamfpro.ResponseDiskEncryptionConfiguration

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		configID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse disk encryption configuration ID: %v", convertErr))
		}

		var apiErr error
		responseConfig, apiErr = conn.GetDiskEncryptionConfigurationByID(configID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				if apiError.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	if responseConfig != nil {
		d.Set("id", strconv.Itoa(responseConfig.ID))
		d.Set("name", responseConfig.Name)
		d.Set("key_type", responseConfig.KeyType)
		d.Set("file_vault_enabled_users", responseConfig.FileVaultEnabledUsers)

		if responseConfig.InstitutionalRecoveryKey != nil {
			irk := map[string]interface{}{
				"key":              responseConfig.InstitutionalRecoveryKey.Key,
				"certificate_type": responseConfig.InstitutionalRecoveryKey.CertificateType,
				"password":         responseConfig.InstitutionalRecoveryKey.Password,
				"data":             responseConfig.InstitutionalRecoveryKey.Data,
			}
			d.Set("institutional_recovery_key", []interface{}{irk})
		}
	} else {
		d.SetId("")
	}

	return diags
}

// ResourceJamfProDockItemsUpdate is responsible for updating an existing Jamf Pro Dock Item on the remote system.
func ResourceJamfProDockItemsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the dock item
		dockItem, err := constructDiskEncryptionConfiguration(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the department for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		dockItemID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to update the resource by ID
		_, apiErr := conn.UpdateDockItemByID(dockItemID, dockItem)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			dockItemName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			_, apiErr = conn.UpdateDockItemByName(dockItemName, dockItem)
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
		readDiags := ResourceJamfProDockItemsRead(ctx, d, meta)
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

// ResourceJamfProDockItemsDelete is responsible for deleting a Jamf Pro Dock Item.
func ResourceJamfProDockItemsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		dockItemID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to DELETE the resource
		apiErr := conn.DeleteDockItemByID(dockItemID)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			dockItemName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeleteDockItemByName(dockItemName)
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
