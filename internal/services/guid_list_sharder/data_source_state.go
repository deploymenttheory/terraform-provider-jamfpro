package guid_list_sharder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func setStateToTerraform(ctx context.Context, state *GuidListSharderDataSourceModel, shards [][]string) error {
	shardsMap := make(map[string]types.List, len(shards))
	for i, shard := range shards {
		shardList, diags := types.ListValueFrom(ctx, types.StringType, shard)
		if diags.HasError() {
			return fmt.Errorf("failed to convert shard %d to list: %v", i, diags.Errors())
		}
		shardsMap[fmt.Sprintf("shard_%d", i)] = shardList
	}

	shardsMapValue, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, shardsMap)
	if diags.HasError() {
		return fmt.Errorf("failed to convert shards map to state: %v", diags.Errors())
	}
	shardCount := len(shards)
	idString := fmt.Sprintf("%s-%d-%s-%s-%s",
		state.SourceType.ValueString(),
		shardCount,
		state.Strategy.ValueString(),
		state.Seed.ValueString(),
		state.GroupId.ValueString(),
	)
	hash := sha256.Sum256([]byte(idString))
	state.Id = types.StringValue(hex.EncodeToString(hash[:]))
	state.Shards = shardsMapValue

	return nil
}
