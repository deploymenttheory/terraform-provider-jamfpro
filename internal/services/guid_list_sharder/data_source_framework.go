package guid_list_sharder

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &guidListSharderDataSource{}
	_ datasource.DataSourceWithConfigure      = &guidListSharderDataSource{}
	_ datasource.DataSourceWithValidateConfig = &guidListSharderDataSource{}
)

func NewGuidListSharderDataSource() datasource.DataSource {
	return &guidListSharderDataSource{}
}

type guidListSharderDataSource struct {
	client *jamfpro.Client
}

type guidListSharderModel struct {
	ID               types.String `tfsdk:"id"`
	SourceType       types.String `tfsdk:"source_type"`
	GroupID          types.String `tfsdk:"group_id"`
	ExcludeGroupID   types.String `tfsdk:"exclude_group_id"`
	ReserveGroupID0  types.String `tfsdk:"reserve_group_id_shard_0"`
	Strategy         types.String `tfsdk:"strategy"`
	ShardCount       types.Int64  `tfsdk:"shard_count"`
	ShardPercentages types.List   `tfsdk:"shard_percentages"`
	ShardSizes       types.List   `tfsdk:"shard_sizes"`
	Seed             types.String `tfsdk:"seed"`
	ExcludeIDs       types.List   `tfsdk:"exclude_ids"`
	ReservedIDs      types.Map    `tfsdk:"reserved_ids"`

	Shards   types.Map `tfsdk:"shards"`
	Metadata types.Map `tfsdk:"metadata"`
}

func (d *guidListSharderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guid_list_sharder"
}

func (d *guidListSharderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *guidListSharderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Deterministically shards Jamf Pro IDs into rollout waves (shard_0..shard_N) for progressive deployments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier for this data source instance.",
			},
			"source_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID source: `computer_inventory`, `mobile_device_inventory`, `computer_group_membership`, `mobile_device_group_membership`, `user_accounts`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"computer_inventory",
						"mobile_device_inventory",
						"computer_group_membership",
						"mobile_device_group_membership",
						"user_accounts",
					),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group ID required for *_group_membership sources.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(numericIDRe, "group_id must be a numeric ID"),
				},
			},
			"exclude_group_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Convenience: exclude IDs that are members of this group. For computer sources, treated as a computer group ID; for mobile sources, treated as a mobile device group ID.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(numericIDRe, "exclude_group_id must be a numeric ID"),
				},
			},
			"reserve_group_id_shard_0": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Convenience: reserve all IDs in this group into shard_0 (first shard). For computer sources, treated as a computer group ID; for mobile sources, treated as a mobile device group ID.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(numericIDRe, "reserve_group_id_shard_0 must be a numeric ID"),
				},
			},
			"strategy": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Sharding strategy: `round-robin`, `percentage`, `size`, `rendezvous`.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(StrategyRoundRobin), string(StrategyPercentage), string(StrategySize), string(StrategyRendezvous)),
				},
			},
			"shard_count": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of shards (required for `round-robin` and `rendezvous`).",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"shard_percentages": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "Percentages (sum to 100). Required for `percentage` strategy.",
			},
			"shard_sizes": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "Absolute sizes. Required for `size` strategy. Use `-1` as final element for remainder.",
			},
			"seed": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Seed for deterministic distribution.",
			},
			"exclude_ids": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "IDs to exclude from all shards.",
			},
			"reserved_ids": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.ListType{ElemType: types.StringType},
				MarkdownDescription: "Map of shard name -> list of IDs pinned to that shard.",
			},
			"shards": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.ListType{ElemType: types.StringType},
				MarkdownDescription: "Computed shards: `shard_N` -> list of IDs.",
			},
			"metadata": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Computed metadata (counts, seed, strategy, etc.).",
			},
		},
	}
}

func (d *guidListSharderDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data guidListSharderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceType := data.SourceType.ValueString()
	groupID := data.GroupID.ValueString()
	excludeGroupID := data.ExcludeGroupID.ValueString()
	reserveGroupID0 := data.ReserveGroupID0.ValueString()
	strategy := data.Strategy.ValueString()

	shardCount := int(data.ShardCount.ValueInt64())
	shardPercentages := expandInt64List(ctx, data.ShardPercentages)
	shardSizes := expandInt64List(ctx, data.ShardSizes)
	excludeIDs := expandStringListFramework(ctx, data.ExcludeIDs)
	reservedIDs := expandReservedIDsFramework(ctx, data.ReservedIDs)

	errMsg := validateConfigInputs(sourceType, groupID, excludeGroupID, reserveGroupID0, strategy, shardCount, shardPercentages, shardSizes, excludeIDs, reservedIDs)

	if errMsg != "" {
		resp.Diagnostics.AddError("Invalid configuration", errMsg)
	}
}

