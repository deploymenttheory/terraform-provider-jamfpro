resource "jamfpro_icon" "icon_from_local" {
  icon_file_base64 = filebase64("${path.module}/icon.png")
}
