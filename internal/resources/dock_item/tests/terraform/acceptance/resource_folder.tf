resource "jamfpro_dock_item_framework" "test_acc_folder" {
  name = "Acceptance Test - Folder Dock Item"
  type = "Folder"
  path = "~/Downloads"
  
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}