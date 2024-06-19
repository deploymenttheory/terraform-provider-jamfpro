// networksegments_resource.go
package networksegments

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

// ResourceJamfProNetworkSegments defines the schema and CRUD operations for managing Jamf Pro NetworkSegments in Terraform.
func ResourceJamfProNetworkSegments() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProNetworkSegmentsCreate,
		ReadContext:   ResourceJamfProNetworkSegmentsRead,
		UpdateContext: ResourceJamfProNetworkSegmentsUpdate,
		DeleteContext: ResourceJamfProNetworkSegmentsDelete,
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
				Description: "The unique identifier of the network segment.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network segment.",
			},
			"starting_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The starting IP address of the network segment.",
			},
			"ending_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ending IP address of the network segment.",
			},
			"distribution_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution server associated with the network segment.",
			},
			"distribution_point": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution point associated with the network segment.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL associated with the network segment.",
			},
			"swu_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The software update server associated with the network segment.",
			},
			"building": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The building associated with the network segment.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The department associated with the network segment.",
			},
			"override_buildings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if building assignments are overridden for this network segment.",
			},
			"override_departments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if department assignments are overridden for this network segment.",
			},
		},
	}
}

// ResourceJamfProNetworkSegmentsCreate is responsible for creating a new Jamf Network segment in the remote system.
// The function:
// 1. Constructs the Network Segment data using the provided Terraform configuration.
// 2. Calls the API to create the Network Segment in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created Network Segment.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProNetworkSegmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProNetworkSegment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Network Segment: %v", err))
	}

	var creationResponse *jamfpro.ResponseNetworkSegmentCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateNetworkSegment(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Network Segment '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetNetworkSegmentByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Network Segment", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProNetworkSegmentsRead(ctx, d, meta)...)
}

// ResourceJamfProNetworkSegmentsRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProNetworkSegmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetNetworkSegmentByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProNetworkSegmentsUpdate is responsible for updating an existing Jamf Pro Network Segment on the remote system.
func ResourceJamfProNetworkSegmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProNetworkSegment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Network Segment for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateNetworkSegmentByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Network Segment '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, ResourceJamfProNetworkSegmentsRead(ctx, d, meta)...)
}

// ResourceJamfProNetworkSegmentsDeleteis responsible for deleting a Jamf Pro network segment.
func ResourceJamfProNetworkSegmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteNetworkSegmentByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteNetworkSegmentByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Network Segment '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
