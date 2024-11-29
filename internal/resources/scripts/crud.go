package scripts

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Create requires a mutex need to lock Create requests during parallel runs

// create is responsible for creating a new Jamf Pro Script in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// mu.Lock()
	// defer mu.Unlock()
	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateScript,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetScriptByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating an existing Jamf Pro Department on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateScriptByID,
		readNoCleanup,
	)
}

// delete is responsible for deleting a Jamf Pro Department.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// mu.Lock()
	// defer mu.Unlock()
	log.Print(("INSIDE"))

	log.Println("LOGHERE - DELETING")

	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		var resp *http.Response

		resp, apiErr := meta.(*jamfpro.Client).DeleteScriptByID(resourceID)

		if apiErr != nil {
			log.Printf("ERROR FOUND - %v", apiErr)
			return retry.RetryableError(apiErr)
		}

		log.Printf("NO ERROR FOUND - %v, %v", resp, apiErr)
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro resource '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
