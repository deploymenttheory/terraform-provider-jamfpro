resource "jamfpro_computer_inventory_collection" "example" {
  local_user_accounts               = true
  home_directory_sizes              = true
  hidden_accounts                   = true
  printers                          = true
  active_services                   = true
  mobile_device_app_purchasing_info = true
  computer_location_information     = true
  package_receipts                  = true
  available_software_updates        = true
  include_applications              = true
  include_fonts                     = true
  include_plugins                   = true

  applications {
    path     = "/Applications/ExampleApp.app"
    platform = "macOS"
  }

  applications {
    path     = "/Applications/AnotherApp.app"
    platform = "macOS"
  }

  fonts {
    path     = "/Library/Fonts/ExampleFont.ttf"
    platform = "macOS"
  }

  fonts {
    path     = "/Library/Fonts/AnotherFont.ttf"
    platform = "macOS"
  }

  plugins {
    path     = "/Library/Internet Plug-Ins/ExamplePlugin.plugin"
    platform = "macOS"
  }

  plugins {
    path     = "/Library/Internet Plug-Ins/AnotherPlugin.plugin"
    platform = "macOS"
  }
}
