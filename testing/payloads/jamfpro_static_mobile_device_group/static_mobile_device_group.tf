// ========================================================================== //
// Static Mobile Device Groups
// ========================================================================== //

data "local_file" "site_and_mobile_device_ids" {
  filename = "testing/data_sources/site_and_mobile_device_ids.json"
}

resource "jamfpro_static_mobile_device_group" "name" {
  name                       = "tf-testing-${var.testing_id}-script-max-${random_id.rng.hex}"
  site_id                    = jsondecode(data.local_file.site_and_mobile_device_ids.content).site
  assigned_mobile_device_ids = jsondecode(data.local_file.site_and_mobile_device_ids.content).mobile_devices
}

