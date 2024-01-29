// sites_resource.go
package sites

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProSite defines the schema and CRUD operations for managing Jamf Pro Sites in Terraform.
func ResourceJamfProSites() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProSitesCreate,
		ReadContext:   ResourceJamfProSitesRead,
		UpdateContext: ResourceJamfProSitesUpdate,
		DeleteContext: ResourceJamfProSitesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the site.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro site.",
			},
		},
	}
}

// constructJamfProSite constructs a SharedResourceSite object from the provided schema data.
func constructJamfProSite(ctx context.Context, d *schema.ResourceData) (*jamfpro.SharedResourceSite, error) {
	site := &jamfpro.SharedResourceSite{
		Name: util.GetStringFromInterface(d.Get("name")),
	}

	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	// Serialize and pretty-print the site object as XML for logging
	siteXML, err := xml.MarshalIndent(site, "", "  ")
	if err != nil {
		logging.Error(subCtx, logging.SubsystemConstruct, "Failed to marshal site to XML", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// Log the pretty-printed XML of the constructed site
	logging.Debug(subCtx, logging.SubsystemConstruct, "Constructed Site XML", map[string]interface{}{"xml": string(siteXML)})

	return site, nil
}

// ResourceJamfProSitesCreate is responsible for creating a new Jamf Pro Site in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProSitesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize diagnostics for collecting any issues to report back to Terraform
	var diags diag.Diagnostics

	// Initialize tflog
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)

	// Initialize variables
	var createdSite *jamfpro.SharedResourceSite

	// construct the resource object
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		site, err := constructJamfProSite(subCtx, d)
		if err != nil {
			logging.Error(subCtx, logging.SubsystemCreate, "Failed to construct site", map[string]interface{}{
				"name":  d.Get("name"),
				"error": err.Error(),
			})
			return retry.NonRetryableError(err)
		}

		createdSite, err = conn.CreateSite(site)
		if err != nil {
			if apiErr, ok := err.(*http_client.APIError); ok {
				logging.Error(subCtx, logging.SubsystemAPI, "API Error during site creation", map[string]interface{}{
					"name":       site.Name,
					"error":      err.Error(),
					"error_code": apiErr.StatusCode,
				})
				return retry.NonRetryableError(err)
			}
			logging.Error(subCtx, logging.SubsystemCreate, "Failed to create site", map[string]interface{}{
				"name":  site.Name,
				"error": err.Error(),
			})
			return retry.RetryableError(err)
		}
		return nil
	})

	// Log any errors to tf diagnostics
	if err != nil {
		logging.Error(subCtx, logging.SubsystemCreate, "Failed to create site", map[string]interface{}{
			"error": err.Error(),
		})
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdSite.ID))

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProSitesRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.Error(subCtx, logging.SubsystemRead, "Failed to read the created site", map[string]interface{}{
				"name":    d.Get("name"),
				"summary": readDiags[0].Summary,
			})
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemCreate, "Failed to update the Terraform state for the created site", map[string]interface{}{
			"error": err.Error(),
		})
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProSitesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProSitesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	siteIDStr := d.Id()
	siteName := d.Get("name").(string)
	var site *jamfpro.SharedResourceSite

	// Use the retry function for the read operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		siteID, convertErr := strconv.Atoi(siteIDStr)
		if convertErr != nil {
			logging.Error(subCtx, logging.SubsystemRead, "Failed to parse site ID", map[string]interface{}{
				"id":    siteIDStr,
				"error": convertErr.Error(),
			})
			return retry.NonRetryableError(fmt.Errorf("failed to parse site ID: %v", convertErr))
		}

		var apiErr error
		site, apiErr = conn.GetSiteByID(siteID)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.Error(subCtx, logging.SubsystemRead, "Error fetching site by ID, trying by name", map[string]interface{}{
				"id":         siteIDStr,
				"name":       siteName,
				"error":      apiErr.Error(),
				"error_code": apiErrorCode,
			})

			site, apiErr = conn.GetSiteByName(siteName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.Error(subCtx, logging.SubsystemRead, "Error fetching site by name", map[string]interface{}{
					"name":       siteName,
					"error":      apiErr.Error(),
					"error_code": apiErrByNameCode,
				})
				return retry.RetryableError(apiErr)
			}
		}

		if site != nil {
			logging.Info(subCtx, logging.SubsystemRead, "Successfully fetched site", map[string]interface{}{
				"id":   siteIDStr,
				"name": site.Name,
			})

			// Set resource values into Terraform state
			if err := d.Set("id", strconv.Itoa(site.ID)); err != nil {
				return retry.RetryableError(err)
			}
			if err := d.Set("name", site.Name); err != nil {
				return retry.RetryableError(err)
			}
		}

		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemRead, "Failed to read site", map[string]interface{}{
			"id":    siteIDStr,
			"name":  siteName,
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	return nil
}

// ResourceJamfProSitesUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func ResourceJamfProSitesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)

	// Initialize variables
	siteID, err := strconv.Atoi(d.Id())
	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to parse site ID for update", map[string]interface{}{
			"error": err.Error(),
			"id":    d.Id(),
		})
		return diag.FromErr(fmt.Errorf("failed to parse site ID: %v", err))
	}

	siteName := d.Get("name").(string)

	// construct the resource object
	site, err := constructJamfProSite(ctx, d)
	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to construct the site for Terraform update", map[string]interface{}{
			"error": err.Error(),
		})
		return diag.FromErr(fmt.Errorf("failed to construct the site for Terraform update: %w", err))
	}

	// update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateSiteByID(siteID, site)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.Error(subCtx, logging.SubsystemUpdate, "Failed to update site by ID, trying by name", map[string]interface{}{
				"error":      apiErr.Error(),
				"error_code": apiErrorCode,
				"id":         siteID,
				"name":       siteName,
			})

			_, apiErrByName := conn.UpdateSiteByName(siteName, site)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.Error(subCtx, logging.SubsystemUpdate, "API error during site update by name", map[string]interface{}{
					"error":      apiErrByName.Error(),
					"error_code": apiErrByNameCode,
					"name":       siteName,
				})
				return retry.RetryableError(apiErrByName)
			}
		}

		logging.Info(subCtx, logging.SubsystemUpdate, "Successfully updated site", map[string]interface{}{
			"name": site.Name,
			"id":   siteID,
		})
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to update site", map[string]interface{}{
			"error": err.Error(),
			"id":    siteID,
			"name":  siteName,
		})
		return diag.FromErr(err)
	}

	// Retry reading the site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProSitesRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.Error(subCtx, logging.SubsystemUpdate, "Failed to read site after update", map[string]interface{}{
				"summary": readDiags[0].Summary,
				"id":      siteID,
			})
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to synchronize Terraform state after site update", map[string]interface{}{
			"error": err.Error(),
			"id":    siteID,
		})
		return diag.FromErr(err)
	}

	return nil
}

// ResourceJamfProSitesDelete is responsible for deleting a Jamf Pro Site.
func ResourceJamfProSitesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize variables
	siteIDStr := d.Id()
	siteName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			logging.Error(subCtx, logging.SubsystemDelete, "Failed to parse site ID for deletion", map[string]interface{}{
				"error": err.Error(),
				"id":    siteIDStr,
			})
			return retry.NonRetryableError(fmt.Errorf("failed to parse site ID: %v", err))
		}

		// Attempt to delete the site by ID first
		apiErr := conn.DeleteSiteByID(siteID)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.Error(subCtx, logging.SubsystemDelete, "Failed to delete site by ID, trying by name", map[string]interface{}{
				"error":      apiErr.Error(),
				"error_code": apiErrorCode,
				"id":         siteIDStr,
				"name":       siteName,
			})

			// If the delete by ID fails, try deleting by name
			apiErr = conn.DeleteSiteByName(siteName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode // Extract the HTTP status code for name-based error
				}

				logging.Error(subCtx, logging.SubsystemDelete, "API error during site deletion by name", map[string]interface{}{
					"error":      apiErr.Error(),
					"error_code": apiErrByNameCode,
					"name":       siteName,
				})
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		// Log the final error using the initialized subsystem logger
		logging.Error(subCtx, logging.SubsystemDelete, "Failed to delete site", map[string]interface{}{
			"id":    siteIDStr,
			"name":  siteName,
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Log the successful removal of the site from the Terraform state
	logging.Info(subCtx, logging.SubsystemDelete, "Successfully removed site from Terraform state", map[string]interface{}{
		"id":   siteIDStr,
		"name": siteName,
	})

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
