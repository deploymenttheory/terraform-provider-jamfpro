// common/constructobject.go
package sharedschemas

import "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

const (
	// Values required to unset, not just empty values.
	EmptySiteName     = ""
	EmptySiteId       = -1
	EmptyCategoryName = ""
	EmptyCategoryId   = 0
)

// ConstructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func ConstructSharedResourceSite(site_id int) *jamfpro.SharedResourceSite {
	if site_id == 0 || site_id == -1 {
		return &jamfpro.SharedResourceSite{
			ID:   EmptySiteId,
			Name: EmptySiteName,
		}
	}

	return &jamfpro.SharedResourceSite{ID: site_id}
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
