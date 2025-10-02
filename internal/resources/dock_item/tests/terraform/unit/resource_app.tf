resource "jamfpro_dock_item_framework" "test_app" {
  name = "Test App Dock Item"
  type = "App"
  path = "file://localhost/Applications/iTunes.app/"
  
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "2m"
  }
}