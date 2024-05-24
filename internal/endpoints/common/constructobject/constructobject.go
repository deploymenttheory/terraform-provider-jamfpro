// common/constructobject.go
package constructobject

import "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

// ConstructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func ConstructSharedResourceSite(data []interface{}) jamfpro.SharedResourceSite {
	// Check if 'site' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		site := data[0].(map[string]interface{})

		// Return the 'site' object with data from the schema
		return jamfpro.SharedResourceSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	// Return default 'site' values if no data is provided or it is empty
	return jamfpro.SharedResourceSite{
		ID:   -1,     // Default ID
		Name: "None", // Default name
	}
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