func validateConfigInputs(
	sourceType string,
	groupID string,
	excludeGroupID string,
	reserveGroupID0 string,
	strategy string,
	shardCount int,
	shardPercentages []int,
	shardSizes []int,
	excludeIDs []string,
	reservedIDs map[string][]string,
) string {
	errs := make([]string, 0, 8)

	errs = append(errs, validateSourceInputs(sourceType, groupID, excludeGroupID, reserveGroupID0)...)
	errs = append(errs, validateShardInputs(strategy, shardCount, shardPercentages, shardSizes)...)
	errs = append(errs, validateIDs(excludeIDs, reservedIDs)...)
	errs = append(errs, validateExcludeReserveConflict(excludeIDs, reservedIDs)...)

	if len(errs) == 0 {
		return ""
	}
	return joinErrors(errs)
}

func validateSourceInputs(sourceType, groupID, excludeGroupID, reserveGroupID0 string) []string {
	var errs []string
	if (sourceType == "computer_group_membership" || sourceType == "mobile_device_group_membership") && groupID == "" {
		errs = append(errs, fmt.Sprintf("group_id is required when source_type is %q", sourceType))
	}

	if excludeGroupID != "" {
		if _, err := groupMembershipSourceTypeFor(sourceType); err != nil {
			errs = append(errs, fmt.Sprintf("exclude_group_id cannot be used with source_type %q", sourceType))
		}
	}

	if reserveGroupID0 != "" {
		if _, err := groupMembershipSourceTypeFor(sourceType); err != nil {
			errs = append(errs, fmt.Sprintf("reserve_group_id_shard_0 cannot be used with source_type %q", sourceType))
		}
	}

	return errs
}

func validateShardInputs(strategy string, shardCount int, shardPercentages []int, shardSizes []int) []string {
	var errs []string

	set := 0
	if shardCount > 0 {
		set++
	}
	if len(shardPercentages) > 0 {
		set++
	}
	if len(shardSizes) > 0 {
		set++
	}
	if set != 1 {
		errs = append(errs, "exactly one of shard_count, shard_percentages, or shard_sizes must be set")
	}

	switch Strategy(strategy) {
	case StrategyRoundRobin, StrategyRendezvous:
		if shardCount < 1 {
			errs = append(errs, fmt.Sprintf("strategy %q requires shard_count >= 1", strategy))
		}
	case StrategyPercentage:
		errs = append(errs, validatePercentageInputs(shardPercentages)...)
	case StrategySize:
		errs = append(errs, validateSizeInputs(shardSizes)...)
	}

	return errs
}

func validatePercentageInputs(shardPercentages []int) []string {
	if len(shardPercentages) == 0 {
		return []string{"strategy 'percentage' requires shard_percentages"}
	}

	sum := 0
	for _, p := range shardPercentages {
		if p < 0 {
			return []string{"shard_percentages values must be >= 0"}
		}
		sum += p
	}
	if sum != 100 {
		return []string{"shard_percentages must sum to 100"}
	}

	return nil
}

func validateSizeInputs(shardSizes []int) []string {
	if len(shardSizes) == 0 {
		return []string{"strategy 'size' requires shard_sizes"}
	}

	var errs []string
	for i, s := range shardSizes {
		if s != -1 && s < 1 {
			errs = append(errs, fmt.Sprintf("shard_sizes[%d] must be >= 1 or -1", i))
		}
		if s == -1 && i != len(shardSizes)-1 {
			errs = append(errs, fmt.Sprintf("shard_sizes[%d] is -1 but only final element may be -1", i))
		}
	}

	return errs
}

func validateIDs(excludeIDs []string, reservedIDs map[string][]string) []string {
	var errs []string

	for _, id := range excludeIDs {
		if !numericIDRe.MatchString(id) {
			errs = append(errs, fmt.Sprintf("exclude_ids contains non-numeric ID %q", id))
		}
	}

	for shardName, ids := range reservedIDs {
		if !shardNameRe.MatchString(shardName) {
			errs = append(errs, fmt.Sprintf("reserved_ids key %q must match shard_N", shardName))
		}
		for _, id := range ids {
			if !numericIDRe.MatchString(id) {
				errs = append(errs, fmt.Sprintf("reserved_ids[%q] contains non-numeric ID %q", shardName, id))
			}
		}
	}

	return errs
}

func validateExcludeReserveConflict(excludeIDs []string, reservedIDs map[string][]string) []string {
	if len(excludeIDs) == 0 || len(reservedIDs) == 0 {
		return nil
	}

	excludeSet := make(map[string]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excludeSet[id] = true
	}

	var errs []string
	for shardName, ids := range reservedIDs {
		for _, id := range ids {
			if excludeSet[id] {
				errs = append(errs, fmt.Sprintf("ID %q appears in both exclude_ids and reserved_ids[%q]", id, shardName))
			}
		}
	}
	return errs
}

