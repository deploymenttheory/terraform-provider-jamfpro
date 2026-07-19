// ========================================================================== //
// Printers
// ========================================================================== //

resource "jamfpro_printer" "jamfpro_printer_min" {
  name = "tf-testing-${var.testing_id}-min-${random_id.rng.hex}"
}

// Regression test for issue #1145 - heredoc strings in HCL always include a
// trailing newline before EOT, but the API strips it server-side. Without a
// DiffSuppressFunc on info/notes, this produced perpetual drift on every plan
// after apply.
resource "jamfpro_printer" "jamfpro_printer_heredoc" {
  name          = "tf-testing-${var.testing_id}-heredoc-${random_id.rng.hex}"
  category_name = "No category assigned"
  info          = <<-EOT
    Multi-line info field used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
  notes         = <<-EOT
    Multi-line notes field used to verify no drift is
    reported after apply due to the heredoc trailing newline.
  EOT
}
