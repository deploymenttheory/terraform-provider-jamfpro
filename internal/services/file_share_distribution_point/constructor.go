package file_share_distribution_point

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceFileShareDistributionPoint object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceFileShareDistributionPoint, error) {
	resource := &jamfpro.ResourceFileShareDistributionPoint{
		Name:                     d.Get("name").(string),
		IP_Address:               d.Get("ip_address").(string),
		IsMaster:                 d.Get("is_master").(bool),
		FailoverPoint:            d.Get("failover_point").(string),
		ConnectionType:           d.Get("connection_type").(string),
		ShareName:                d.Get("share_name").(string),
		SharePort:                d.Get("share_port").(int),
		EnableLoadBalancing:      d.Get("enable_load_balancing").(bool),
		WorkgroupOrDomain:        d.Get("workgroup_or_domain").(string),
		ReadOnlyUsername:         d.Get("read_only_username").(string),
		ReadOnlyPassword:         d.Get("read_only_password").(string),
		ReadWriteUsername:        d.Get("read_write_username").(string),
		ReadWritePassword:        d.Get("read_write_password").(string),
		NoAuthenticationRequired: d.Get("no_authentication_required").(bool),
		HTTPDownloadsEnabled:     d.Get("https_downloads_enabled").(bool),
		Port:                     d.Get("https_port").(int),
		Context:                  d.Get("https_share_path").(string),
		UsernamePasswordRequired: d.Get("https_username_password_required").(bool),
		HTTPUsername:             d.Get("https_username").(string),
		HTTPPassword:             d.Get("https_password").(string),
		Protocol:                 d.Get("protocol").(string),
		HTTPURL:                  d.Get("http_url").(string),
	}

	resourceXML, err := common.SerializeAndRedactXML(resource, []string{"ReadOnlyPassword", "ReadWritePassword", "HTTPPassword"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro File Share Distribution Point XML:\n%s\n", string(resourceXML))

	return resource, nil
}
