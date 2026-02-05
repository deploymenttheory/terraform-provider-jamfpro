package guid_list_sharder

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GuidListSharderDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	SourceType       types.String `tfsdk:"source_type"`
	GroupId          types.String `tfsdk:"group_id"`
	ShardCount       types.Int64  `tfsdk:"shard_count"`
	ShardPercentages types.List   `tfsdk:"shard_percentages"`
	ShardSizes       types.List   `tfsdk:"shard_sizes"`
	Strategy         types.String `tfsdk:"strategy"`
	Seed             types.String `tfsdk:"seed"`
	Shards           types.Map    `tfsdk:"shards"`
}
