// smartcomputergroup_data_source.go
package smart_computer_group

import (
	"fmt"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
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
