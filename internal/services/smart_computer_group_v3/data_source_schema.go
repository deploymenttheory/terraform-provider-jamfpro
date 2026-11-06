package smart_computer_group_v3

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type smartComputerGroupV3DataSource struct {
	client interface {
		GetSmartComputerGroupByIDV3WithContext(context.Context, string) (*jamfpro.ResourceSmartComputerGroupV3, error)
		GetSmartComputerGroupByNameV3WithContext(
			context.Context, string,
		) (*jamfpro.ResponseSmartComputerGroupListItemV3, error)
	}
}

var (
	_ datasource.DataSource                     = &smartComputerGroupV3DataSource{}
	_ datasource.DataSourceWithConfigure        = &smartComputerGroupV3DataSource{}
	_ datasource.DataSourceWithConfigValidators = &smartComputerGroupV3DataSource{}
)

// NewSmartComputerGroupV3FrameworkDataSource retrieves a smart computer group using v3.
func NewSmartComputerGroupV3FrameworkDataSource() datasource.DataSource {
	return &smartComputerGroupV3DataSource{}
}

func (d *smartComputerGroupV3DataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_smart_computer_group_v3"
}

func (d *smartComputerGroupV3DataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok || client == nil {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got %T.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *smartComputerGroupV3DataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
	}
}

func (d *smartComputerGroupV3DataSource) Schema(
	ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a Jamf Pro smart computer group using v3, by numeric ID or case-insensitive name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true, Computed: true, Description: "The numeric Jamf Pro group ID. Specify either id or name.",
				Validators: []validator.String{stringvalidator.RegexMatches(positiveID, "must be a positive numeric group ID")},
			},
			"name": schema.StringAttribute{
				Optional: true, Computed: true, Description: "The group name. Specify either id or name.",
				Validators: []validator.String{stringvalidator.LengthBetween(1, 255)},
			},
			"description": schema.StringAttribute{Computed: true, Description: "The group description."},
			"site_id":     schema.StringAttribute{Computed: true, Description: "The site ID, or -1 for None."},
			"criteria": schema.ListNestedAttribute{
				Computed: true, Description: "Group criteria in priority order.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: commonschema.CriteriaDataSource(ctx).NestedObject.Attributes,
				},
			},
		},
	}
}
