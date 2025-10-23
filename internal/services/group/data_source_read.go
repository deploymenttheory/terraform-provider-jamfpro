package group

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of a specific group from Jamf Pro using either its unique Name or its Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	name := d.Get("name").(string)
	groupPlatformId := d.Get("group_platform_id").(string)
	groupJamfProId := d.Get("group_jamfpro_id").(string)
	groupType := d.Get("group_type").(string)

	if name == "" && groupPlatformId == "" && groupJamfProId == "" {
		return diag.FromErr(errMustProvideOne)
	}

	if name != "" && groupJamfProId != "" {
		return diag.FromErr(errNameAndIDConflict)
	}

	if (name != "" || groupJamfProId != "") && groupType == "" {
		return diag.FromErr(errGroupTypeRequired)
	}

	var resource *jamfpro.ResourceGroup
	var lookupMethod, lookupValue string
	var apiErr error

	switch {
	case name != "" && groupType == "COMPUTER":
		resource, apiErr = client.GetComputerGroupByJamfProName(name)
		lookupMethod = "name (COMPUTER)"
		lookupValue = name
	case name != "" && groupType == "MOBILE":
		resource, apiErr = client.GetMobileGroupByJamfProName(name)
		lookupMethod = "name (MOBILE)"
		lookupValue = name
	case groupPlatformId != "":
		resource, apiErr = client.GetGroupByID(groupPlatformId)
		lookupMethod = "group_platform_id"
		lookupValue = groupPlatformId
	case groupJamfProId != "" && groupType == "COMPUTER":
		resource, apiErr = client.GetComputerGroupByJamfProID(groupJamfProId)
		lookupMethod = "group_jamfpro_id (COMPUTER)"
		lookupValue = groupJamfProId
	case groupJamfProId != "" && groupType == "MOBILE":
		resource, apiErr = client.GetMobileGroupByJamfProID(groupJamfProId)
		lookupMethod = "group_jamfpro_id (MOBILE)"
		lookupValue = groupJamfProId
	}

	if apiErr != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Group with %s '%s': %w", lookupMethod, lookupValue, apiErr))
	}

	d.SetId(resource.GroupPlatformId)
	return updateState(d, resource)
}

// updateState sets the Terraform state from the ResourceGroup object.
func updateState(d *schema.ResourceData, resource *jamfpro.ResourceGroup) diag.Diagnostics {
	fields := map[string]interface{}{
		"group_platform_id": resource.GroupPlatformId,
		"group_jamfpro_id":  resource.GroupJamfProId,
		"name":              resource.GroupName,
		"group_type":        resource.GroupType,
		"group_description": resource.GroupDescription,
		"smart":             resource.Smart,
		"membership_count":  resource.MembershipCount,
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set %s: %w", k, err))
		}
	}

	return nil
}
