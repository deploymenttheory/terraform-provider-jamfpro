package provider

import (
	"context"

	jamfProStaticComputerGroup "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/static_computer_group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *frameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		jamfProStaticComputerGroup.NewStaticComputerGroupFrameworkDataSource,
	}
}
