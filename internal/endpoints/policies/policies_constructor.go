package policies

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructJamfProPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {

	// Non computed values first
	policy := &jamfpro.ResourcePolicy{
		General: jamfpro.PolicySubsetGeneral{
			Name: d.Get("name").(string),
		},
	}

	return policy, nil
}
