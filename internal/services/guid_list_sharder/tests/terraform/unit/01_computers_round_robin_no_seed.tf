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

data "jamfpro_guid_list_sharder" "round_robin" {
  source_type = "computer_inventory"
  shard_count = 4
  strategy    = "round-robin"
  # No seed - uses API order (non-deterministic)
}

# Create static groups for each shard
resource "jamfpro_static_computer_group" "shard_0" {
  name    = "Test Round Robin - Shard 0 - ${random_string.test_suffix.result}"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.round_robin.shards["shard_0"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "shard_1" {
  name    = "Test Round Robin - Shard 1 - ${random_string.test_suffix.result}"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.round_robin.shards["shard_1"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "shard_2" {
  name    = "Test Round Robin - Shard 2 - ${random_string.test_suffix.result}"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.round_robin.shards["shard_2"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "shard_3" {
  name    = "Test Round Robin - Shard 3 - ${random_string.test_suffix.result}"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.round_robin.shards["shard_3"] : tonumber(id)
  ]
}

output "test_suffix" {
  description = "Random suffix for round_robin group names"
  value       = random_string.test_suffix.result
}

output "round_robin_shard_0_count" {
  description = "Computers in shard 0 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_0"])
}

output "round_robin_shard_1_count" {
  description = "Computers in shard 1 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_1"])
}

output "round_robin_shard_2_count" {
  description = "Computers in shard 2 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_2"])
}

output "round_robin_shard_3_count" {
  description = "Computers in shard 3 (should be ~25%)"
  value       = length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_3"])
}

output "round_robin_size_variance" {
  description = "Max difference between largest and smallest shard"
  value       = max(length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_0"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_1"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_2"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_3"])) - min(length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_0"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_1"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_2"]), length(data.jamfpro_guid_list_sharder.round_robin.shards["shard_3"]))
}
