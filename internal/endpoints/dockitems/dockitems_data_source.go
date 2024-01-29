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
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var dockItem *jamfpro.ResourceDockItem
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		dockItemName := v.(string)
		dockItem, err = conn.GetDockItemByName(dockItemName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch dock item by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		dockItemID, err := strconv.Atoi(v.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse dock item ID: %v", err))
		}
		dockItem, err = conn.GetDockItemByID(dockItemID)
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
	if err := d.Set("name", dockItem.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'name': %v", err))
	}
	if err := d.Set("type", dockItem.Type); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'type': %v", err))
	}
	if err := d.Set("path", dockItem.Path); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'path': %v", err))
	}
	if err := d.Set("contents", dockItem.Contents); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'contents': %v", err))
	}

	// Set the Terraform state ID for the dock item
	d.SetId(fmt.Sprintf("%d", dockItem.ID))

	return nil

}
