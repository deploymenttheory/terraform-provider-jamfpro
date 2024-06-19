// filesharedistributionpoints_state.go
package filesharedistributionpoints

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest MacOS Configuration Profile information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceFileShareDistributionPoint) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if distribution point data exists
	if resource != nil {
		// Organize state updates into a map
		resourceData := map[string]interface{}{
			"id":                    strconv.Itoa(resource.ID),
			"name":                  resource.Name,
			"ip_address":            resource.IPAddress,
			"is_master":             resource.IsMaster,
			"failover_point":        resource.FailoverPoint,
			"failover_point_url":    resource.FailoverPointURL,
			"enable_load_balancing": resource.EnableLoadBalancing,
			"local_path":            resource.LocalPath,
			"ssh_username":          resource.SSHUsername,
			// "password": resource.Password,  // sensitive field, not included in state
			"connection_type":                  resource.ConnectionType,
			"share_name":                       resource.ShareName,
			"workgroup_or_domain":              resource.WorkgroupOrDomain,
			"share_port":                       resource.SharePort,
			"read_only_username":               resource.ReadOnlyUsername,
			"https_downloads_enabled":          resource.HTTPDownloadsEnabled,
			"http_url":                         resource.HTTPURL,
			"https_share_path":                 resource.Context,
			"protocol":                         resource.Protocol,
			"https_port":                       resource.Port,
			"no_authentication_required":       resource.NoAuthenticationRequired,
			"https_username_password_required": resource.UsernamePasswordRequired,
			"https_username":                   resource.HTTPUsername,
		}

		for key, val := range resourceData {
			if err := d.Set(key, val); err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return diags
}
