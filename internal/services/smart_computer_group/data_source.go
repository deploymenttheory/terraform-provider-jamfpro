// smartcomputergroup_data_source.go
package smart_computer_group

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceJamfProSmartComputerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the computer group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique name of the Jamf Pro computer group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean selection to state if the group is a Smart group or not. If false then the group is a static group.",
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the smart computer group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the smart computer group.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either 'and', 'or', or blank.",
						},
						"search_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: fmt.Sprintf("The type of smart group search operator. Allowed values are '%v'", getCriteriaOperators()),
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Search value for the smart group criteria to match with.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Opening parenthesis flag used during smart group construction.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Closing parenthesis flag used during smart group construction.",
						},
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceComputerGroup
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetComputerGroupByName(name)
		} else {
			resource, apiErr = client.GetComputerGroupByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		return diag.FromErr(fmt.Errorf("failed to read Smart Computer Group with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		//nolint:err113
		return diag.FromErr(fmt.Errorf("the Jamf Pro Smart Computer Group was not found"))
	}

	d.SetId(fmt.Sprintf("%d", resource.ID))
	return updateState(d, resource)
}
