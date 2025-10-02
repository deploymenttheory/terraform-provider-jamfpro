resource "jamfpro_dock_item_framework" "test_acc_file" {
  name = "Acceptance Test - File Dock Item"
  type = "File"
  path = "/etc/hosts"
  
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}