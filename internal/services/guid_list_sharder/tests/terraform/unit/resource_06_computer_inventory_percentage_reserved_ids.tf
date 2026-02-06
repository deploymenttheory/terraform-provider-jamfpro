# ==============================================================================
# Test 10: Percentage Strategy with Reserved IDs
#
# Purpose: Verify percentage-based distribution correctly handles reserved IDs
# while maintaining target percentage distribution
#
# Test Design:
# - Reserve specific IDs for shard_0 and shard_2
# - Distribute remaining IDs across 3 shards with percentages (20%, 30%, 50%)
# - Verify reserved IDs are in correct shards
# - Verify total shard sizes still match percentage targets
#
# Note: Update reserved IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "percentage_reserved" {
  source_type       = "computer_inventory"
  strategy          = "percentage"
  shard_percentages = [20, 30, 50]
  seed              = "percentage-reserved-2026"

  # Update these IDs to match actual IDs from your inventory
  reserved_ids = {
    "shard_0" = ["1", "2", "3"]     # Early adopters
    "shard_2" = ["98", "99", "100"] # Production hold-backs
  }
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "percentage_reserved_shard_0_count" {
  description = "Shard 0 count (should be ~20% of total)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_0"])
}

output "percentage_reserved_shard_1_count" {
  description = "Shard 1 count (should be ~30% of total)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"])
}

output "percentage_reserved_shard_2_count" {
  description = "Shard 2 count (should be ~50% of total)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_2"])
}

output "percentage_reserved_total" {
  description = "Total IDs distributed (should equal inventory count)"
  value       = length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_2"])
}

output "percentage_reserved_ids_in_shard_0" {
  description = "Verify reserved IDs are in shard_0 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_0"], "1"),
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_0"], "2"),
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_0"], "3")
  ])
}

output "percentage_reserved_ids_in_shard_2" {
  description = "Verify reserved IDs are in shard_2 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_2"], "98"),
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_2"], "99"),
    contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_2"], "100")
  ])
}

output "percentage_reserved_shard_1_no_reserved" {
  description = "Verify shard_1 has no reserved IDs (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "1"),
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "2"),
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "3"),
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "98"),
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "99"),
    !contains(data.jamfpro_guid_list_sharder.percentage_reserved.shards["shard_1"], "100")
  ])
}
