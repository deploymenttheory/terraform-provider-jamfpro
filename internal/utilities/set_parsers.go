package utilities

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"
)

func SetParserJamfProSite(inputs []interface{}) ([]*jamfpro.SharedResourceSite, error) {
	out := make([]*jamfpro.SharedResourceSite, len(inputs))
	for i, input := range inputs { // Should only be one

		if i > 1 {
			return nil, fmt.Errorf(logging.MsgTooManyResourcesProvided, "site")
		}

		param := input.(map[string]interface{})
		config := &jamfpro.SharedResourceSite{}

		if v, ok := param["id"]; ok {
			config.ID = v.(int)
		}

		if v, ok := param["name"]; ok {
			config.Name = v.(string)
		}

		out[i] = config
	}

	return out, nil
}
