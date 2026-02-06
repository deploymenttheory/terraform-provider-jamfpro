# ==============================================================================
# Test 08: Rendezvous with Reserved IDs - Stability with Reservations
#
# Purpose: Prove that rendezvous hashing works correctly with reserved IDs
# while maintaining its core stability property
#
# Important: Rendezvous does NOT guarantee balanced distribution with reserved IDs.
# It prioritizes stability (minimal churn on shard count changes) over balance.
# Each ID's shard assignment is determined by hash weights, not target quotas.
#
# Test Design:
# - Reserve specific IDs for shard_0 (first ring - IT test devices)
# - Reserve specific IDs for shard_2 (last ring - executive devices)
# - Distribute remaining IDs across 3 shards with rendezvous (natural distribution)
# - Verify reserved IDs are in correct shards
# - Accept whatever distribution results (balance may vary from target)
#
# Note: Update reserved IDs to match actual IDs in your Jamf Pro inventory
# ==============================================================================

data "jamfpro_guid_list_sharder" "rendezvous_reserved_test" {
  source_type = "computer_inventory"
  strategy    = "rendezvous"
  shard_count = 3
  seed        = "rendezvous-reserved-2026"
  
  # Update these IDs to match actual IDs from your inventory
  reserved_ids = {
    "shard_0" = ["1", "2"]      # IT test devices -> first ring
    "shard_2" = ["99", "100"]   # Executive devices -> last ring
  }
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "rendezvous_reserved_shard_0" {
  description = "First ring - should contain reserved IDs plus distributed IDs"
  value       = data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"]
}

output "rendezvous_reserved_shard_1" {
  description = "Middle ring - distributed devices only (no reserved IDs)"
  value       = data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_1"]
}

output "rendezvous_reserved_shard_2" {
  description = "Last ring - should contain reserved IDs plus distributed IDs"
  value       = data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"]
}

output "rendezvous_reserved_shard_0_count" {
  description = "Count of IDs in shard_0 (includes 2 reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"])
}

output "rendezvous_reserved_shard_1_count" {
  description = "Count of IDs in shard_1 (no reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_1"])
}

output "rendezvous_reserved_shard_2_count" {
  description = "Count of IDs in shard_2 (includes 2 reserved IDs)"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"])
}

output "rendezvous_reserved_ids_in_shard_0" {
  description = "Verify reserved IDs are in shard_0 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"], "1"),
    contains(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"], "2")
  ])
}

output "rendezvous_reserved_ids_in_shard_2" {
  description = "Verify reserved IDs are in shard_2 (should be true)"
  value = alltrue([
    contains(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"], "99"),
    contains(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"], "100")
  ])
}

output "rendezvous_reserved_total_count" {
  description = "Total IDs distributed (should equal original inventory count)"
  value       = length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"])
}

output "rendezvous_reserved_distribution_variance" {
  description = "Distribution variance (max - min shard size). Note: Rendezvous prioritizes stability over balance, so higher variance is expected."
  value = max(
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"]),
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_1"]),
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"])
  ) - min(
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_0"]),
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_1"]),
    length(data.jamfpro_guid_list_sharder.rendezvous_reserved_test.shards["shard_2"])
  )
}
