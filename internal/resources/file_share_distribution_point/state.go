// filesharedistributionpoints_state.go
package file_share_distribution_point

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest File Share Distribution Point information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceFileShareDistributionPoint) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":                           resp.ID,
		"name":                         resp.Name,
		"server_name":                  resp.ServerName,
		"principal":                    resp.Principal,
		"backup_distribution_point_id": resp.BackupDistributionPointID,
		"enable_load_balancing":        resp.EnableLoadBalancing,
		"local_path_to_share":          resp.LocalPathToShare,
		"ssh_username":                 resp.SSHUsername,
		"file_sharing_connection_type": resp.FileSharingConnectionType,
		"share_name":                   resp.ShareName,
		"workgroup":                    resp.Workgroup,
		"port":                         resp.Port,
		"read_only_username":           resp.ReadOnlyUsername,
		"https_enabled":                resp.HTTPSEnabled,
		"https_port":                   resp.HTTPSPort,
		"https_context":                resp.HTTPSContext,
		"https_security_type":          resp.HTTPSSecurityType,
		"https_username":               resp.HTTPSUsername,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
