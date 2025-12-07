package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_computer_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_mobile_device_group"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		dock_item.NewDockItemFrameworkResource,
		smart_computer_group.NewSmartComputerGroupFrameworkResource,
		smart_mobile_device_group.NewSmartMobileDeviceGroupFrameworkResource,
	}
}
