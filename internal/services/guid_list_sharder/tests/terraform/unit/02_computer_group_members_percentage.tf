# ==============================================================================
# Test 02: Computer Group Members - Percentage Strategy
#
# Purpose: Verify percentage-based distribution from a computer group
#
# Use Case: Split existing group membership into phased rollout shards
# (10% pilot, 30% broader, 60% full deployment)
#
# Expected Behavior:
# - Accurate percentage distribution
# - All group members accounted for
# - Last shard gets any remainder
# ==============================================================================

data "jamfpro_guid_list_sharder" "test" {
  source_type       = "computer_group_membership"
  group_id          = "123" # Replace with actual group ID
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "deployment-2024"
}

output "pilot_count" {
  description = "Pilot computers (10%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_0"])
}

output "broader_count" {
  description = "Broader rollout computers (30%)"
  value       = length(data.jamfpro_guid_list_sharder.test.shards["shard_1"])
}

output "full_count" {
  description = "Full deployment computers (60%)"
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

resource "jamfpro_static_computer_group" "broader_group" {
  name    = "Broader Deployment - 30%"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.test.shards["shard_1"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "full_group" {
  name    = "Full Deployment - 60%"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.test.shards["shard_2"] : tonumber(id)
  ]
}
