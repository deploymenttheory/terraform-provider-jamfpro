package jamf_protect_plans

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrFailedToReadPlans = fmt.Errorf("failed to read Jamf Protect Plans after retries")
	ErrNoPlansFound      = fmt.Errorf("no Jamf Protect Plans found")
)

// DataSourceJamfProtectPlans provides information about Jamf Protect Plans
func DataSourceJamfProtectPlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"plans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all Jamf Protect Plans",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Jamf Protect Plan",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the Jamf Protect Plan",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Jamf Protect Plan",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the Jamf Protect Plan",
						},
						"profile_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Profile ID associated with the plan",
						},
						"profile_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Profile name associated with the plan",
						},
						"scope_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the plan's scope",
						},
					},
				},
			},
		},
	}
}

// dataSourceRead fetches Jamf Protect Plan details
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var plansList *jamfpro.ResponseJamfProtectPlansList
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		params := url.Values{}
		var apiErr error
		plansList, apiErr = client.GetJamfProtectPlans(params)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%w", ErrFailedToReadPlans))
	}

	if plansList == nil || len(plansList.Results) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("%w", ErrNoPlansFound))
	}

	plans := make([]map[string]any, 0, len(plansList.Results))
	for _, plan := range plansList.Results {
		plans = append(plans, map[string]any{
			"id":                plan.ID,
			"uuid":              plan.UUID,
			"name":              plan.Name,
			"description":       plan.Description,
			"profile_id":        plan.ProfileID,
			"profile_name":      plan.ProfileName,
			"scope_description": plan.ScopeDescription,
		})
	}

	if err := d.Set("plans", plans); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("jamf_protect_plans")

	return nil
}
