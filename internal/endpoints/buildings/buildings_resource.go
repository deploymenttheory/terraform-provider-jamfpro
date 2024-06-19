// buildings_resource.go
package buildings

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProBuildings defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProBuildings() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProBuildingCreate,
		ReadContext:   ResourceJamfProBuildingRead,
		UpdateContext: ResourceJamfProBuildingUpdate,
		DeleteContext: ResourceJamfProBuildingDelete,
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
				Description: "The unique identifier of the building.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the building.",
			},
			"street_address1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first line of the street address of the building.",
			},
			"street_address2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The second line of the street address of the building.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The city in which the building is located.",
			},
			"state_province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The state or province in which the building is located.",
			},
			"zip_postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ZIP or postal code of the building.",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The country in which the building is located.",
			},
		},
	}
}

// ResourceJamfProBuildingCreate is responsible for creating a new Building in the remote system.
// The function:
// 1. Constructs the building data using the provided Terraform configuration.
// 2. Calls the API to create the building in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created building.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProBuildingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProBuilding(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Building: %v", err))
	}

	var creationResponse *jamfpro.ResponseBuildingCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateBuilding(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Building '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	return client.GetBuildingByID(id.(string))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Building", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProBuildingRead(ctx, d, meta)...)
}

// ResourceJamfProBuildingRead is responsible for reading the current state of a Building Resource from the remote system.
// The function:
// 1. Fetches the building's current state using its ID. If it fails, then obtain the building's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the building being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProBuildingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := client.GetBuildingByID(resourceID)
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProBuildingUpdate is responsible for updating an existing Building on the remote system.
func ResourceJamfProBuildingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProBuilding(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Building for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateBuildingByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Building '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	return append(diags, ResourceJamfProBuildingRead(ctx, d, meta)...)
}

// ResourceJamfProBuildingDelete is responsible for deleting a Building.
func ResourceJamfProBuildingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteBuildingByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteBuildingByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Building '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
