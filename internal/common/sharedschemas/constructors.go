// common/constructobject.go
package sharedschemas

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

const (
	// Values required to unset, not just empty values.
	EmptySiteName     = ""
	EmptySiteId       = -1
	EmptyCategoryName = ""
	EmptyCategoryId   = 0
)

// ConstructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
func ConstructSharedResourceSite(siteId int) *jamfpro.SharedResourceSite {
	if siteId == 0 || siteId == -1 {
		return &jamfpro.SharedResourceSite{
			ID:   EmptySiteId,
			Name: EmptySiteName,
		}
	}

	return &jamfpro.SharedResourceSite{ID: siteId}
}

// ConstructSharedResourceCategory constructs a SharedResourceCategory object from the provided schema data,
func ConstructSharedResourceCategory(categoryId int) *jamfpro.SharedResourceCategory {
	if categoryId == 0 || categoryId == -1 {
		return &jamfpro.SharedResourceCategory{
			ID:   EmptySiteId,
			Name: EmptySiteName,
		}
	}

	return &jamfpro.SharedResourceCategory{ID: categoryId}
}
