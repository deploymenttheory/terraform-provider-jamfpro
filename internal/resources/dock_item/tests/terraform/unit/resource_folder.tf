resource "jamfpro_dock_item_framework" "test_folder" {
  name = "Test Folder Dock Item"
  type = "Folder"
  path = "~/Downloads"
  
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "2m"
  }
}