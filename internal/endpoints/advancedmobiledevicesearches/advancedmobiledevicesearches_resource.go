// advancedmobiledevicesearches_resource.go
package advancedmobiledevicesearches

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProAdvancedMobileDeviceSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Mobile Device Search: %v", err))
	}

	// Retry the API call to create the site in Jamf Pro
	var creationResponse *jamfpro.ResponseAdvancedMobileDeviceSearchCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAdvancedMobileDeviceSearch(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Advanced Mobile Device Search '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchRead is responsible for reading the current state of a Jamf Pro mobile device Search from the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var resource *jamfpro.ResourceAdvancedMobileDeviceSearch

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetAdvancedMobileDeviceSearchByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro Advanced Mobile Device Search with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Advanced Mobile Device Search with ID '%s' after retries: %v", resourceID, err))
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("view_as", resource.ViewAs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort1", resource.Sort1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort2", resource.Sort2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort3", resource.Sort3); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "criteria" field
	criteriaList := make([]interface{}, len(resource.Criteria.Criterion))
	for i, crit := range resource.Criteria.Criterion {
		criteriaMap := map[string]interface{}{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        crit.AndOr,
			"search_type":   crit.SearchType,
			"value":         crit.Value,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
		criteriaList[i] = criteriaMap
	}
	if err := d.Set("criteria", criteriaList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "display_fields" field
	if len(resource.DisplayFields) == 0 || len(resource.DisplayFields[0].DisplayField) == 0 {
		if err := d.Set("display_fields", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		displayFieldsList := make([]map[string]interface{}, len(resource.DisplayFields[0].DisplayField))
		for i, displayField := range resource.DisplayFields[0].DisplayField {
			displayFieldMap := map[string]interface{}{
				"name": displayField.Name,
			}
			displayFieldsList[i] = displayFieldMap
		}
		if err := d.Set("display_fields", displayFieldsList); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Handle "site" field
	site := map[string]interface{}{
		"id":   resource.Site.ID,
		"name": resource.Site.Name,
	}
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchUpdate is responsible for updating an existing Jamf Pro mobile device Search on the remote system.
func ResourceJamfProAdvancedMobileDeviceSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProAdvancedMobileDeviceSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Mobile Device Search for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateAdvancedMobileDeviceSearchByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Advanced User Search '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAdvancedMobileDeviceSearchRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAdvancedMobileDeviceSearchDelete is responsible for deleting a Jamf Pro AdvancedMobileDeviceSearch.
func ResourceJamfProAdvancedMobileDeviceSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteAdvancedMobileDeviceSearchByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteAdvancedMobileDeviceSearchByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Advanced Mobile Device Search '%s' (ID: %s) after retries: %v", d.Get("extension").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
