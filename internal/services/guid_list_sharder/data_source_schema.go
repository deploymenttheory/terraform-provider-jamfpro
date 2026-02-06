package guid_list_sharder

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/validate"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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
		MarkdownDescription: "Retrieves managed device and user IDs from Jamf Pro and distributes them into configurable shards for progressive rollouts and phased deployments. " +
			"Supports four sharding strategies: **round-robin** (equal distribution ±1), **percentage** (proportional distribution), **size** (fixed sizes with remainder support), " +
			"and **rendezvous** (HRW algorithm with minimal disruption when scaling). Optional features include deterministic seeding for reproducibility, ID exclusions, " +
			"and reserved ID assignment to specific shards. Automatically filters out unmanaged computer and mobile device sources as they cannot be allocated to a jamf pro group.\n\n" +
			"**Sources:** Computer Inventory (Pro API), Mobile Device Inventory (Pro API), Computer Groups (Classic API), Mobile Device Groups (Classic API), Users (Classic API)\n\n" +
			"**Sharding Strategies:** Round-robin, Percentage, Size, Rendezvous (HRW)\n\n" +
			"**Optional Features:** Seed (determinism), Exclude IDs, Reserved IDs\n\n" +
			"**Common Use Cases:** OS update rings, pilot testing, A/B testing, phased policy rollouts, canary deployments\n\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"source_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The source type to query IDs from before applying exclusions, reservations and applyingthe chosen sharding strategy. " +
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
					"Use this to split an existing group's membership into multiple new groups for targeted policy application. e.g `group_id = \"10\"`",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"must be a valid numeric ID",
					),
					validate.RequiredWhenOneOf("source_type", "computer_group_membership", "mobile_device_group_membership"),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				MarkdownDescription: "Number of shards to create (minimum 1). " +
					"Required for `round-robin` and `rendezvous` strategies. Mutually exclusive with `shard_percentages` and `shard_sizes`. " +
					"Round-robin guarantees equal distribution (±1 ID variance). Rendezvous provides superior stability when shard count changes (~1/n disruption when scaling). " +
					"Creates output shards named `shard_0`, `shard_1`, ..., `shard_N-1`. " +
					"**Example:** `shard_count = 3` with round-robin creates 3 equal shards. " +
					"**Example:** `shard_count = 4` with rendezvous creates 4 shards optimized for minimal redistribution if later scaled to 5 shards. " +
					"For proportional distribution (e.g., 10% pilot, 90% production), use `shard_percentages` instead. " +
					"For fixed-size shards (e.g., 25 dev, 75 staging, rest production), use `shard_sizes` instead.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.ExactlyOneOf(path.MatchRoot("shard_percentages"), path.MatchRoot("shard_sizes")),
					validate.Int64RequiredWhenOneOf("strategy", "round-robin", "rendezvous"),
				},
			},
			"shard_percentages": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of percentages for proportional shard distribution. Required for `percentage` strategy. Mutually exclusive with `shard_count` and `shard_sizes`. " +
					"Array length determines number of shards (e.g., `[10, 30, 60]` creates 3 shards). Values must be non-negative integers summing to exactly 100. " +
					"Percentages are calculated against total ID count (after exclusions). When `reserved_ids` are used, reserved counts are subtracted from targets to maintain percentage accuracy. " +
					"Last shard automatically receives any remainder from rounding to ensure all IDs are distributed. " +
					"**Example:** `[10, 30, 60]` with 1000 IDs creates shard_0=100, shard_1=300, shard_2=600. " +
					"**Example:** `[5, 15, 80]` for OS update rings: 5% test, 15% early adopters, 80% production. " +
					"**Example:** `[33, 33, 34]` for A/B/C testing with near-equal distribution. " +
					"**With Reserved IDs:** Optionally use in conjuction with `reserved_ids` to set specific Ids to specific shards. Scenario: if 1000 total IDs with 5 Ids reserved are set to shard_0, " +
					"with percentages `[10, 90]` , then this would result in shard_0=100 (95 distributed + 5 reserved), shard_1=900. " +
					"**Determinism:** Optional use in conjunction with `seed` for reproducible shard assignments across Terraform runs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.AtLeast(0)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_sizes")),
					validate.ListInt64SumEquals(100),
					validate.ListRequiredWhenEquals("strategy", "percentage"),
				},
			},
			"shard_sizes": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of absolute shard sizes (exact number of IDs per shard). Required for `size` strategy. Mutually exclusive with `shard_count` and `shard_percentages`. " +
					"Array length determines number of shards (e.g., `[50, 200, -1]` creates 3 shards). Values must be positive integers or -1 (which means 'all remaining IDs'). " +
					"Only the last element can be -1. Sizes are allocated in order - first shard gets first size, second shard gets second size, etc. " +
					"When `reserved_ids` are used, reserved IDs count toward the specified shard size - the algorithm distributes fewer unreserved IDs to reach the exact target, then adds reserved IDs. " +
					"If total IDs are fewer than requested sizes, earlier shards get their full allocation and later shards receive remaining IDs. " +
					"**Example:** `[50, 200, -1]` with 1000 IDs creates shard_0=50, shard_1=200, shard_2=750 (remainder). " +
					"**Example:** `[10, 30, -1]` for controlled pilot expansion: 10 internal testers, 30 early adopters, remainder for general deployment. " +
					"**Example:** `[100, 100, 100, -1]` for fixed-size rings with equal initial waves and unlimited final wave. " +
					"**Example:** `[25, 25, 25, 25]` for exactly 4 waves of 25 devices each (no remainder wave). " +
					"**With Reserved IDs:** Optionally use in conjunction with `reserved_ids` to set specific IDs to specific shards. Scenario: if 1000 total IDs with 5 IDs reserved are set to shard_0, " +
					"with sizes `[50, -1]`, then this would result in shard_0=50 (45 distributed + 5 reserved), shard_1=950. " +
					"**Determinism:** Optional use in conjunction with `seed` for reproducible shard assignments across Terraform runs. " +
					"**Use this when:** You need exact capacity constraints (e.g., support team handles exactly 50 pilot devices, beta program limited to 100 users), " +
					"or when deployment waves must have specific sizes regardless of total device count.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.Any(
						int64validator.AtLeast(1),
						int64validator.OneOf(-1),
					)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_percentages")),
					validate.ListRequiredWhenEquals("strategy", "size"),
				},
			},
			"strategy": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The distribution strategy for sharding IDs. " +
					"`round-robin` distributes in circular order (guarantees equal sizes ±1, optional seed for reproducibility). " +
					"`percentage` distributes by specified percentages (guarantees exact percentages, requires `shard_percentages`, optional seed for reproducibility). " +
					"`size` distributes by absolute sizes (guarantees exact sizes, requires `shard_sizes`, optional seed for reproducibility). " +
					"`rendezvous` uses Highest Random Weight (HRW) algorithm for stable consistent hashing (always deterministic, requires seed). " +
					"**When to use `rendezvous`:** Choose this strategy when you have pre-existing active deployments with IDs already mapped to shards/groups and need to add new shards later. " +
					"Example: 3 existing computer groups with policies deployed, now adding a 4th group. Rendezvous keeps ~75% of IDs in their current shards (only ~25% redistribute to the new shard), " +
					"while other strategies would shuffle nearly 100% of IDs across all groups, disrupting active deployments. " +
					"**Trade-off:** Prioritizes stability over balanced distribution. Expected variance by dataset size: " +
					"100 IDs ≈12%, 500 IDs ≈15%, 1000 IDs ≈10%, 5000 IDs ≈6%, 10000+ IDs ≈1%. " +
					"When using `reserved_ids`, variance may increase moderately (e.g., 6% → 9% with significant reservations), but impact is typically minor. ",
				Validators: []validator.String{
					stringvalidator.OneOf("round-robin", "percentage", "size", "rendezvous"),
				},
			},
			"seed": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional seed value for deterministic ID distribution and can be used with all shardingstrategies. " +
					"Makes results reproducible across Terraform runs and enables controlled variance across different rollout scenarios. " +
					"**Without seed (`round-robin`, `percentage`, `size` strategies):** IDs are distributed in API response order. Since API order may vary between calls (server-side changes, pagination, etc.), " +
					"shard assignments can drift over time, causing IDs to unexpectedly move between shards on subsequent `terraform plan` operations. " +
					"**With seed (`round-robin`, `percentage`, `size` strategies):** IDs are sorted numerically, then shuffled deterministically using the seed before distribution. " +
					"Same seed always produces identical shard assignments, ensuring stability across Terraform runs. " +
					"**`rendezvous` strategy:** Always deterministic regardless of seed (uses consistent hashing). Seed affects which shard each ID maps to via Highest Random Weight algorithm. " +
					"**Multiple rollout use case:** Use different seeds for different rollout types to distribute user pilot burden fairly across devices. " +
					"Example: With seed=\"os-updates\", device ID \"12345\" lands in shard_0. With seed=\"app-deployments\", same device lands in shard_2. " +
					"This prevents the same devices from always being early adopters for every change type (OS updates, app rollouts, policy changes, etc.). " +
					"**Recommendation:** Always use seeds for production deployments to avoid unexpected shard membership changes triggering unintended policy/group reassignments.",
			},
			"exclude_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "Optional list of IDs to **completely exclude** from all shards. These IDs are filtered out before any sharding strategy is " +
					"applied and will not appear in any shard output whatsoever. " +
					"The sharding process will then proceed as if these IDs never existed in the source. Works with all source types (computers, mobile devices, users). " +
					"**Common use cases:** " +
					"Computers: Remove C-suite laptops requiring manual change management, exclude production servers. " +
					"Mobile devices: Exclude executive iPhones/iPads from beta iOS rollouts. " +
					"Users: Exclude service accounts, external consultants, or admin users from user-based policies. " +
					"**Example:** `exclude_ids = [\"1001\", \"1002\", \"1003\"]` - these three IDs will be completely absent from all shards. ",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^\d+$`),
							"must be a valid numeric ID",
						),
					),
				},
			},
			"reserved_ids": schema.MapAttribute{
				ElementType: types.ListType{ElemType: types.StringType},
				Optional:    true,
				MarkdownDescription: "Optional map of shard names to lists of IDs that should always be assigned to specific shards. " +
					"These IDs are removed from the main pool before sharding, then directly assigned to their designated shards after sharding completes. " +
					"Works with all source types (computers, mobile devices, users). " +
					"**Shard names:** Use `shard_0`, `shard_1`, `shard_2`, etc. (must match actual shard count). " +
					"**Processing order:** " +
					"1. `exclude_ids` are removed completely. " +
					"2. `reserved_ids` are extracted and set aside. " +
					"3. Remaining IDs are distributed using the selected strategy. " +
					"4. Reserved IDs are added to their designated shards. " +
					"**Distribution impact:** " +
					"`round-robin`, `percentage`, `size`: Adjusts targets to maintain balanced distribution (e.g., 74 IDs with 4 reserved across 3 shards → targets of 25, 25, 24 are maintained). " +
					"`rendezvous`: Does NOT adjust for balance. Adds reserved IDs on top of natural hash-based distribution (may increase variance by ±5-10%). " +
					"**Common use cases:** " +
					"Computers: Assign C-suite laptops to the final deployment ring (manual approval), place test devices in first ring. " +
					"Mobile devices: Assign executive iPhones/iPads to last ring for conservative iOS rollouts, IT team devices to first ring. " +
					"Users: Assign IT admins to first ring for policy testing, VIP users to last ring for stability. " +
					"**Example:** `reserved_ids = { \"shard_0\" = [\"101\", \"102\"], \"shard_2\" = [\"201\", \"202\"] }` - IDs 101/102 always in first shard, 201/202 always in third shard. " +
					"**Conflicts:** If an ID appears in both `exclude_ids` and `reserved_ids`, exclusion takes precedence (ID is completely removed). " +
					"If an ID appears in multiple shards within `reserved_ids`, validation will fail.",
				Validators: []validator.Map{
					mapvalidator.KeysAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^shard_\d+$`),
							"must be a valid shard name (e.g., 'shard_0', 'shard_1')",
						),
					),
					mapvalidator.ValueListsAre(
						listvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d+$`),
								"must be a valid numeric ID",
							),
						),
					),
				},
			},
			"shards": schema.MapAttribute{
				ElementType: types.ListType{ElemType: types.StringType},
				Computed:    true,
				MarkdownDescription: "Computed map of shard names (`shard_0`, `shard_1`, ...) to lists of IDs. " +
					"Each value is a `list(string)` type that preserves the deterministic numerical order of IDs. " +
					"Compatible with resource attributes expecting ID lists " +
					"(e.g., static group members, policy scope). " +
					"Access with `data.example.shards[\"shard_0\"]`, check size with `length(data.example.shards[\"shard_0\"])`.",
			},
		},
	}
}
