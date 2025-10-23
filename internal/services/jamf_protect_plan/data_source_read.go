package jamf_protect_plan

import (
	"context"
	"fmt"
	"net/url"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches Jamf Protect Plan details
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
