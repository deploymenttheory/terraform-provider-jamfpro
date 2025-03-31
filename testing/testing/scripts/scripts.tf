// ========================================================================== //
// Scripts
// ========================================================================== //


// ========================================================================== //
// Single scripts

resource "jamfpro_script" "min_script" {
  name            = "tf-testing-min-${random_id.rng.hex}"
  script_contents = "script_contents_field"
  priority        = "BEFORE"
}

resource "jamfpro_script" "max_script" {
  name            = "tf-testing-bw-max-${random_id.rng.hex}"
  category_id     = "9"
  info            = "info_field"
  notes           = "notes_field"
  os_requirements = "os_requirements_field"
  priority        = "BEFORE"
  script_contents = "script_contents_field"
  parameter4      = "parameter4_field"
  parameter5      = "parameter5_field"
  parameter6      = "parameter6_field"
  parameter7      = "parameter7_field"
  parameter8      = "parameter8_field"
  parameter9      = "parameter9_field"
  parameter10     = "parameter10_field"
  parameter11     = "parameter11_field"
}
// ========================================================================== //
// Multiple scripts

resource "jamfpro_script" "multiple_script_min" {
  count           = 100
  name            = "tf-testing-min-${count.index}-${random_id.rng.hex}"
  script_contents = "echo hello world"
  priority        = "BEFORE"
}

resource "jamfpro_script" "multiple_script_max" {
  count           = 100
  name            = "tf-testing-max-${count.index}-${random_id.rng.hex}"
  category_id     = "9"
  info            = "info_field"
  notes           = "notes_field"
  os_requirements = "os_requirements_field"
  priority        = "BEFORE"
  script_contents = "script_contents_field"
  parameter4      = "parameter4_field"
  parameter5      = "parameter5_field"
  parameter6      = "parameter6_field"
  parameter7      = "parameter7_field"
  parameter8      = "parameter8_field"
  parameter9      = "parameter9_field"
  parameter10     = "parameter10_field"
  parameter11     = "parameter11_field"
}