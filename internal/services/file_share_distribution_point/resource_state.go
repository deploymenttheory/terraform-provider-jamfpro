// filesharedistributionpoints_state.go
package file_share_distribution_point

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest File Share Distribution Point information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceFileShareDistributionPoint) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]any{
		"id":                               strconv.Itoa(resp.ID),
		"name":                             resp.Name,
		"ip_address":                       resp.IPAddress,
		"is_master":                        resp.IsMaster,
		"failover_point":                   resp.FailoverPoint,
		"failover_point_url":               resp.FailoverPointURL,
		"enable_load_balancing":            resp.EnableLoadBalancing,
		"local_path":                       resp.LocalPath,
		"ssh_username":                     resp.SSHUsername,
		"connection_type":                  resp.ConnectionType,
		"share_name":                       resp.ShareName,
		"workgroup_or_domain":              resp.WorkgroupOrDomain,
		"share_port":                       resp.SharePort,
		"read_only_username":               resp.ReadOnlyUsername,
		"https_downloads_enabled":          resp.HTTPDownloadsEnabled,
		"http_url":                         resp.HTTPURL,
		"https_share_path":                 resp.Context,
		"protocol":                         resp.Protocol,
		"https_port":                       resp.Port,
		"no_authentication_required":       resp.NoAuthenticationRequired,
		"https_username_password_required": resp.UsernamePasswordRequired,
		"https_username":                   resp.HTTPUsername,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
