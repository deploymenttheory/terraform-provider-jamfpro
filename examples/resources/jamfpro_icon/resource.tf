# Example using a local file
resource "jamfpro_icon" "icon_from_local" {
  icon_file_path = "/Users/dafyddwatkins/localtesting/terraform/support_files/icons/firefox_logo_icon_170152.png"
}

# Example using a web source
resource "jamfpro_icon" "icon_from_web" {
  icon_file_web_source = "https://upload.wikimedia.org/wikipedia/commons/1/16/Firefox_logo%2C_2017.png"
}