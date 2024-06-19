// advancedmobiledevicesearches_resource.go
package advancedmobiledevicesearches

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

// ResourceJamfProAdvancedMobileDeviceSearches defines the schema for managing mobile device Searches in Terraform.
func ResourceJamfProAdvancedMobileDeviceSearches() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAdvancedMobileDeviceSearchCreate,
		ReadContext:   ResourceJamfProAdvancedMobileDeviceSearchRead,
		UpdateContext: ResourceJamfProAdvancedMobileDeviceSearchUpdate,
		DeleteContext: ResourceJamfProAdvancedMobileDeviceSearchDelete,
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
				Description: "The unique identifier of the advanced mobile device search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the advanced mobile device search",
			},
			"view_as": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "View type of the mobile device search",
			},
			"sort1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First sorting criteria for the mobile device search",
			},
			"sort2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second sorting criteria for the mobile device search",
			},
			"sort3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Third sorting criteria for the mobile device search",
			},
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"and_or": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"search_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"opening_paren": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"closing_paren": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"display_fields": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "display field in the mobile device search",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display field item in the mobile device search",
						},
					},
				},
			},
			"site": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the site",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the site",
						},
					},
				},
			},
		},
	}
}

// ResourceJamfProAdvancedMobileDeviceSearchCreate is responsible for creating a new Jamf Pro mobile device Search in the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProAdvancedMobileDeviceSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Mobile Device Search: %v", err))
	}

	var creationResponse *jamfpro.ResponseAdvancedMobileDeviceSearchCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAdvancedMobileDeviceSearch(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Advanced Mobile Device Search '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetAdvancedComputerSearchByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Advanced Mobile Device Search", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)...)
}

// ResourceJamfProAdvancedMobileDeviceSearchRead is responsible for reading the current state of a Jamf Pro mobile device Search from the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetAdvancedMobileDeviceSearchByID(resourceIDInt)
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProAdvancedMobileDeviceSearchUpdate is responsible for updating an existing Jamf Pro mobile device Search on the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProAdvancedMobileDeviceSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Mobile Device Search for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateAdvancedMobileDeviceSearchByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Advanced Mobile Device Search '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	return append(diags, ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)...)
}

// ResourceJamfProAdvancedMobileDeviceSearchDelete is responsible for deleting a Jamf Pro AdvancedMobileDeviceSearch.
func ResourceJamfProAdvancedMobileDeviceSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteAdvancedMobileDeviceSearchByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteAdvancedMobileDeviceSearchByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Advanced Mobile Device Search '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
