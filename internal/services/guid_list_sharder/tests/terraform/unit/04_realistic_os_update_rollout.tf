# ==============================================================================
# Test 04: Realistic OS Update Rollout
#
# Purpose: Demonstrate a real-world use case - phased macOS update deployment
# with pilot, early adopters, and production groups
#
# Scenario: IT team wants to roll out macOS Sonoma to 500 computers
# - Test Ring: 25 computers (5%)
# - Early Adopters: 75 computers (15%)
# - Production: 400 computers (80%)
#
# ==============================================================================

# Query all macOS computers and shard them into update rings
data "jamfpro_guid_list_sharder" "macos_update_rings" {
  source_type       = "computer_inventory"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "macos-sonoma-rollout-2024"
}

# Create static groups for each update ring
resource "jamfpro_static_computer_group" "test_ring" {
  name    = "macOS Sonoma - Test Ring"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_0"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "early_adopters" {
  name    = "macOS Sonoma - Early Adopters"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_1"] : tonumber(id)
  ]
}

resource "jamfpro_static_computer_group" "production" {
  name    = "macOS Sonoma - Production"
  site_id = -1
  assigned_computer_ids = [
    for id in data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_2"] : tonumber(id)
  ]
}

# Create policies for each ring (example - simplified)
# Note: You would create actual jamfpro_policy resources here
# with different deployment schedules for each group

# ==============================================================================
# Outputs for monitoring rollout progress
# ==============================================================================

output "test_ring_count" {
  description = "Number of computers in test ring"
  value       = length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_0"])
}

output "early_adopters_count" {
  description = "Number of computers in early adopters ring"
  value       = length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_1"])
}

output "production_count" {
  description = "Number of computers in production ring"
  value       = length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_2"])
}

output "total_computers" {
  description = "Total computers in rollout"
  value       = length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_0"]) + length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_1"]) + length(data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_2"])
}

output "test_ring_ids" {
  description = "Computer IDs in test ring (for verification)"
  value       = data.jamfpro_guid_list_sharder.macos_update_rings.shards["shard_0"]
}
