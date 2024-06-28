
// Script example an uploaded script taken from a file path with parameters
resource "jamfpro_script" "scripts_0001" {
  name            = "tf-example-script-fileupload"
  script_contents = file("support_files/scripts/Add or Remove Group Membership.zsh")
  category_id = 5
  os_requirements = "13"
  priority        = "BEFORE"
  info            = "Adds target user or group to specified group membership, or removes said membership."
  notes           = "Jamf Pro script parameters 4 -> 7"
  parameter4  = "100" // targetID
  parameter5  = "group" // Target Type - Must be either "user" or "group"
  parameter6  = "someGroupName" // targetMembership
  parameter7  = "add" // Script Action - Must be either "add" or "remove"
}

// Script example with an inline script
resource "jamfpro_script" "scripts_0002" {
  name            = "tf-example-script-inline"
  script_contents = "hello world"
  os_requirements = "13.1"
  priority        = "BEFORE"
  info = "Your script info here."
  notes = ""

}



