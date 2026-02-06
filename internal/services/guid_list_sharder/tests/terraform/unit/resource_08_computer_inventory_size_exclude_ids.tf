# ==============================================================================
# Test 12: Size Strategy with Exclude IDs
#
# Purpose: Verify size-based distribution correctly excludes specific IDs
#
# Test Design:
# - Exclude specific IDs from distribution
# - Distribute remaining IDs across 3 shards with sizes (20, 30, -1)
# - Verify excluded IDs are not present in any shard
# - Verify exact sizes are maintained for first two shards
# - Verify last shard gets all remaining non-excluded IDs
#
# Note: Update excluded IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "size_exclude" {
  source_type = "computer_inventory"
  strategy    = "size"
  shard_sizes = [20, 30, -1]
  seed        = "size-exclude-2026"

  # Update these IDs to match actual IDs from your inventory
  exclude_ids = ["5", "15", "25", "35"]
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "size_exclude_shard_0_count" {
  description = "Shard 0 count (should be exactly 20)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"])
}

output "size_exclude_shard_1_count" {
  description = "Shard 1 count (should be exactly 30)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"])
}

output "size_exclude_shard_2_count" {
  description = "Shard 2 count (all remaining non-excluded IDs)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"])
}

output "size_exclude_total" {
  description = "Total IDs distributed (should be inventory - 4 excluded)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"])
}

output "size_exclude_shard_0_meets_target" {
  description = "Verify shard_0 has exactly 20 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"]) == 20
}

output "size_exclude_shard_1_meets_target" {
  description = "Verify shard_1 has exactly 30 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"]) == 30
}

output "size_exclude_ids_absent" {
  description = "Verify excluded IDs not in any shard (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"], "5"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"], "5"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"], "5"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"], "15"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"], "15"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"], "15"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"], "25"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"], "25"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"], "25"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_0"], "35"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_1"], "35"),
    !contains(data.jamfpro_guid_list_sharder.size_exclude.shards["shard_2"], "35")
  ])
}
