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
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
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
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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

const (
	JamfProResourceDiskEncryptionConfiguration = "Disk Encryption Configuration"
)

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

// ResourceJamfProDiskEncryptionConfigurationsCreate is responsible for creating a new Jamf Pro Disk Encryption Configuration in the remote system.
// The function:
// 1. Constructs the disk encryption configuration data using the provided Terraform configuration.
// 2. Calls the API to create the disk encryption configuration in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created disk encryption configuration.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDiskEncryptionConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the object outside the retry loop to avoid reconstructing it on each retry
	diskEncryptionConfig, err := constructDiskEncryptionConfiguration(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDiskEncryptionConfiguration, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDiskEncryptionConfiguration)

	// Retry the API call to create the resource in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateDiskEncryptionConfiguration(diskEncryptionConfig)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceName, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceDiskEncryptionConfiguration, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the site and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceDiskEncryptionConfiguration, strconv.Itoa(creationResponse.ID))

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the site to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDiskEncryptionConfigurationsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the site, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, d.Id())
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsRead is responsible for reading the current state of a Jamf Pro Disk Encryption Configuration resource from Jamf Pro and updating the Terraform state with the retrieved data.
func ResourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var apiErrorCode int
	var diskEncryptionConfig *jamfpro.ResourceDiskEncryptionConfiguration

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDiskEncryptionConfiguration, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		diskEncryptionConfig, apiErr = conn.GetDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the account group, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID)
		return diag.FromErr(err)
	}

	// Assuming successful retrieval, proceed to set the resource attributes in Terraform state
	d.SetId(strconv.Itoa(resourceIDInt)) // Update the ID in the state
	d.Set("name", diskEncryptionConfig.Name)
	d.Set("key_type", diskEncryptionConfig.KeyType)
	d.Set("file_vault_enabled_users", diskEncryptionConfig.FileVaultEnabledUsers)

	// Institutional Recovery Key
	if diskEncryptionConfig.InstitutionalRecoveryKey == nil ||
		(diskEncryptionConfig.InstitutionalRecoveryKey.Key == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType == "" &&
			//diskEncryptionConfig.InstitutionalRecoveryKey.Password == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.Data == "") {

		// If InstitutionalRecoveryKey is nil or empty, ensure it is not set in the Terraform state
		d.Set("institutional_recovery_key", []interface{}{})
	} else {
		// If InstitutionalRecoveryKey has data, set it in the Terraform state
		irk := make(map[string]interface{})
		irk["certificate_type"] = diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType
		//irk["password"] = diskEncryptionConfig.InstitutionalRecoveryKey.Password
		irk["data"] = diskEncryptionConfig.InstitutionalRecoveryKey.Data

		d.Set("institutional_recovery_key", []interface{}{irk})
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsUpdate is responsible for updating an existing Jamf Pro Disk Encryption Configuration on the remote system.
func ResourceJamfProDiskEncryptionConfigurationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDiskEncryptionConfiguration, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Construct the resource object
	diskEncryptionConfig, err := constructDiskEncryptionConfiguration(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDiskEncryptionConfiguration, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDiskEncryptionConfiguration)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDiskEncryptionConfigurationByID(resourceIDInt, diskEncryptionConfig)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateDiskEncryptionConfigurationByName(resourceName, diskEncryptionConfig)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the Site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDiskEncryptionConfigurationsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDiskEncryptionConfiguration, resourceID)
	}

	return nil
}

// ResourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func ResourceJamfProDiskEncryptionConfigurationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDiskEncryptionConfiguration, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteDiskEncryptionConfigurationByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
