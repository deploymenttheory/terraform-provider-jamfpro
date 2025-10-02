resource "jamfpro_dock_item_framework" "test_file" {
  name = "Test File Dock Item"
  type = "File"
  path = "/etc/hosts"
  
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "2m"
  }
}