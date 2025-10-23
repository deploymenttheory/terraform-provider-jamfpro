package jamf_protect_plan

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
	ErrProvideEitherIDOrName = fmt.Errorf("please provide either 'id' or 'name', not both")
	ErrPlanNotFound          = fmt.Errorf("jamf Protect Plan not found using the provided identifier")
	ErrFailedToReadPlan      = fmt.Errorf("failed to read Jamf Protect Plan after retries")
)

// DataSourceJamfProtectPlan provides information about a Jamf Protect Plan
func DataSourceJamfProtectPlan() *schema.Resource {
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
				Description: "The ID of the Jamf Protect Plan",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Jamf Protect Plan",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Jamf Protect Plan",
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
	}
}

// dataSourceRead fetches Jamf Protect Plan details
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	planID := d.Get("id").(string)
	planName := d.Get("name").(string)

	if planID != "" && planName != "" {
		return diag.FromErr(fmt.Errorf("%w", ErrProvideEitherIDOrName))
	}

	var resource *jamfpro.ResourceJamfProtectPlan
	var identifier string

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		params := url.Values{}
		var plansList *jamfpro.ResponseJamfProtectPlansList
		var apiErr error
		plansList, apiErr = client.GetJamfProtectPlans(params)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		found := false
		for _, plan := range plansList.Results {
			if planID != "" && plan.ID == planID {
				resource = &plan
				d.SetId(plan.ID)
				identifier = planID
				found = true
				break
			}
			if planName != "" && plan.Name == planName {
				resource = &plan
				d.SetId(plan.ID)
				identifier = planName
				found = true
				break
			}
		}
		if !found {
			return retry.NonRetryableError(fmt.Errorf("%w: '%s'", ErrPlanNotFound, planID+planName))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: '%s', %w", ErrFailedToReadPlan, identifier, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("%w: '%s'", ErrPlanNotFound, identifier))
	}

	if err := d.Set("uuid", resource.UUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("profile_id", resource.ProfileID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("profile_name", resource.ProfileName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope_description", resource.ScopeDescription); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
