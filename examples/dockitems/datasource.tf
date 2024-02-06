data "jamfpro_dock_items" "dock_item_001_data" {
  id = jamfpro_dock_items.dock_item_001.id
}

output "jamfpro_dock_item_001_id" {
  value = data.jamfpro_dock_items.dock_item_001_data.id
}

output "jamfpro_dock_item_001_name" {
  value = data.jamfpro_dock_items.dock_item_001_data.name
}

data "jamfpro_dock_items" "dock_item_002_data" {
  id = jamfpro_dock_items.dock_item_002.id
}

output "jamfpro_dock_item_002_id" {
  value = data.jamfpro_dock_items.dock_item_002_data.id
}

output "jamfpro_dock_item_002_name" {
  value = data.jamfpro_dock_items.dock_item_002_data.name
}