// dockitems_resource.go
package dockitems

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProDockItems defines the schema and CRUD operations for managing Jamf Pro Dock Items in Terraform.
func ResourceJamfProDockItems() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProDockItemsCreate,
		ReadContext:   ResourceJamfProDockItemsRead,
		UpdateContext: ResourceJamfProDockItemsUpdate,
		DeleteContext: ResourceJamfProDockItemsDelete,
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
				Description: "The unique identifier of the dock item.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the dock item.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the dock item (App/File/Folder).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("expected a string for %q but got a different type", key))
						return
					}
					validTypes := map[string]bool{
						"App":    true,
						"File":   true,
						"Folder": true,
					}
					if !validTypes[v] {
						errs = append(errs, fmt.Errorf("%q must be one of 'App', 'File', or 'Folder', got: %s", key, v))
					}
					return
				},
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path of the dock item.",
			},
			"contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the dock item.",
			},
		},
	}
}

// ResourceJamfProDockItemsCreate is responsible for creating a new Jamf Pro Dock Item in the remote system.
// The function:
// 1. Constructs the dock item data using the provided Terraform configuration.
// 2. Calls the API to create the dock item in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created dock item.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDockItemsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProDockItem(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Dock Item: %v", err))
	}

	var creationResponse *jamfpro.ResourceDockItem
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDockItem(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Dock Item '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetDockItemByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Dock Item", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProDockItemsRead(ctx, d, meta)...)
}

// ResourceJamfProDockItemsRead is responsible for reading the current state of a Jamf Pro Dock Item Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProDockItemsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetDockItemByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProDockItemsUpdate is responsible for updating a Jamf Pro dock item.
func ResourceJamfProDockItemsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProDockItem(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Dock Item for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDockItemByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Dock Item '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, ResourceJamfProDockItemsRead(ctx, d, meta)...)
}

// ResourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func ResourceJamfProDockItemsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDockItemByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteDockItemByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Dock Item '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
