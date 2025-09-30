resource "jamfpro_self_service_branding_image" "local_example" {
  self_service_branding_image_file_path = "/path/to/branding_image.png"
}

resource "jamfpro_self_service_branding_image" "web_example" {
  self_service_branding_image_file_web_source = "https://example.com/branding_image.png"
}
