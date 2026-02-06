# ==============================================================================
# Test 13: Size Strategy with Reserved IDs
#
# Purpose: Verify size-based distribution correctly handles reserved IDs
# while maintaining exact target sizes
#
# Test Design:
# - Reserve specific IDs for shard_0 and shard_1
# - Distribute remaining IDs to meet exact shard sizes (25, 40, -1)
# - Verify reserved IDs are in correct shards
# - Verify shard sizes match targets exactly (reserved + distributed = target)
#
# Note: Update reserved IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "size_reserved" {
  source_type = "computer_inventory"
  strategy    = "size"
  shard_sizes = [25, 40, -1]
  seed        = "size-reserved-2026"

  # Update these IDs to match actual IDs from your inventory
  reserved_ids = {
    "shard_0" = ["1", "2", "3", "4", "5"] # 5 reserved for dev/test
    "shard_1" = ["10", "11"]              # 2 reserved for staging
  }
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "size_reserved_shard_0_count" {
  description = "Shard 0 count (should be exactly 25: 5 reserved + 20 distributed)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"])
}

output "size_reserved_shard_1_count" {
  description = "Shard 1 count (should be exactly 40: 2 reserved + 38 distributed)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_1"])
}

output "size_reserved_shard_2_count" {
  description = "Shard 2 count (all remaining)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"])
}

output "size_reserved_total" {
  description = "Total IDs distributed (should equal inventory count)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"])
}

output "size_reserved_shard_0_meets_target" {
  description = "Verify shard_0 has exactly 25 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"]) == 25
}

output "size_reserved_shard_1_meets_target" {
  description = "Verify shard_1 has exactly 40 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_1"]) == 40
}

output "size_reserved_ids_in_shard_0" {
  description = "Verify reserved IDs are in shard_0 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"], "1"),
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"], "2"),
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"], "3"),
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"], "4"),
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_0"], "5")
  ])
}

output "size_reserved_ids_in_shard_1" {
  description = "Verify reserved IDs are in shard_1 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_1"], "10"),
    contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_1"], "11")
  ])
}

output "size_reserved_shard_2_no_reserved" {
  description = "Verify shard_2 has no reserved IDs (should be true)"
  value = alltrue([
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "1"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "2"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "3"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "4"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "5"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "10"),
    !contains(data.jamfpro_guid_list_sharder.size_reserved.shards["shard_2"], "11")
  ])
}
