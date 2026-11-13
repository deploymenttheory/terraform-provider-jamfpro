package advanced_mobile_device_search

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCriteriaSearchTypeAcceptsDateOperators(t *testing.T) {
	criteria := ResourceJamfProAdvancedMobileDeviceSearches().Schema["criteria"].Elem.(*schema.Resource)
	validate := criteria.Schema["search_type"].ValidateFunc

	for _, searchType := range []string{"before (yyyy-mm-dd)", "after (yyyy-mm-dd)", "is", "in less than x days"} {
		if _, errs := validate(searchType, "search_type"); len(errs) > 0 {
			t.Errorf("search_type %q rejected: %v", searchType, errs)
		}
	}

	if _, errs := validate("before", "search_type"); len(errs) == 0 {
		t.Error("expected an unknown search_type to be rejected")
	}
}
