// dockitems_data_source.go
package dockitems

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDockItems provides information about specific Jamf Pro Dock Items by their ID or Name.
func DataSourceJamfProDockItems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProDockItemsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the dock item.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the dock item.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the dock item (App/File/Folder).",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the dock item.",
			},
			"contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the dock item.",
			},
		},
	}
}

// dataSourceJamfProDockItemsRead fetches the details of specific dock items from Jamf Pro using either their unique Name or Id.
func dataSourceJamfProDockItemsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDockItem, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// read operation
	dockItem, err := conn.GetDockItemByID(resourceIDInt)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceDockItem, resourceID, err.Error(), apiErrorCode)
		return diags
	}

	// Check if dockItem data exists
	if dockItem != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", strconv.Itoa(dockItem.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", dockItem.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("type", dockItem.Type); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("path", dockItem.Path); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("contents", dockItem.Contents); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
