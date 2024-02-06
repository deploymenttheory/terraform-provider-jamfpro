data "jamfpro_dock_item" "example_dock_item" {
  id = resource.jamfpro_dock_item.example_dock_item.id
}

output "jamfpro_dock_item_id" {
  value = data.jamfpro_dock_item.example_dock_item.id
}

output "jamfpro_dock_item_name" {
  value = data.jamfpro_dock_item.example_dock_item.name
}
