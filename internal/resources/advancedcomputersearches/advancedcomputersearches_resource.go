// advancedcomputersearches_resource.go
package advancedcomputersearches

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the advanced computer search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the advanced computer search",
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
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the criteria list",
						},
						"criterion": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of the criteria",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Priority of the criteria",
									},
									"and_or": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Logical operator (AND or OR)",
									},
									"search_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Search operator",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value for the criteria",
									},
									"opening_paren": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates if an opening parenthesis is used",
									},
									"closing_paren": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates if a closing parenthesis is used",
									},
								},
							},
						},
					},
				},
			},
			"display_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the display fields list",
						},
						"display_field": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of the display field",
									},
								},
							},
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

// constructJamfProAdvancedComputerSearch constructs a ResourceAdvancedComputerSearch object from the provided schema data.
func constructJamfProAdvancedComputerSearch(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAdvancedComputerSearch, error) {
	search := &jamfpro.ResourceAdvancedComputerSearch{}

	// Utilize type assertion helper functions for direct field extraction
	search.Name = util.GetStringFromInterface(d.Get("name"))
	search.ViewAs = util.GetStringFromInterface(d.Get("view_as"))
	search.Sort1 = util.GetStringFromInterface(d.Get("sort1"))
	search.Sort2 = util.GetStringFromInterface(d.Get("sort2"))
	search.Sort3 = util.GetStringFromInterface(d.Get("sort3"))

	// Handle nested "criteria" field
	if criteriaList, ok := d.GetOk("criteria"); ok {
		var criteria []jamfpro.SharedSubsetCriteria
		for _, crit := range criteriaList.([]interface{}) {
			criterionMap := util.ConvertToMapFromInterface(crit)
			newCriterion := jamfpro.SharedSubsetCriteria{
				Name:         util.GetStringFromMap(criterionMap, "name"),
				Priority:     util.GetIntFromMap(criterionMap, "priority"),
				AndOr:        util.GetStringFromMap(criterionMap, "and_or"),
				SearchType:   util.GetStringFromMap(criterionMap, "search_type"),
				Value:        util.GetStringFromMap(criterionMap, "value"),
				OpeningParen: util.GetBoolFromMap(criterionMap, "opening_paren"),
				ClosingParen: util.GetBoolFromMap(criterionMap, "closing_paren"),
			}
			criteria = append(criteria, newCriterion)
		}
		search.Criteria = jamfpro.SharedContainerCriteria{
			Size:      len(criteria),
			Criterion: criteria,
		}
	}

	// Handle nested "display_fields" field
	if displayFieldsList, ok := d.GetOk("display_fields"); ok {
		var displayFields []jamfpro.SharedAdvancedSearchSubsetDisplayField
		for _, field := range displayFieldsList.([]interface{}) {
			displayFieldMap := util.ConvertToMapFromInterface(field)
			newDisplayField := jamfpro.SharedAdvancedSearchSubsetDisplayField{
				Name: util.GetStringFromMap(displayFieldMap, "name"),
			}
			displayFields = append(displayFields, newDisplayField)
		}
		search.DisplayFields = displayFields
	}

	// Handle nested "site" field
	if siteList, ok := d.GetOk("site"); ok && len(siteList.([]interface{})) > 0 {
		siteData := util.ConvertToMapFromInterface(siteList.([]interface{})[0])
		search.Site = jamfpro.SharedResourceSite{
			ID:   util.GetIntFromMap(siteData, "id"),
			Name: util.GetStringFromMap(siteData, "name"),
		}
	}

	// Logging the constructed object in debug mode
	tflog.Debug(ctx, fmt.Sprintf("Constructed AdvancedComputerSearch Object: %+v", search))

	// Log the successful construction of the Jamf Pro AdvancedComputerSearch
	log.Printf("[INFO] Successfully constructed AdvancedComputerSearch with name: %s", search.Name)

	return search, nil
}

// Helper function to generate diagnostics based on the error type.
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProAdvancedComputerSearchCreate is responsible for creating a new Jamf Pro Advanced Computer Search in the remote system.
func ResourceJamfProAdvancedComputerSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the advanced computer search
		search, err := constructJamfProAdvancedComputerSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the advanced computer search for terraform create: %w", err))
		}

		// Log the details of the search that is about to be created
		log.Printf("[INFO] Attempting to create AdvancedComputerSearch with name: %s", search.Name)

		// Directly call the API to create the resource
		response, err := conn.CreateAdvancedComputerSearch(search)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating AdvancedComputerSearch with name: %s. Error: %s", search.Name, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created AdvancedComputerSearch with ID: %d and name: %s", response.ID, search.Name)

		// Set the ID of the created resource in the Terraform state
		d.SetId(strconv.Itoa(response.ID))

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAdvancedComputerSearchRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProAdvancedComputerSearchRead is responsible for reading the current state of a Jamf Pro Advanced Computer Search from the remote system.
func ResourceJamfProAdvancedComputerSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var search *jamfpro.ResourceAdvancedComputerSearch

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Try fetching the advanced computer search using the ID
		var apiErr error
		search, apiErr = conn.GetAdvancedComputerSearchByID(searchID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			searchName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			search, apiErr = conn.GetAdvancedComputerSearchByName(searchName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Set attributes in the Terraform state
	if err := d.Set("name", search.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("view_as", search.ViewAs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort1", search.Sort1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort2", search.Sort2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sort3", search.Sort3); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "criteria" field
	criteriaList := make([]interface{}, len(search.Criteria.Criterion))
	for i, crit := range search.Criteria.Criterion {
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
	displayFieldsList := make([]interface{}, len(search.DisplayFields))
	for i, displayField := range search.DisplayFields {
		displayFieldMap := map[string]interface{}{
			"name": displayField.Name,
		}
		displayFieldsList[i] = displayFieldMap
	}
	if err := d.Set("display_fields", displayFieldsList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle "site" field
	site := map[string]interface{}{
		"id":   search.Site.ID,
		"name": search.Site.Name,
	}
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProAdvancedComputerSearchUpdate is responsible for updating an existing Jamf Pro Advanced Computer Search on the remote system.
func ResourceJamfProAdvancedComputerSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated advanced computer search
		search, err := constructJamfProAdvancedComputerSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the advanced computer search for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateAdvancedComputerSearchByID(searchID, search)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			searchName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
			_, apiErr = conn.UpdateAdvancedComputerSearchByName(searchName, search)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAdvancedComputerSearchRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProAdvancedComputerSearchDelete is responsible for deleting a Jamf Pro AdvancedComputerSearch.
func ResourceJamfProAdvancedComputerSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		groupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse group ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteAdvancedComputerSearchByID(groupID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
			apiErr = conn.DeleteAdvancedComputerSearchByName(groupName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
