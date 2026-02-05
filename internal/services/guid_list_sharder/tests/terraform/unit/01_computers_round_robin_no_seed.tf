# ==============================================================================
# Test 01: Computers - Round-Robin Strategy (No Seed)
#
# Purpose: Verify round-robin distribution produces exactly equal shard sizes
# using API order (non-deterministic between runs)
#
# Use Case: Quick one-time equal split where reproducibility isn't needed
#
# Expected Behavior:
# - Exactly equal shard sizes (within ±1)
# - Uses API order (may change between Terraform runs)
# - Fast processing (no shuffle overhead)
# ==============================================================================

data "jamfpro_guid_list_sharder" "test" {
  source_type = "computer_inventory"
  shard_count = 4
  strategy    = "round-robin"
  # No seed - uses API order (non-deterministic)
}

output "shard_0_count" {
  description = "Computers in shard 0 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Computers in shard 1 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Computers in shard 2 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_2"])
}

output "shard_3_count" {
  description = "Computers in shard 3 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_3"])
}

output "size_variance" {
  description = "Max difference between largest and smallest shard"
  value       = max(length(data.jamfpro_guid_list_sharder.test.shards["shard_0"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_1"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_2"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_3"])) - min(length(data.jamfpro_guid_list_sharder.test.shards["shard_0"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_1"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_2"]), length(data.jamfpro_guid_list_sharder.test.shards["shard_3"]))
}
