// filesharedistributionpoints_object.go
package file_share_distribution_point

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructs JamfProFileShareDistributionPoint from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceFileShareDistributionPoint, error) {
	resource := &jamfpro.ResourceFileShareDistributionPoint{
		ShareName:                 d.Get("share_name").(string),
		Workgroup:                 d.Get("workgroup").(string),
		Port:                      d.Get("port").(int),
		ReadWriteUsername:         d.Get("read_write_username").(string),
		ReadWritePassword:         d.Get("read_write_password").(string),
		ReadOnlyUsername:          d.Get("read_only_username").(string),
		ReadOnlyPassword:          d.Get("read_only_password").(string),
		ID:                        d.Get("id").(string),
		Name:                      d.Get("name").(string),
		ServerName:                d.Get("server_name").(string),
		Principal:                 d.Get("principal").(bool),
		BackupDistributionPointID: d.Get("backup_distribution_point_id").(string),
		SSHUsername:               d.Get("ssh_username").(string),
		SSHPassword:               d.Get("ssh_password").(string),
		LocalPathToShare:          d.Get("local_path_to_share").(string),
		FileSharingConnectionType: d.Get("file_sharing_connection_type").(string),
		HTTPSEnabled:              d.Get("https_enabled").(bool),
		HTTPSPort:                 d.Get("https_port").(int),
		HTTPSContext:              d.Get("https_context").(string),
		HTTPSSecurityType:         d.Get("https_security_type").(string),
		HTTPSUsername:             d.Get("https_username").(string),
		HTTPSPassword:             d.Get("https_password").(string),
		EnableLoadBalancing:       d.Get("enable_load_balancing").(bool),
	}

	resourceJSON, err := common.SerializeAndRedactJSON(resource, []string{"ReadOnlyPassword", "ReadWritePassword", "HTTPSPassword"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro File Share Distribution Point JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
