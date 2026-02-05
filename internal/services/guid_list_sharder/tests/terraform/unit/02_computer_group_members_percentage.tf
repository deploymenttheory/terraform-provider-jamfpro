# ==============================================================================
# Test 02: Computer Group Members - Percentage Strategy
#
# Purpose: Verify percentage-based distribution from a computer group
#
# Use Case: Split existing group membership into phased rollout shards
# (10% pilot, 30% limited, 60% broad deployment)
#
# Expected Behavior:
# - Accurate percentage distribution
# - All group members accounted for
# - Last shard gets any remainder
# ==============================================================================

data "jamfpro_guid_list_sharder" "test" {
  source_type       = "computer_inventory"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "deployment-2026"
}

output "pilot_count" {
  description = "01 Pilot computers (10%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_0"])
}

output "limited_count" {
  description = "02 Limited rollout computers (30%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_1"])
}

output "broad_count" {
  description = "03 Broad rollout computers (60%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_2"])
}

output "total_distributed" {
  description = "Total computers distributed"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.test.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.test.shards["shard_2"])
}

# Example usage: Create static groups for each shard
resource "jamfpro_static_computer_group" "pilot_group" {
  name    = "Pilot Deployment - 10%"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.test.shards["shard_0"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "limited_group" {
  name    = "Limited Deployment - 30%"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.test.shards["shard_1"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "broad_group" {
  name    = "Broad Deployment - 60%"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.test.shards["shard_2"] : tonumber(id)
  ]
}
