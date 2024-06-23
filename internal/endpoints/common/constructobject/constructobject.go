// common/constructobject.go
package constructobject

import "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

// ConstructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func ConstructSharedResourceSite(suppliedSite []interface{}) *jamfpro.SharedResourceSite {
	var outSite *jamfpro.SharedResourceSite
	if len(suppliedSite) > 0 && suppliedSite[0].(map[string]interface{})["id"].(int) > 0 {
		outSite = &jamfpro.SharedResourceSite{
			ID: suppliedSite[0].(map[string]interface{})["id"].(int),
		}
	} else {
		outSite = nil
	}

	return outSite
}

// ConstructSharedResourceCategory constructs a SharedResourceCategory object from the provided schema data,
// setting default values if none are presented.
func ConstructSharedResourceCategory(data []interface{}) jamfpro.SharedResourceCategory {
	// Check if 'category' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		category := data[0].(map[string]interface{})

		// Return the 'category' object with data from the schema
		return jamfpro.SharedResourceCategory{
			ID:   category["id"].(int),
			Name: category["name"].(string),
		}
	}

	// Return default 'category' values if no data is provided or it is empty
	return jamfpro.SharedResourceCategory{
		ID:   -1,                     // Default ID
		Name: "No category assigned", // Default name
	}
}
