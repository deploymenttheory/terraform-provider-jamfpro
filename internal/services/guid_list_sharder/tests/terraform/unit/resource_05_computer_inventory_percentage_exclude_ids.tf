# ==============================================================================
# Test 09: Percentage Strategy with Exclude IDs
#
# Purpose: Verify percentage-based distribution correctly excludes specific IDs
#
# Test Design:
# - Exclude specific IDs from distribution
# - Distribute remaining IDs across 3 shards with percentages (20%, 30%, 50%)
# - Verify excluded IDs are not present in any shard
# - Verify percentage distribution is maintained on remaining IDs
#
# Note: Update excluded IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "percentage_exclude" {
  source_type       = "computer_inventory"
  strategy          = "percentage"
  shard_percentages = [20, 30, 50]
  seed              = "percentage-exclude-2026"
  
  # Update these IDs to match actual IDs from your inventory
  exclude_ids = ["10", "20", "30"]
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "percentage_exclude_shard_0_count" {
  description = "Shard 0 count (should be ~20% of non-excluded IDs)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_0"])
}

output "percentage_exclude_shard_1_count" {
  description = "Shard 1 count (should be ~30% of non-excluded IDs)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_1"])
}

output "percentage_exclude_shard_2_count" {
  description = "Shard 2 count (should be ~50% of non-excluded IDs)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_2"])
}

output "percentage_exclude_total" {
  description = "Total IDs distributed (should be inventory - 3 excluded)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_2"])
}

output "percentage_exclude_ids_absent" {
  description = "Verify excluded IDs not in any shard (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_0"], "10"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_1"], "10"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_2"], "10"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_0"], "20"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_1"], "20"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_2"], "20"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_0"], "30"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_1"], "30"),
    !contains(data.jamfpro_guid_list_sharder.percentage_exclude.shards["shard_2"], "30")
  ])
}
