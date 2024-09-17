resource "jamfpro_managed_software_update" "specific_version_update" {
  group {
    group_id    = "2"
    object_type = "COMPUTER"
  }

  config {
    update_action    = "DOWNLOAD_INSTALL_ALLOW_DEFERRAL"
    version_type     = "SPECIFIC_VERSION"
    specific_version = "15.1"
    max_deferrals    = 3
  }
}

resource "jamfpro_managed_software_update" "specific_version_update" {
  group {
    group_id    = "2"
    object_type = "MOBILE_DEVICE"
  }

  config {
    update_action    = "DOWNLOAD_INSTALL_ALLOW_DEFERRAL"
    version_type     = "SPECIFIC_VERSION"
    specific_version = "15.1"
    max_deferrals    = 3
  }
}