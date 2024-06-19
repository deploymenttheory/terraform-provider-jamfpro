// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

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

// ResourceJamfProDiskEncryptionConfigurationsCreate is responsible for creating a new Jamf Pro Disk Encryption Configuration in the remote system.
// The function:
// 1. Constructs the disk encryption configuration data using the provided Terraform configuration.
// 2. Calls the API to create the disk encryption configuration in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created disk encryption configuration.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDiskEncryptionConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProDiskEncryptionConfiguration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Disk Encryption Configuration: %v", err))
	}

	var creationResponse *jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDiskEncryptionConfiguration(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Disk Encryption Configuration '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetDiskEncryptionConfigurationByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Disk Encryption Configuration", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta)...)
}

// ResourceJamfProDiskEncryptionConfigurationsRead is responsible for reading the current state of a Jamf Pro Disk Encryption Configuration resource from Jamf Pro and updating the Terraform state with the retrieved data.
func ResourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetDiskEncryptionConfigurationByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProDiskEncryptionConfigurationsUpdate is responsible for updating an existing Jamf Pro Disk Encryption Configuration on the remote system.
func ResourceJamfProDiskEncryptionConfigurationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProDiskEncryptionConfiguration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Disk Encryption Configuration for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDiskEncryptionConfigurationByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Disk Encryption Configuration '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, ResourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta)...)
}

// ResourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func ResourceJamfProDiskEncryptionConfigurationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteDiskEncryptionConfigurationByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Disk Encryption Configuration '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
