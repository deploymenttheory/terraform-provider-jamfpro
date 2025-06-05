package volume_purchasing_locations

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct creates a VolumePurchasingLocationCreateUpdateRequest from the provided schema data
func construct(d *schema.ResourceData) (*jamfpro.VolumePurchasingLocationCreateUpdateRequest, error) {
	request := &jamfpro.VolumePurchasingLocationCreateUpdateRequest{
		Name:                                  d.Get("name").(string),
		ServiceToken:                          d.Get("service_token").(string),
		AutomaticallyPopulatePurchasedContent: d.Get("automatically_populate_purchased_content").(bool),
		SendNotificationWhenNoLongerAssigned:  d.Get("send_notification_when_no_longer_assigned").(bool),
		AutoRegisterManagedUsers:              d.Get("auto_register_managed_users").(bool),
		SiteID:                                d.Get("site_id").(string),
	}

	return request, nil
}
