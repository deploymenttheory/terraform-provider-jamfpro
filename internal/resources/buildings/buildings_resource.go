// buildings_resource.go
package buildings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProBuilding defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProBuilding() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProBuildingCreate,
		ReadContext:   ResourceJamfProBuildingRead,
		UpdateContext: ResourceJamfProBuildingUpdate,
		DeleteContext: ResourceJamfProBuildingDelete,
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

// constructBuilding constructs a Building object from the provided schema data.
func constructBuilding(d *schema.ResourceData) *jamfpro.ResponseBuilding {
	building := &jamfpro.ResponseBuilding{}

	if v, ok := d.GetOk("name"); ok {
		building.Name = v.(string)
	}
	if v, ok := d.GetOk("street_address1"); ok {
		building.StreetAddress1 = v.(string)
	}
	if v, ok := d.GetOk("street_address2"); ok {
		building.StreetAddress2 = v.(string)
	}
	if v, ok := d.GetOk("city"); ok {
		building.City = v.(string)
	}
	if v, ok := d.GetOk("state_province"); ok {
		building.StateProvince = v.(string)
	}
	if v, ok := d.GetOk("zip_postal_code"); ok {
		building.ZipPostalCode = v.(string)
	}
	if v, ok := d.GetOk("country"); ok {
		building.Country = v.(string)
	}

	return building
}

// Helper function to generate diagnostics based on the error type
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

// ResourceJamfProBuildingCreate is responsible for creating a new Building in the remote system.
// The function:
// 1. Constructs the building data using the provided Terraform configuration.
// 2. Calls the API to create the building in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created building.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProBuildingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the create operation
	var createdBuilding *jamfpro.ResponseBuildingCreate
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the building
		building := constructBuilding(d)

		// Check if the building is nil
		if building == nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the building"))
		}

		// Directly call the API to create the resource
		createdBuilding, err = conn.CreateBuilding(building)
		if err != nil {
			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(createdBuilding.ID)

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProBuildingRead(ctx, d, meta)
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

// ResourceJamfProBuildingRead is responsible for reading the current state of a Building Resource from the remote system.
// The function:
// 1. Fetches the building's current state using its ID. If it fails, then obtain the building's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the building being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProBuildingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	var building *jamfpro.ResponseBuilding

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// The ID in Terraform state is already a string, so we use it directly for the API request
		buildingID := d.Id()

		// Try fetching the building using the ID
		var apiErr error
		building, apiErr = conn.GetBuildingByID(buildingID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			buildingName := d.Get("name").(string)
			building, apiErr = conn.GetBuildingByNameByID(buildingName)
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

	// Safely set all attributes in the Terraform state
	if err := d.Set("name", building.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("street_Address1", building.StreetAddress1); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("street_Address2", building.StreetAddress2); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("city", building.City); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("state_province", building.StateProvince); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("zip_postal_code", building.ZipPostalCode); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("country", building.Country); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProBuildingUpdate is responsible for updating an existing Building on the remote system.
func ResourceJamfProBuildingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the building
		building := constructBuilding(d)

		// The ID in Terraform state is already a string, so we use it directly for the API request
		buildingID := d.Id()

		// Directly call the API to update the resource by ID
		_, apiErr := conn.UpdateBuildingByID(buildingID, building)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			buildingName := d.Get("name").(string)
			_, apiErr = conn.UpdateBuildingByNameByID(buildingName, building)
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
		readDiags := ResourceJamfProBuildingRead(ctx, d, meta)
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

// ResourceJamfProBuildingDelete is responsible for deleting a Building.
func ResourceJamfProBuildingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the DELETE operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// The ID in Terraform state is already a string, so we use it directly for the API request
		buildingID := d.Id()

		// Directly call the API to DELETE the resource by ID
		apiErr := conn.DeleteBuildingByID(buildingID)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			buildingName := d.Get("name").(string)
			apiErr = conn.DeleteBuildingByNameByID(buildingName)
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
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