func (d *guidListSharderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data guidListSharderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceType := data.SourceType.ValueString()
	groupID := data.GroupID.ValueString()
	excludeGroupID := data.ExcludeGroupID.ValueString()
	reserveGroupID0 := data.ReserveGroupID0.ValueString()
	strategy := data.Strategy.ValueString()
	seed := data.Seed.ValueString()

	shardCount := int(data.ShardCount.ValueInt64())
	shardPercentages := expandInt64List(ctx, data.ShardPercentages)
	shardSizes := expandInt64List(ctx, data.ShardSizes)
	excludeIDs := expandStringListFramework(ctx, data.ExcludeIDs)
	reservedIDs := expandReservedIDsFramework(ctx, data.ReservedIDs)

	ids, err := fetchSourceIDsFramework(ctx, d.client, sourceType, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch IDs", err.Error())
		return
	}

	excludedGroupCount := 0
	if excludeGroupID != "" {
		groupSourceType, err := groupMembershipSourceTypeFor(sourceType)
		if err != nil {
			resp.Diagnostics.AddError("Invalid exclude_group_id", err.Error())
			return
		}
		excludeGroupIDs, err := fetchSourceIDsFramework(ctx, d.client, groupSourceType, excludeGroupID)
		if err != nil {
			resp.Diagnostics.AddError("Failed to fetch exclude_group_id membership", err.Error())
			return
		}
		excludedGroupCount = len(excludeGroupIDs)
		excludeIDs = mergeUniqueStrings(excludeIDs, excludeGroupIDs)
	}

	reservedGroupCount := 0
	if reserveGroupID0 != "" {
		groupSourceType, err := groupMembershipSourceTypeFor(sourceType)
		if err != nil {
			resp.Diagnostics.AddError("Invalid reserve_group_id_shard_0", err.Error())
			return
		}
		reserveGroupIDs, err := fetchSourceIDsFramework(ctx, d.client, groupSourceType, reserveGroupID0)
		if err != nil {
			resp.Diagnostics.AddError("Failed to fetch reserve_group_id_shard_0 membership", err.Error())
			return
		}
		reservedGroupCount = len(reserveGroupIDs)
		reservedIDs["shard_0"] = mergeUniqueStrings(reservedIDs["shard_0"], reserveGroupIDs)
	}

	totalFetched := len(ids)
	filtered := ApplyExclusions(ids, excludeIDs)
	excludedCount := totalFetched - len(filtered)

	// Runtime conflict: IDs can't be both excluded and reserved.
	if len(excludeIDs) > 0 && len(reservedIDs) > 0 {
		excludeSet := make(map[string]bool, len(excludeIDs))
		for _, id := range excludeIDs {
			excludeSet[id] = true
		}
		for shardName, ids := range reservedIDs {
			for _, id := range ids {
				if excludeSet[id] {
					resp.Diagnostics.AddError(
						"Invalid configuration",
						fmt.Sprintf("ID %q appears in both exclude list and reserved_ids[%q]", id, shardName),
					)
					return
				}
			}
		}
	}

	resolvedShardCount := resolveShardCountFramework(shardCount, shardPercentages, shardSizes)
	reservations, err := ApplyReservations(filtered, reservedIDs, resolvedShardCount)
	if err != nil {
		resp.Diagnostics.AddError("Failed to apply reservations", err.Error())
		return
	}
	reservedCount := len(filtered) - len(reservations.UnreservedIDs)

	shards, err := Shard(filtered, Strategy(strategy), resolvedShardCount, shardPercentages, shardSizes, seed, reservations)
	if err != nil {
		resp.Diagnostics.AddError("Failed to shard IDs", err.Error())
		return
	}

	shardsMap := make(map[string]types.List, len(shards))
	for i, shard := range shards {
		vals := make([]types.String, 0, len(shard))
		for _, id := range shard {
			vals = append(vals, types.StringValue(id))
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, vals)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		shardsMap[fmt.Sprintf("shard_%d", i)] = listVal
	}

	shardsTF, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, shardsMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	generatedAt := time.Now().UTC().Format(time.RFC3339)
	metadataMap := map[string]types.String{
		"generated_at":               types.StringValue(generatedAt),
		"source_type":                types.StringValue(sourceType),
		"group_id":                   types.StringValue(groupID),
		"exclude_group_id":           types.StringValue(excludeGroupID),
		"reserve_group_id_shard_0":   types.StringValue(reserveGroupID0),
		"strategy":                   types.StringValue(strategy),
		"seed":                       types.StringValue(seed),
		"total_ids_fetched":          types.StringValue(strconv.Itoa(totalFetched)),
		"excluded_id_count":          types.StringValue(strconv.Itoa(excludedCount)),
		"excluded_group_id_count":    types.StringValue(strconv.Itoa(excludedGroupCount)),
		"reserved_id_count":          types.StringValue(strconv.Itoa(reservedCount)),
		"reserved_group_id_count":    types.StringValue(strconv.Itoa(reservedGroupCount)),
		"unreserved_ids_distributed": types.StringValue(strconv.Itoa(len(reservations.UnreservedIDs))),
		"shard_count":                types.StringValue(strconv.Itoa(len(shards))),
	}
	metadataTF, diags := types.MapValueFrom(ctx, types.StringType, metadataMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s:%s:%s:%s", sourceType, groupID, strategy, seed))
	data.Shards = shardsTF
	data.Metadata = metadataTF

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fetchSourceIDsFramework(ctx context.Context, client *jamfpro.Client, sourceType, groupID string) ([]string, error) {
	switch sourceType {
	case "computer_inventory":
		params := url.Values{}
		params.Set("section", "GENERAL")
		resp, err := client.GetComputersInventory(params)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve computer inventory: %w", err)
		}
		ids := make([]string, 0, len(resp.Results))
		for _, c := range resp.Results {
			if c.General.RemoteManagement.Managed {
				ids = append(ids, c.ID)
			}
		}
		return ids, nil

	case "mobile_device_inventory":
		resp, err := client.GetMobileDevices()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve mobile devices: %w", err)
		}
		ids := make([]string, 0, len(resp.MobileDevices))
		for _, md := range resp.MobileDevices {
			if md.Managed {
				ids = append(ids, strconv.Itoa(md.ID))
			}
		}
		return ids, nil

	case "computer_group_membership":
		group, err := client.GetComputerGroupByID(groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve computer group %s: %w", groupID, err)
		}
		if group.Computers == nil {
			return []string{}, nil
		}
		ids := make([]string, 0, len(*group.Computers))
		for _, c := range *group.Computers {
			ids = append(ids, strconv.Itoa(c.ID))
		}
		return ids, nil

	case "mobile_device_group_membership":
		group, err := client.GetMobileDeviceGroupByID(groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve mobile device group %s: %w", groupID, err)
		}
		if group.MobileDevices == nil {
			return []string{}, nil
		}
		ids := make([]string, 0, len(*group.MobileDevices))
		for _, md := range *group.MobileDevices {
			ids = append(ids, strconv.Itoa(md.ID))
		}
		return ids, nil

	case "user_accounts":
		resp, err := client.GetUsers()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve users: %w", err)
		}
		ids := make([]string, 0, len(resp.Users))
		for _, u := range resp.Users {
			ids = append(ids, strconv.Itoa(u.ID))
		}
		return ids, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSourceType, sourceType)
	}
}

