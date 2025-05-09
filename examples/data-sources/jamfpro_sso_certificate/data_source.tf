data "jamfpro_sso_certificate" "current" {}

output "certificate_basic_info" {
  value = {
    type              = data.jamfpro_sso_certificate.current.keystore[0].type
    keystore_filename = data.jamfpro_sso_certificate.current.keystore[0].keystore_file_name
    setup_type        = data.jamfpro_sso_certificate.current.keystore[0].keystore_setup_type
  }
}

output "certificate_keys" {
  value = data.jamfpro_sso_certificate.current.keystore[0].keys
}

output "certificate_details" {
  value = {
    issuer        = data.jamfpro_sso_certificate.current.keystore_details[0].issuer
    subject       = data.jamfpro_sso_certificate.current.keystore_details[0].subject
    expiration    = data.jamfpro_sso_certificate.current.keystore_details[0].expiration
    serial_number = data.jamfpro_sso_certificate.current.keystore_details[0].serial_number
  }
}
