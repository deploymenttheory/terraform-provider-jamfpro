package guid_list_sharder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setStateToTerraform sets all values in the state object
func setStateToTerraform(ctx context.Context, state *GuidListSharderDataSourceModel, shards [][]string, sourceType string, shardCount int, strategy string) error {

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

	// Generate deterministic ID based on configuration
	// This ensures the datasource ID remains stable across refreshes
	idString := fmt.Sprintf("%s-%d-%s-%s-%s",
		sourceType,
		shardCount,
		strategy,
		state.Seed.ValueString(),
		state.GroupId.ValueString(),
	)
	hash := sha256.Sum256([]byte(idString))
	state.Id = types.StringValue(hex.EncodeToString(hash[:]))
	state.Shards = shardsMapValue

	return nil
}
