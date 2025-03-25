// ---- by name ---- //

data "jamfpro_device_enrollments" "by_name" {
  name = "Apple School Manager"
}

output "device_enrollment_id" {
  value = data.jamfpro_device_enrollments.by_name.id
}

output "device_enrollment_name" {
  value = data.jamfpro_device_enrollments.by_name.name
}

output "supervision_identity_id" {
  value = data.jamfpro_device_enrollments.by_name.supervision_identity_id
}

output "site_id" {
  value = data.jamfpro_device_enrollments.by_name.site_id
}

output "server_name" {
  value = data.jamfpro_device_enrollments.by_name.server_name
}

output "server_uuid" {
  value = data.jamfpro_device_enrollments.by_name.server_uuid
}

output "admin_id" {
  value = data.jamfpro_device_enrollments.by_name.admin_id
}

output "org_name" {
  value = data.jamfpro_device_enrollments.by_name.org_name
}

output "org_email" {
  value = data.jamfpro_device_enrollments.by_name.org_email
}

output "org_phone" {
  value = data.jamfpro_device_enrollments.by_name.org_phone
}

output "org_address" {
  value = data.jamfpro_device_enrollments.by_name.org_address
}

output "token_expiration_date" {
  value = data.jamfpro_device_enrollments.by_name.token_expiration_date
}
// ---- by id ---- //

data "jamfpro_device_enrollments" "by_id" {
  id = "3"
}

output "device_enrollment_id" {
  value = data.jamfpro_device_enrollments.by_id.id
}

output "device_enrollment_name" {
  value = data.jamfpro_device_enrollments.by_id.name
}

output "supervision_identity_id" {
  value = data.jamfpro_device_enrollments.by_id.supervision_identity_id
}

output "site_id" {
  value = data.jamfpro_device_enrollments.by_id.site_id
}

output "server_name" {
  value = data.jamfpro_device_enrollments.by_id.server_name
}

output "server_uuid" {
  value = data.jamfpro_device_enrollments.by_id.server_uuid
}

output "admin_id" {
  value = data.jamfpro_device_enrollments.by_id.admin_id
}

output "org_name" {
  value = data.jamfpro_device_enrollments.by_id.org_name
}

output "org_email" {
  value = data.jamfpro_device_enrollments.by_id.org_email
}

output "org_phone" {
  value = data.jamfpro_device_enrollments.by_id.org_phone
}

output "org_address" {
  value = data.jamfpro_device_enrollments.by_id.org_address
}

output "token_expiration_date" {
  value = data.jamfpro_device_enrollments.by_id.token_expiration_date
}
