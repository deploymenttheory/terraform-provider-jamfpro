package sharedschema_constructors

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetSite(d *schema.ResourceData) (*jamfpro.SharedResourceSite, error) {
	var out *jamfpro.SharedResourceSite

	hcl := d.Get("site_id").([]interface{})

	if len(hcl) > 0 {
		out = &jamfpro.SharedResourceSite{
			ID: hcl[0].(map[string]interface{})["id"].(int),
		}
	} else {
		out = &jamfpro.SharedResourceSite{
			ID: -1,
		}
	}

	return out, nil

}
