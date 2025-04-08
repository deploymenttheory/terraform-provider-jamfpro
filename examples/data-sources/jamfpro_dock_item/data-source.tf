# Test data source by ID
data "jamfpro_dock_item" "test_by_id" {
  id = jamfpro_dock_item.test.id
}

# Test data source by name 
data "jamfpro_dock_item" "test_by_name" {
 name = jamfpro_dock_item.test.name
}

# Outputs
output "dock_item_by_id" {
  value = {
    id   = data.jamfpro_dock_item.test_by_id.id
    name = data.jamfpro_dock_item.test_by_id.name
    type = data.jamfpro_dock_item.test_by_id.type
    path = data.jamfpro_dock_item.test_by_id.path
  }
}

output "dock_item_by_name" {
  value = {
    id   = data.jamfpro_dock_item.test_by_name.id
    name = data.jamfpro_dock_item.test_by_name.name
    type = data.jamfpro_dock_item.test_by_name.type
    path = data.jamfpro_dock_item.test_by_name.path
  }
}

output "dock_item_by_name" {
 value = {
   id   = data.jamfpro_dock_item.test_by_name.id 
   name = data.jamfpro_dock_item.test_by_name.name
   type = data.jamfpro_dock_item.test_by_name.type
   path = data.jamfpro_dock_item.test_by_name.path
 }
}