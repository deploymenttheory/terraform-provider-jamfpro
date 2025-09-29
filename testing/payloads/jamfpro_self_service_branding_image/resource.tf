resource "jamfpro_self_service_branding_image" "local_example" {
  self_service_branding_image_file_path = "${path.module}/cat.png"
}
