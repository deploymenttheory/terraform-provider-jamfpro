package dockitems

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDockItems provides information about specific Jamf Pro Dock Items
func DataSourceJamfProDockItems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the dock item.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the dock item.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the dock item.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the dock item.",
			},
			"contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The contents of the dock item.",
			},
		},
	}
}

// dataSourceRead fetches dock items details from Jamf Pro
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	searchID := d.Get("id").(string)
	searchName := d.Get("name").(string)

	if searchID == "" && searchName == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var dockItemsList *jamfpro.ResponseDockItemsList
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		dockItemsList, apiErr = client.GetDockItems()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch list of dock items: %v", err))
	}

	var matchedID string
	if searchID != "" {
		for _, item := range dockItemsList.DockItems {
			if fmt.Sprintf("%d", item.ID) == searchID {
				matchedID = searchID
				break
			}
		}
	} else {
		for _, item := range dockItemsList.DockItems {
			if item.Name == searchName {
				matchedID = fmt.Sprintf("%d", item.ID)
				break
			}
		}
	}

	if matchedID == "" {
		return diag.FromErr(fmt.Errorf("no dock item found matching the provided criteria"))
	}

	var resource *jamfpro.ResourceDockItem
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetDockItemByID(matchedID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read dock item with ID '%s': %v", matchedID, err))
	}

	d.SetId(matchedID)
	return updateState(d, resource)
}
