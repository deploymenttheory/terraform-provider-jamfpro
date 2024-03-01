// networksegments_resource.go
package networksegments

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

// ResourceJamfProNetworkSegments defines the schema and CRUD operations for managing Jamf Pro NetworkSegments in Terraform.
func ResourceJamfProNetworkSegments() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProNetworkSegmentsCreate,
		ReadContext:   ResourceJamfProNetworkSegmentsRead,
		UpdateContext: ResourceJamfProNetworkSegmentsUpdate,
		DeleteContext: ResourceJamfProNetworkSegmentsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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
// 1. Constructs the printer data using the provided Terraform configuration.
// 2. Calls the API to create the printer in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created printer.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProNetworkSegmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProNetworkSegment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Site: %v", err))
	}

	// Retry the API call to create the site in Jamf Pro
	var creationResponse *jamfpro.ResponseNetworkSegmentCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateNetworkSegment(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Network Segment '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProNetworkSegmentsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProNetworkSegmentsRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProNetworkSegmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var resource *jamfpro.ResourceNetworkSegment

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetNetworkSegmentByID(resourceIDInt)
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
				Detail:   fmt.Sprintf("Jamf Pro Site with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Site with ID '%s' after retries: %v", resourceID, err))
	}

	// Update Terraform state with the resource information
	if resource != nil {
		if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("starting_address", resource.StartingAddress); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ending_address", resource.EndingAddress); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("distribution_server", resource.DistributionServer); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("distribution_point", resource.DistributionPoint); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("url", resource.URL); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("swu_server", resource.SWUServer); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("building", resource.Building); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("department", resource.Department); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("override_buildings", resource.OverrideBuildings); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("override_departments", resource.OverrideDepartments); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// ResourceJamfProNetworkSegmentsUpdate is responsible for updating an existing Jamf Pro Network Segment on the remote system.
func ResourceJamfProNetworkSegmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProNetworkSegment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Network Segment for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateNetworkSegmentByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Network Segment '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProNetworkSegmentsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProNetworkSegmentsDeleteis responsible for deleting a Jamf Pro network segment.
func ResourceJamfProNetworkSegmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteNetworkSegmentByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteNetworkSegmentByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Network Segment '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
