# ==============================================================================
# Test 11: Size Strategy - Explicit Shard Sizes
#
# Purpose: Verify size-based distribution with explicit shard sizes
#
# Use Case: Create fixed-size deployment groups
# - Dev/Test: 25 computers
# - Staging: 75 computers
# - Production: All remaining computers
#
# Expected Behavior:
# - Exact sizes for defined shards
# - Last shard with -1 gets all remaining IDs
# - All IDs accounted for
# ==============================================================================

data "jamfpro_guid_list_sharder" "size_basic" {
  source_type = "computer_inventory"
  strategy    = "size"
  shard_sizes = [25, 75, -1]
  seed        = "size-basic-2026"
}

# Create static groups for each environment
resource "jamfpro_static_computer_group" "dev_test" {
  name    = "Size Test - Dev/Test (25)"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.size_basic.shards["shard_0"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "staging" {
  name    = "Size Test - Staging (75)"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.size_basic.shards["shard_1"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "production" {
  name    = "Size Test - Production (Remainder)"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.size_basic.shards["shard_2"] : tonumber(id)
  ]
}

# ==============================================================================
# Verification Outputs
# ==============================================================================

output "size_shard_0_count" {
  description = "Dev/Test count (should be exactly 25)"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_0"])
}

output "size_shard_1_count" {
  description = "Staging count (should be exactly 75)"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_1"])
}

output "size_shard_2_count" {
  description = "Production count (all remaining)"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_2"])
}

output "size_total_distributed" {
  description = "Total computers distributed"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_2"])
}

output "size_shard_0_meets_target" {
  description = "Verify shard_0 has exactly 25 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_0"]) == 25
}

output "size_shard_1_meets_target" {
  description = "Verify shard_1 has exactly 75 IDs (should be true)"
  value       = length(data.jamfpro_guid_list_sharder.size_basic.shards["shard_1"]) == 75
}
