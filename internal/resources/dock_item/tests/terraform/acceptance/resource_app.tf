resource "jamfpro_dock_item_framework" "test_acc_app" {
  name = "Acceptance Test - App Dock Item"
  type = "App"
  path = "file://localhost/Applications/iTunes.app/"
  
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}