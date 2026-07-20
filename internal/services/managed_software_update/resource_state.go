// state.go
package managed_software_update

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResponseManagedSoftwareUpdatePlan
// information from the Jamf Pro API. 'group' and 'device' are intentionally left untouched:
// the plan detail response always echoes the concrete target device (Jamf fans a group out
// into one plan per member at creation time), so re-deriving 'group' vs 'device' from it would
// misclassify every group-scoped plan as device-scoped on the very next read.
func updateState(d *schema.ResourceData, plan *jamfpro.ResponseManagedSoftwareUpdatePlan) diag.Diagnostics {
	if plan == nil {
		return diag.Errorf("no managed software update plan found in the response")
	}

	d.SetId(plan.PlanUuid)
	if err := d.Set("plan_uuid", plan.PlanUuid); err != nil {
		return diag.FromErr(fmt.Errorf("error setting plan_uuid: %v", err))
	}

	if err := d.Set("update_action", plan.UpdateAction); err != nil {
		return diag.FromErr(fmt.Errorf("error setting update_action: %v", err))
	}

	if err := d.Set("version_type", plan.VersionType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting version_type: %v", err))
	}

	if err := d.Set("specific_version", plan.SpecificVersion); err != nil {
		return diag.FromErr(fmt.Errorf("error setting specific_version: %v", err))
	}

	if err := d.Set("build_version", plan.BuildVersion); err != nil {
		return diag.FromErr(fmt.Errorf("error setting build_version: %v", err))
	}

	if err := d.Set("max_deferrals", plan.MaxDeferrals); err != nil {
		return diag.FromErr(fmt.Errorf("error setting max_deferrals: %v", err))
	}

	if err := d.Set("force_install_local_date_time", plan.ForceInstallLocalDateTime); err != nil {
		return diag.FromErr(fmt.Errorf("error setting force_install_local_date_time: %v", err))
	}

	return nil
}
