// dockitems_data_source.go
package dockitems

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDockItems provides information about a specific dock item in Jamf Pro.
func DataSourceJamfProDockItems() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProDockItemsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the dock item.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique name of the jamf pro dock item.",
				Computed:    true,
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

// DataSourceJamfProDockItemsRead fetches the details of a specific dock item from Jamf Pro using either its unique Name or its Id.
func DataSourceJamfProDockItemsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var dockItem *jamfpro.ResponseDockItem
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		dockItemName := v.(string)
		dockItem, err = conn.GetDockItemsByName(dockItemName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch dock item by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		dockItemID, convertErr := strconv.Atoi(v.(string))
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert dock item ID to integer: %v", convertErr))
		}
		dockItem, err = conn.GetDockItemsByID(dockItemID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch dock item by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if dockItem == nil {
		return diag.FromErr(fmt.Errorf("dock item not found"))
	}

	// Set the data source attributes using the fetched data
	d.SetId(fmt.Sprintf("%d", dockItem.ID))
	d.Set("name", dockItem.Name)
	d.Set("type", dockItem.Type)
	d.Set("path", dockItem.Path)
	d.Set("contents", dockItem.Contents)

	return nil
}
