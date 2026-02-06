# ==============================================================================
# Test 07: Rendezvous with Exclude IDs - Stability with Exclusions
#
# Purpose: Prove that rendezvous hashing maintains stability even when
# IDs are excluded from the distribution
#
# Test Design:
# - Baseline: Rendezvous without exclusions (3 shards)
# - Test: Rendezvous with specific IDs excluded (3 shards, same seed)
# - Verify: Excluded IDs absent, non-excluded IDs remain stable
# - Verify: Rendezvous stability maintained despite exclusions
#
# Note: Update excluded IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

# Baseline: No exclusions
data "jamfpro_guid_list_sharder" "rendezvous_baseline" {
  source_type = "computer_inventory"
  strategy    = "rendezvous"
  shard_count = 3
  seed        = "rendezvous-exclude-2026"
}

# Test: With exclusions
data "jamfpro_guid_list_sharder" "rendezvous_with_exclusions" {
  source_type = "computer_inventory"
  strategy    = "rendezvous"
  shard_count = 3
  seed        = "rendezvous-exclude-2026"
  
  # Update these IDs to match actual IDs from your inventory
  exclude_ids = ["50", "51", "52"]
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "rendezvous_exclude_baseline_total" {
  description = "Total IDs without exclusions"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_2"])
}

output "rendezvous_exclude_with_exclusions_total" {
  description = "Total IDs with 3 exclusions (should be baseline - 3)"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"])
}

output "rendezvous_exclude_ids_not_in_any_shard" {
  description = "Verify excluded IDs completely absent (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"], "50"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"], "50"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"], "50"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"], "51"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"], "51"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"], "51"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"], "52"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"], "52"),
    !contains(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"], "52")
  ])
}

# Verify that non-excluded IDs maintain their shard assignments
# This proves rendezvous stability is preserved with exclusions
output "rendezvous_exclude_shard_0_stable" {
  description = "IDs that remained in shard_0 after exclusions"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_0"],
    data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"]
  ))
}

output "rendezvous_exclude_shard_1_stable" {
  description = "IDs that remained in shard_1 after exclusions"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_1"],
    data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"]
  ))
}

output "rendezvous_exclude_shard_2_stable" {
  description = "IDs that remained in shard_2 after exclusions"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_2"],
    data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"]
  ))
}

output "rendezvous_exclude_stability_percentage" {
  description = "% of non-excluded IDs that stayed in same shard (should be ~100%)"
  value = floor((
    (length(setintersection(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_0"], data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"])) +
     length(setintersection(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_1"], data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"])) +
     length(setintersection(data.jamfpro_guid_list_sharder.rendezvous_baseline.shards["shard_2"], data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"]))) /
    (length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_0"]) + 
     length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_1"]) + 
     length(data.jamfpro_guid_list_sharder.rendezvous_with_exclusions.shards["shard_2"]))
  ) * 100)
}
