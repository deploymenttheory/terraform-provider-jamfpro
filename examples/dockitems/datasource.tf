data "jamfpro_dock_item" "example_dock_item" {
  id = "12345"  # Replace this with the actual ID of the dock item you want to retrieve
}

output "dock_item_id" {
  value = data.jamfpro_dock_item.example_dock_item.id
}

output "dock_item_name" {
  value = data.jamfpro_dock_item.example_dock_item.name
}

output "dock_item_type" {
  value = data.jamfpro_dock_item.example_dock_item.type
}

output "dock_item_path" {
  value = data.jamfpro_dock_item.example_dock_item.path
}

output "dock_item_contents" {
  value = data.jamfpro_dock_item.example_dock_item.contents
}
