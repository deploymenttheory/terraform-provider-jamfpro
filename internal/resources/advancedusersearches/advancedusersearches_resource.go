// advancedusersearches_resource.go
package advancedusersearches

import (
	"context"
	"encoding/xml"
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

// ResourceJamfProAdvancedUserSearches defines the schema for managing advanced user Searches in Terraform.
func ResourceJamfProAdvancedUserSearches() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAdvancedUserSearchCreate,
		ReadContext:   ResourceJamfProAdvancedUserSearchRead,
		UpdateContext: ResourceJamfProAdvancedUserSearchUpdate,
		DeleteContext: ResourceJamfProAdvancedUserSearchDelete,
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
				Description: "The unique identifier of the advanced user search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the advanced mobile device search",
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
				Description: "display field in the advanced user search",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display field item in the advanced user search",
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

// constructJamfProAdvancedUserSearch constructs an advanced user search object for create and update oeprations
func constructJamfProAdvancedUserSearch(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAdvancedUserSearch, error) {
	search := &jamfpro.ResourceAdvancedUserSearch{
		Name: util.GetStringFromInterface(d.Get("name")),
	}

	if v, ok := d.GetOk("criteria"); ok {
		criteriaList := v.([]interface{})
		var criteria []jamfpro.SharedSubsetCriteria
		for _, crit := range criteriaList {
			criterionMap, ok := crit.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to parse criterion: %+v", crit)
			}

			newCriterion := jamfpro.SharedSubsetCriteria{
				Name:         util.GetStringFromMap(criterionMap, "name"),
				Priority:     util.GetIntFromMap(criterionMap, "priority"),
				AndOr:        util.GetStringFromMap(criterionMap, "and_or"),
				SearchType:   util.GetStringFromMap(criterionMap, "search_type"),
				Value:        util.GetStringFromMap(criterionMap, "value"),
				OpeningParen: util.GetBoolFromMap(criterionMap, "opening_paren"),
				ClosingParen: util.GetBoolFromMap(criterionMap, "closing_paren"),
			}
			tflog.Debug(ctx, fmt.Sprintf("Processing criterion: %+v", newCriterion))
			criteria = append(criteria, newCriterion)
		}
		search.Criteria.Criterion = criteria
	}

	// Initialize displayFields as an empty slice
	displayFields := []jamfpro.SharedAdvancedSearchSubsetDisplayField{}

	// Check if display_fields are provided in the configuration
	if v, ok := d.GetOk("display_fields"); ok {
		// Since display_fields is now a TypeSet, use the *schema.Set methods to access the data
		for _, field := range v.(*schema.Set).List() {
			displayFieldMap, ok := field.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to parse display field: %+v", field)
			}

			displayField := jamfpro.SharedAdvancedSearchSubsetDisplayField{
				Name: util.GetStringFromMap(displayFieldMap, "name"),
			}
			displayFields = append(displayFields, displayField)
		}
	}

	search.DisplayFields = []jamfpro.SharedAdvancedSearchContainerDisplayField{{DisplayField: displayFields}}

	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData, ok := v.([]interface{})[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse site data")
		}

		search.Site = jamfpro.SharedResourceSite{
			ID:   util.GetIntFromMap(siteData, "id"),
			Name: util.GetStringFromMap(siteData, "name"),
		}
	}

	// Marshal the search object into XML for logging
	xmlData, err := xml.MarshalIndent(search, "", "  ")
	if err != nil {
		// Handle the error if XML marshaling fails
		log.Printf("[ERROR] Error marshaling AdvancedUserSearch object to XML: %s", err)
		return nil, fmt.Errorf("error marshaling AdvancedUserSearch object to XML: %v", err)
	}

	// Log the XML formatted search object
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Constructed AdvancedUserSearch Object:\n%s", string(xmlData)))

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Successfully constructed AdvancedUserSearch with name: %s", search.Name))

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

// ResourceJamfProAdvancedUserSearchCreate is responsible for creating a new Jamf Pro advanced user Search in the remote system.
func ResourceJamfProAdvancedUserSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the advanced user search
		search, err := constructJamfProAdvancedUserSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the advanced user search for terraform create: %w", err))
		}

		// Log the details of the search that is about to be created
		log.Printf("[INFO] Attempting to create AdvancedUserSearch with name: %s", search.Name)

		// Directly call the API to create the resource
		response, err := conn.CreateAdvancedUserSearch(search)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating AdvancedUserSearch with name: %s. Error: %s", search.Name, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created AdvancedUserSearch with ID: %d and name: %s", response.ID, search.Name)

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
		readDiags := ResourceJamfProAdvancedUserSearchRead(ctx, d, meta)
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

// ResourceJamfProAdvancedUserSearchRead is responsible for reading the current state of a Jamf Pro advanced user Search from the remote system.
func ResourceJamfProAdvancedUserSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var search *jamfpro.ResourceAdvancedUserSearch

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Try fetching the advanced user search using the ID
		var apiErr error
		search, apiErr = conn.GetAdvancedUserSearchByID(searchID)
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

			search, apiErr = conn.GetAdvancedUserSearchByName(searchName)
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
	if len(search.DisplayFields) == 0 || len(search.DisplayFields[0].DisplayField) == 0 {
		if err := d.Set("display_fields", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		displayFieldsList := make([]map[string]interface{}, len(search.DisplayFields[0].DisplayField))
		for i, displayField := range search.DisplayFields[0].DisplayField {
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
		"id":   search.Site.ID,
		"name": search.Site.Name,
	}
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProAdvancedUserSearchUpdate is responsible for updating an existing Jamf Pro advanced user Search on the remote system.
func ResourceJamfProAdvancedUserSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		// Construct the updated advanced user search
		search, err := constructJamfProAdvancedUserSearch(ctx, d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the advanced user search for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		searchID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse search ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateAdvancedUserSearchByID(searchID, search)
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
			_, apiErr = conn.UpdateAdvancedUserSearchByName(searchName, search)
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
		readDiags := ResourceJamfProAdvancedUserSearchRead(ctx, d, meta)
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

// ResourceJamfProAdvancedUserSearchDelete is responsible for deleting a Jamf Pro AdvancedUserSearch.
func ResourceJamfProAdvancedUserSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteAdvancedUserSearchByID(groupID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
			apiErr = conn.DeleteAdvancedUserSearchByName(groupName)
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
