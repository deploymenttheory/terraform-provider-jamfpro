package computer_inventory_collection_settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errConstructSettings         = errors.New("failed to construct Jamf Pro Computer Inventory Collection Settings")
	errApplySettings             = errors.New("failed to apply Jamf Pro Computer Inventory Collection Settings")
	errConstructCustomPaths      = errors.New("failed to construct custom paths")
	errCreateCustomPath          = errors.New("failed to create custom path")
	errDeleteCustomPath          = errors.New("failed to delete custom path")
	errDeleteCustomPathOnDestroy = errors.New("failed to delete custom path during resource destruction")
)

// create is responsible for initializing the Jamf Pro Computer Inventory Collection Settings in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	settings, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %w", errConstructSettings, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateComputerInventoryCollectionSettings(settings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %w", errApplySettings, err))
	}

	customPaths, err := constructCustomPaths(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %w", errConstructCustomPaths, err))
	}

	for _, customPath := range customPaths {
		pathToCreate := customPath
		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			_, apiErr := client.CreateComputerInventoryCollectionSettingsCustomPath(&pathToCreate)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s: %w", errCreateCustomPath, customPath.Path, err))
		}
	}

	d.SetId("jamfpro_computer_inventory_collection_settings_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Computer Inventory Collection Settings.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	var err error

	d.SetId("jamfpro_computer_inventory_collection_settings_singleton")
	var response *jamfpro.ResourceComputerInventoryCollectionSettings
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetComputerInventoryCollectionSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the Jamf Pro Computer Inventory Collection Settings.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	if d.HasChange("computer_inventory_collection_preferences") {
		settings, err := construct(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %w", errConstructSettings, err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, apiErr := client.UpdateComputerInventoryCollectionSettings(settings)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %w", errApplySettings, err))
		}
	}

	pathsToAdd, pathIDsToRemove := constructPathUpdates(d)

	for _, idToDelete := range pathIDsToRemove {
		id := idToDelete // Create a copy for the closure
		log.Printf("[DEBUG] Deleting custom path with ID: %s", id)

		// Skip built-in paths which use ID "-1"
		if id == "-1" {
			log.Printf("[DEBUG] Skipping deletion of built-in custom path with ID: %s", id)
			continue
		}

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			apiErr := client.DeleteComputerInventoryCollectionSettingsCustomPathByID(id)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s: %w", errDeleteCustomPath, id, err))
		}
	}

	// Add new paths
	for _, pathToAdd := range pathsToAdd {
		addPath := pathToAdd // Create a copy for the closure
		pathJSON, _ := json.MarshalIndent(addPath, "", "  ")
		log.Printf("[DEBUG] Creating custom path:\n%s", string(pathJSON))

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, apiErr := client.CreateComputerInventoryCollectionSettingsCustomPath(&addPath)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s: %w", errCreateCustomPath, addPath.Path, err))
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Computer Inventory Collection Settings.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state after cleaning up custom paths.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	// Clean up all custom paths
	for _, pt := range pathTypes {
		if v, ok := d.GetOk(pt.key); ok {
			pathSet := v.(*schema.Set)
			for _, p := range pathSet.List() {
				pathMap := p.(map[string]interface{})
				path := pathMap["path"].(string)
				id := pathMap["id"].(string)

				log.Printf("[DEBUG] Deleting custom path during destruction:\n  Path: %s\n  ID: %s", path, id)

				// Skip built-in paths which use ID "-1"
				if id == "-1" {
					log.Printf("[DEBUG] Skipping deletion of built-in custom path during destruction (ID: %s, Path: %s)", id, path)
					continue
				}

				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					apiErr := client.DeleteComputerInventoryCollectionSettingsCustomPathByID(id)
					if apiErr != nil {
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					return diag.FromErr(fmt.Errorf("%w: (ID: %s, Path: %s): %w",
						errDeleteCustomPathOnDestroy, id, path, err))
				}

				log.Printf("[DEBUG] Successfully deleted custom path:\n  Path: %s\n  ID: %s", path, id)
			}
		}
	}

	log.Printf("[DEBUG] All custom paths deleted, removing resource from state")

	// Remove from state (base settings aren't actually deletable)
	d.SetId("")
	return nil
}
