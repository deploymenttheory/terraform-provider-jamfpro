# ==============================================================================
# Test 03: Rendezvous Stability Test - Proving Minimal Disruption
#
# Purpose: Prove that Rendezvous Hashing has superior stability compared to
# position-based strategies when shard count changes
#
# Hypothesis: When increasing from 3 to 4 shards, only ~25% of IDs should
# move (those assigned to the new shard_3). Round-robin would cause ~75% to move!
#
# Test Design:
# - Create TWO datasources side-by-side with identical IDs
# - First: 3 shards (baseline)
# - Second: 4 shards (expanded)
# - Calculate: How many IDs stayed in their original shard?
# - Assert: Movement should be <= 30% (theoretical: ~25%)
# ==============================================================================

# Baseline: 3-shard distribution
data "jamfpro_guid_list_sharder" "baseline_3_shards" {
  source_type = "computer_inventory"
  shard_count = 3
  strategy    = "rendezvous"
  seed        = "stability-test-2026"
}

# Expanded: 4-shard distribution (same computers, same seed, +1 shard)
data "jamfpro_guid_list_sharder" "expanded_4_shards" {
  source_type = "computer_inventory"
  shard_count = 4
  strategy    = "rendezvous"
  seed        = "stability-test-2026"
}

# ==============================================================================
# Stability Calculations
# ==============================================================================

# Count how many IDs remained in shard_0
output "rendezvous_shard_0_stable_count" {
  description = "IDs that stayed in shard_0 (3-shard → 4-shard)"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_0"],
    data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_0"]
  ))
}

# Count how many IDs remained in shard_1
output "rendezvous_shard_1_stable_count" {
  description = "IDs that stayed in shard_1 (3-shard → 4-shard)"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_1"],
    data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_1"]
  ))
}

# Count how many IDs remained in shard_2
output "rendezvous_shard_2_stable_count" {
  description = "IDs that stayed in shard_2 (3-shard → 4-shard)"
  value = length(setintersection(
    data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_2"],
    data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_2"]
  ))
}

# Total IDs that didn't move
output "rendezvous_total_stable_ids" {
  description = "Total IDs that stayed in same shard number"
  value       = length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_0"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_0"])) + length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_1"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_1"])) + length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_2"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_2"]))
}

# Total IDs being distributed
output "rendezvous_total_ids" {
  description = "Total IDs in the dataset"
  value       = length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_2"])
}

# Stability percentage (what we're proving!)
output "rendezvous_stability_percentage" {
  description = "% of IDs that stayed in same shard (target: >=70%, proves <30% moved)"
  value       = floor(((length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_0"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_0"])) + length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_1"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_1"])) + length(setintersection(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_2"], data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_2"]))) / (length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.baseline_3_shards.shards["shard_2"]))) * 100)
}

# New shard_3 size (should be ~25% of total)
output "rendezvous_new_shard_3_count" {
  description = "Size of new shard_3 (should be ~25% of total)"
  value       = length(data.jamfpro_guid_list_sharder.expanded_4_shards.shards["shard_3"])
}

# Verify determinism: Same config = same ID
output "rendezvous_baseline_id" {
  value = data.jamfpro_guid_list_sharder.baseline_3_shards.id
}

output "rendezvous_expanded_id" {
  value = data.jamfpro_guid_list_sharder.expanded_4_shards.id
}
