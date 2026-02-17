package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_computer_group_v2"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_mobile_device_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		dock_item.NewDockItemFrameworkResource,
		smart_computer_group_v2.NewSmartComputerGroupV2FrameworkResource,
		smart_mobile_device_group_v1.NewSmartMobileDeviceGroupV1FrameworkResource,
	}
}
