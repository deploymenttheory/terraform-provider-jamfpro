# ==============================================================================
# Test 05: Reserved IDs - Assign Specific IDs to Specific Shards
#
# Purpose: Prove that reserved_ids correctly assigns specific IDs to
# designated shards while remaining IDs are distributed normally
#
# Test Design:
# - Reserve specific IDs for shard_0 (first ring)
# - Reserve specific IDs for shard_2 (last ring)
# - Distribute remaining IDs across 3 shards with round-robin
# - Verify reserved IDs are present in correct shards
#
# Note: Update reserved IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "reserved_test" {
  source_type = "computer_inventory"
  strategy    = "round-robin"
  shard_count = 3
  seed        = "reserved-test-2026"

  # Update these IDs to match actual IDs from your inventory
  reserved_ids = {
    "shard_0" = ["1", "2"]    # IT test devices -> first ring
    "shard_2" = ["99", "100"] # Executive devices -> last ring
  }
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "reserved_shard_0" {
  description = "First ring - should contain reserved IDs plus distributed IDs"
  value       = data.jamfpro_guid_list_sharder.reserved_test.shards["shard_0"]
}

output "reserved_shard_1" {
  description = "Middle ring - distributed devices only"
  value       = data.jamfpro_guid_list_sharder.reserved_test.shards["shard_1"]
}

output "reserved_shard_2" {
  description = "Last ring - should contain reserved IDs plus distributed IDs"
  value       = data.jamfpro_guid_list_sharder.reserved_test.shards["shard_2"]
}

output "reserved_shard_0_count" {
  description = "Count of IDs in shard_0 (includes 2 reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_0"])
}

output "reserved_shard_1_count" {
  description = "Count of IDs in shard_1 (no reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_1"])
}

output "reserved_shard_2_count" {
  description = "Count of IDs in shard_2 (includes 2 reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_2"])
}

output "reserved_ids_in_shard_0" {
  description = "Verify reserved IDs are in shard_0 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_0"], "1"),
    contains(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_0"], "2")
  ])
}

output "reserved_ids_in_shard_2" {
  description = "Verify reserved IDs are in shard_2 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_2"], "99"),
    contains(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_2"], "100")
  ])
}

output "reserved_total_count" {
  description = "Total IDs distributed (should be original count)"
  value       = length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.reserved_test.shards["shard_2"])
}
