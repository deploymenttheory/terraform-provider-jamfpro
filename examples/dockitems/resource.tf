
// App Dock Item Example
resource "jamfpro_dock_item" "dock_item_001" {
  name     = "tf-example-dockItem-app-iTunes"
  type     = "App"
  path     = "file://localhost/Applications/iTunes.app/"
}

// File Dock Item Example
resource "jamfpro_dock_item" "dock_item_002" {
  name     = "tf-example-dockItem-file-hosts"
  type     = "File" // App / File / Folder
  path     = "/etc/hosts"
}

// Folder Dock Item Example
resource "jamfpro_dock_item" "dock_item_003" {
  name     = "tf-example-dockItem-folder-downloadsFolder"
  type     = "Folder" // App / File / Folder
  path     = "~/Downloads"
}