resource "jamfpro_device_enrollments" "example" {
  name          = "Test Device Enrollment"
  encoded_token = filebase64("/path/to/device_enrollment_token.p7m")

  # Optional fields
  supervision_identity_id = "-1"
  site_id                 = "-1"
}