func expandInt64List(ctx context.Context, v types.List) []int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var raw []int64
	_ = v.ElementsAs(ctx, &raw, false)
	out := make([]int, 0, len(raw))
	for _, n := range raw {
		out = append(out, int(n))
	}
	return out
}

func expandStringListFramework(ctx context.Context, v types.List) []string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var raw []string
	_ = v.ElementsAs(ctx, &raw, false)
	return raw
}

func expandReservedIDsFramework(ctx context.Context, v types.Map) map[string][]string {
	out := make(map[string][]string)
	if v.IsNull() || v.IsUnknown() {
		return out
	}
	var raw map[string][]string
	_ = v.ElementsAs(ctx, &raw, false)
	if raw == nil {
		return out
	}
	return raw
}

func joinErrors(errs []string) string {
	return strings.Join(errs, "; ")
}

func resolveShardCountFramework(shardCount int, shardPercentages []int, shardSizes []int) int {
	if len(shardPercentages) > 0 {
		return len(shardPercentages)
	}
	if len(shardSizes) > 0 {
		return len(shardSizes)
	}
	return shardCount
}

func groupMembershipSourceTypeFor(sourceType string) (string, error) {
	switch sourceType {
	case "computer_inventory", "computer_group_membership":
		return "computer_group_membership", nil
	case "mobile_device_inventory", "mobile_device_group_membership":
		return "mobile_device_group_membership", nil
	default:
		return "", fmt.Errorf("%w: %q", ErrUnsupportedSourceTypeForGroupConvenience, sourceType)
	}
}

func mergeUniqueStrings(a []string, b []string) []string {
	if len(b) == 0 {
		return a
	}
	seen := make(map[string]bool, len(a)+len(b))
	out := make([]string, 0, len(a)+len(b))
	for _, v := range a {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	for _, v := range b {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
