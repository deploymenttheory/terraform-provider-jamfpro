package guid_list_sharder

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DataSourceName = "jamfpro_guid_list_sharder"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &guidListSharderDataSource{}
	_ datasource.DataSourceWithConfigure = &guidListSharderDataSource{}
)

func NewGuidListSharderDataSource() datasource.DataSource {
	return &guidListSharderDataSource{}
}

type guidListSharderDataSource struct {
	client *jamfpro.Client
}

func (d *guidListSharderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *guidListSharderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

func (d *guidListSharderDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves object IDs from Jamf Pro API and distributes them into configurable shards for progressive rollouts and phased deployments. " +
			"Queries computer inventory, mobile device groups, computer groups, or user list sources with optional filtering, then applies sharding strategies (random, sequential, or percentage-based) " +
			"to distribute results. Output shards are sets that can be directly used in static groups, policies, and other resources requiring ID collections.\n\n" +
			"**API Endpoints:** Computer Inventory (Pro API), Computer Groups (Classic API), Mobile Device Groups (Classic API), Users (Classic API)\n\n" +
			"**Common Use Cases:** Policy rollouts, group management, OS update rings, pilot testing, A/B testing for configurations\n\n" +
			"For detailed examples and best practices, see the provider documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"source_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The source type to query IDs from and shard. " +
					"`computer_inventory` queries computer inventory for computer-based policies and groups. " +
					"`mobile_device_inventory` queries mobile device inventory. " +
					"`computer_group_membership` queries computer group membership (requires `group_id`). " +
					"`mobile_device_group_membership` queries mobile device group membership (requires `group_id`). " +
					"`user_accounts` queries Jamf Pro user accounts for user-based policies.",
				Validators: []validator.String{
					stringvalidator.OneOf("computer_inventory", "mobile_device_inventory", "computer_group_membership", "mobile_device_group_membership", "user_accounts"),
				},
			},
			"group_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the group to query members from. " +
					"Required when `source_type` is `computer_group_membership` or `mobile_device_group_membership`, ignored otherwise. " +
					"Use this to split an existing group's membership into multiple new groups for targeted policy application.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"must be a valid numeric ID",
					),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				MarkdownDescription: "Number of equally-sized shards to create (minimum 1). " +
					"Use with `round-robin` strategy. Conflicts with `shard_percentages` and `shard_sizes`. " +
					"Creates shards named `shard_0`, `shard_1`, ..., `shard_N-1`. " +
					"For custom-sized shards (e.g., 10% pilot, 30% broader, 60% full), use `shard_percentages` with `percentage` strategy instead.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.ExactlyOneOf(path.MatchRoot("shard_percentages"), path.MatchRoot("shard_sizes")),
				},
			},
			"shard_percentages": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of percentages for custom-sized shards. Use with `percentage` strategy. Conflicts with `shard_count` and `shard_sizes`. " +
					"Values must be non-negative integers that sum to exactly 100. " +
					"Example: `[10, 30, 60]` creates 10% pilot, 30% broader pilot, 60% full rollout. " +
					"Common patterns: `[5, 15, 80]` (OS update rings), `[33, 33, 34]` (A/B/C testing). " +
					"Last shard receives all remaining IDs to prevent loss.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.AtLeast(0)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_sizes")),
				},
			},
			"shard_sizes": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of absolute shard sizes (exact number of IDs per shard). Use with `size` strategy. Conflicts with `shard_count` and `shard_percentages`. " +
					"Values must be positive integers or -1 (which means 'all remaining'). Only the last element can be -1. " +
					"Example: `[50, 200, -1]` creates 50 pilot computers, 200 broader rollout, remainder for full deployment. " +
					"Common patterns: `[10, 30, -1]` (controlled pilot expansion), `[100, 100, 100, -1]` (fixed-size rings). " +
					"Use this when you need exact capacity constraints (e.g., support team handles exactly 50 pilot devices).",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.Any(
						int64validator.AtLeast(1),
						int64validator.OneOf(-1),
					)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_percentages")),
				},
			},
			"strategy": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The distribution strategy for sharding IDs. " +
					"`round-robin` distributes in circular order (guarantees equal sizes, optional seed for reproducibility). " +
					"`percentage` distributes by specified percentages (requires `shard_percentages`, optional seed for reproducibility). " +
					"`size` distributes by absolute sizes (requires `shard_sizes`, optional seed for reproducibility). " +
					"`rendezvous` uses Highest Random Weight algorithm (always deterministic, minimal disruption when shard count changes, requires seed). " +
					"See the provider documentation for detailed comparison.",
				Validators: []validator.String{
					stringvalidator.OneOf("round-robin", "percentage", "size", "rendezvous"),
				},
			},
			"seed": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional seed value for deterministic distribution. When provided, makes results reproducible across Terraform runs. " +
					"**`round-robin` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies round-robin (reproducible). " +
					"**`percentage` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies percentage split (reproducible). " +
					"**`size` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies size-based split (reproducible). " +
					"**`rendezvous` strategy**: Always deterministic. Seed affects which shard wins for each ID via Highest Random Weight algorithm. " +
					"Use different seeds for different rollouts to distribute pilot burden: Device X might be in shard_0 for OS updates but shard_2 for app deployments.",
			},
			"shards": schema.MapAttribute{
				ElementType: types.SetType{ElemType: types.StringType},
				Computed:    true,
				MarkdownDescription: "Computed map of shard names (`shard_0`, `shard_1`, ...) to sets of IDs. " +
					"Each value is a `set(string)` type, directly compatible with resource attributes expecting ID sets " +
					"(e.g., static group members, policy scope). " +
					"Access with `data.example.shards[\"shard_0\"]`, check size with `length(data.example.shards[\"shard_0\"])`.",
			},
		},
	}
}
