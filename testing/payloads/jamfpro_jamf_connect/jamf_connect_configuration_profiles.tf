// ========================================================================== //
// Jamf Connect Configuration Profiles
// ========================================================================== //

// ========================================================================== //
// Create configuration profile for testing

resource "jamfpro_macos_configuration_profile_plist" "jamf_connect_license_001" {
  name                = "Jamf Connect License"
  level               = "System"
  distribution_method = "Install Automatically"
  redeploy_on_update  = "Newly Assigned"
  payloads            = base64decode("PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgoJPGRpY3Q+CgkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCTxhcnJheT4KCQkJPGRpY3Q+CgkJCQk8a2V5PlBheWxvYWRDb250ZW50PC9rZXk+CgkJCQk8ZGljdD4KCQkJCQk8a2V5PmNvbS5qYW1mLmNvbm5lY3Q8L2tleT4KCQkJCQk8ZGljdD4KCQkJCQkJPGtleT5Gb3JjZWQ8L2tleT4KCQkJCQkJPGFycmF5PgoJCQkJCQkJPGRpY3Q+CgkJCQkJCQkJPGtleT5tY3hfcHJlZmVyZW5jZV9zZXR0aW5nczwva2V5PgoJCQkJCQkJCTxkaWN0PgoJCQkJCQkJCQk8a2V5PkxpY2Vuc2VGaWxlPC9rZXk+CgkJCQkJCQkJCTxzdHJpbmc+RmFrZUxpY2Vuc2VGb3JUZXN0aW5nPC9zdHJpbmc+CgkJCQkJCQkJPC9kaWN0PgoJCQkJCQkJPC9kaWN0PgoJCQkJCQk8L2FycmF5PgoJCQkJCTwvZGljdD4KCQkJCTwvZGljdD4KCQkJCTxrZXk+UGF5bG9hZERpc3BsYXlOYW1lPC9rZXk+CgkJCQk8c3RyaW5nPkN1c3RvbSBTZXR0aW5nczwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCQkJPHN0cmluZz5FMUMyRjc2Ny02NjA4LTQyNDYtQTZCQy04RTQ4OUE5MEYxRDI8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZE9yZ2FuaXphdGlvbjwva2V5PgoJCQkJPHN0cmluZz5KQU1GIFNvZnR3YXJlPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCQk8c3RyaW5nPmNvbS5hcHBsZS5NYW5hZ2VkQ2xpZW50LnByZWZlcmVuY2VzPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRVVUlEPC9rZXk+CgkJCQk8c3RyaW5nPkUxQzJGNzY3LTY2MDgtNDI0Ni1BNkJDLThFNDg5QTkwRjFEMjwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkVmVyc2lvbjwva2V5PgoJCQkJPGludGVnZXI+MTwvaW50ZWdlcj4KCQkJPC9kaWN0PgoJCQk8ZGljdD4KCQkJCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCQkJCTxkaWN0PgoJCQkJCTxrZXk+Y29tLmphbWYuY29ubmVjdC5sb2dpbjwva2V5PgoJCQkJCTxkaWN0PgoJCQkJCQk8a2V5PkZvcmNlZDwva2V5PgoJCQkJCQk8YXJyYXk+CgkJCQkJCQk8ZGljdD4KCQkJCQkJCQk8a2V5Pm1jeF9wcmVmZXJlbmNlX3NldHRpbmdzPC9rZXk+CgkJCQkJCQkJPGRpY3Q+CgkJCQkJCQkJCTxrZXk+TGljZW5zZUZpbGU8L2tleT4KCQkJCQkJCQkJPHN0cmluZz5GYWtlTGljZW5zZUZvclRlc3Rpbmc8L3N0cmluZz4KCQkJCQkJCQk8L2RpY3Q+CgkJCQkJCQk8L2RpY3Q+CgkJCQkJCTwvYXJyYXk+CgkJCQkJPC9kaWN0PgoJCQkJPC9kaWN0PgoJCQkJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCQkJCTxzdHJpbmc+Q3VzdG9tIFNldHRpbmdzPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+CgkJCQk8c3RyaW5nPkZEQTRBREZGLTUyQkItNEU0RS05OUEwLUExQzFCREI2QzA1QTwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJCQk8c3RyaW5nPkpBTUYgU29mdHdhcmU8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCQkJCTxzdHJpbmc+Y29tLmFwcGxlLk1hbmFnZWRDbGllbnQucHJlZmVyZW5jZXM8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCQkJCTxzdHJpbmc+RkRBNEFERkYtNTJCQi00RTRFLTk5QTAtQTFDMUJEQjZDMDVBPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQk8L2RpY3Q+CgkJPC9hcnJheT4KCQk8a2V5PlBheWxvYWREZXNjcmlwdGlvbjwva2V5PgoJCTxzdHJpbmcvPgoJCTxrZXk+UGF5bG9hZERpc3BsYXlOYW1lPC9rZXk+CgkJPHN0cmluZz5KYW1mIENvbm5lY3QgTGljZW5zZTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZEVuYWJsZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCTxzdHJpbmc+RjZDM0NDRDQtOTNBOC00NTNGLThDREMtMzVEMjRGRjVENTBFPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJPHN0cmluZz5ORVRPUElFPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkUmVtb3ZhbERpc2FsbG93ZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkU2NvcGU8L2tleT4KCQk8c3RyaW5nPlN5c3RlbTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCQk8c3RyaW5nPkNvbmZpZ3VyYXRpb248L3N0cmluZz4KCQk8a2V5PlBheWxvYWRVVUlEPC9rZXk+CgkJPHN0cmluZz5GNkMzQ0NENC05M0E4LTQ1M0YtOENEQy0zNUQyNEZGNUQ1MEU8L3N0cmluZz4KCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJPGludGVnZXI+MTwvaW50ZWdlcj4KCTwvZGljdD4KPC9wbGlzdD4=")
  payload_validate    = true
  user_removable      = false
  scope {
    all_computers = false
    all_jss_users = false
  }
}

// ========================================================================== //
// Data source by configuration profile ID

data "jamfpro_jamf_connect" "by_id" {
  depends_on = [
    jamfpro_macos_configuration_profile_plist.jamf_connect_license_001
  ]
  profile_id = jamfpro_macos_configuration_profile_plist.jamf_connect_license_001.id
}

// ========================================================================== //
// Data source by configuration profile name

data "jamfpro_jamf_connect" "by_name" {
  depends_on = [
    jamfpro_macos_configuration_profile_plist.jamf_connect_license_001
  ]
  profile_name = "Jamf Connect License"
}
