# ==============================================================================
# Test 06: Exclude IDs - Remove Specific IDs from All Shards
#
# Purpose: Prove that exclude_ids completely removes specific IDs from
# the sharding process and they don't appear in any shard output
#
# Test Design:
# - Exclude specific IDs from distribution
# - Distribute remaining IDs across 3 shards with round-robin
# - Verify excluded IDs are not present in any shard
# - Verify total count is reduced by number of excluded IDs
#
# Note: Update excluded IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

# Baseline: No exclusions
data "jamfpro_guid_list_sharder" "baseline_no_exclusions" {
  source_type = "computer_inventory"
  strategy    = "round-robin"
  shard_count = 3
  seed        = "exclude-test-2026"
}

# Test: With exclusions
data "jamfpro_guid_list_sharder" "with_exclusions" {
  source_type = "computer_inventory"
  strategy    = "round-robin"
  shard_count = 3
  seed        = "exclude-test-2026"
  
  # Update these IDs to match actual IDs from your inventory
  exclude_ids = ["50", "51", "52"]
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "exclude_baseline_total" {
  description = "Total IDs without exclusions"
  value       = length(data.jamfpro_guid_list_sharder.baseline_no_exclusions.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.baseline_no_exclusions.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.baseline_no_exclusions.shards["shard_2"])
}

output "exclude_with_exclusions_total" {
  description = "Total IDs with 3 exclusions (should be baseline - 3)"
  value       = length(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"])
}

output "exclude_ids_not_in_shard_0" {
  description = "Verify excluded IDs not in shard_0 (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "52")
  ])
}

output "exclude_ids_not_in_shard_1" {
  description = "Verify excluded IDs not in shard_1 (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "52")
  ])
}

output "exclude_ids_not_in_shard_2" {
  description = "Verify excluded IDs not in shard_2 (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "52")
  ])
}

output "exclude_ids_not_in_any_shard" {
  description = "Verify excluded IDs completely absent (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "50"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "51"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_0"], "52"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_1"], "52"),
    !contains(data.jamfpro_guid_list_sharder.with_exclusions.shards["shard_2"], "52")
  ])
}
