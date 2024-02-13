resource "jamfpro_file_share_distribution_points" "full_example" {
  // Page 1
  name                     = "tf-example-fileshare-distribution-point"
  ip_address               = "ny.company.com"
  is_master                = false
  failover_point           = "Cloud Distribution Point" // Cloud Distribution Point / or any other dp defined with jamf pro
  enable_load_balancing    = false
  // Page 2
  connection_type          = "AFP" // SMB / AFP
  share_name               = "FullExampleShareName"
  share_port               = 445
  workgroup_or_domain      = "EXAMPLEDOMAIN"
  read_only_username       = "fullreadonlyuser"
  read_only_password       = "funky_password_qwerty"
  read_write_username      = "fullreadwriteuser"
  read_write_password      = "funky_password_qwerty"
  no_authentication_required = false
  // page 3
  https_downloads_enabled   = true
  https_port                = 443
  https_share_path           = "/contextpath"
  https_username_password_required = true
  https_username            = "fullhttpuser"
  https_password            = "funky_password_qwerty"

}