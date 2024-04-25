// advancedcomputersearches_resource.go
package advancedcomputersearches

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAdvancedComputerSearches defines the schema for managing Advanced Computer Searches in Terraform.
func ResourceJamfProAdvancedComputerSearches() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAdvancedComputerSearchCreate,
		ReadContext:   ResourceJamfProAdvancedComputerSearchRead,
		UpdateContext: ResourceJamfProAdvancedComputerSearchUpdate,
		DeleteContext: ResourceJamfProAdvancedComputerSearchDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the advanced computer search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the advanced computer search",
			},
			"view_as": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "View type of the advanced computer search",
			},
			"sort1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First sorting criteria for the advanced computer search",
			},
			"sort2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second sorting criteria for the advanced computer search",
			},
			"sort3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Third sorting criteria for the advanced computer search",
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
				Description: "display field in the advanced computer search",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display field item in the advanced computer search",
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

// ResourceJamfProAdvancedComputerSearchCreate is responsible for creating a new Jamf Pro Advanced Computer Search in the remote system.
func ResourceJamfProAdvancedComputerSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProAdvancedComputerSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Computer Search: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseAdvancedComputerSearchCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAdvancedComputerSearch(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Advanced Computer Search '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetAdvancedComputerSearchByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Advanced Computer Search", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAdvancedComputerSearchRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAdvancedComputerSearchRead is responsible for reading the current state of a Jamf Pro Advanced Computer Search from the remote system.
func ResourceJamfProAdvancedComputerSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := conn.GetAdvancedComputerSearchByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProAdvancedComputerSearchUpdate is responsible for updating an existing Jamf Pro Advanced Computer Search on the remote system.
func ResourceJamfProAdvancedComputerSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProAdvancedComputerSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Computer Search for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateAdvancedComputerSearchByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Advanced Computer Search '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAdvancedComputerSearchRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAdvancedComputerSearchDelete is responsible for deleting a Jamf Pro AdvancedComputerSearch.
func ResourceJamfProAdvancedComputerSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resourceName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteAdvancedComputerSearchByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			apiErrByName := conn.DeleteAdvancedComputerSearchByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Advanced Computer Search '%s' (ID: %s) after retries: %v", resourceName, resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
