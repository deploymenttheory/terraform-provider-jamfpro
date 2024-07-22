// Specific printer ppd
resource "jamfpro_printer" "jamfpro_printers_001" {
  name          = "tf-example-printer-specific_ppd-01"
  category_name = "No category assigned"
  uri           = "lpd://10.1.20.204/"
  cups_name     = "HP_DesignJet_1050C_PS3"
  location      = "string"
  model         = "HP DesignJet 1050C PS3"
  info          = "string"
  notes         = "string"
  make_default  = true
  use_generic   = false
  ppd           = "HP_DesignJet_1050C_PS3.ppd"
  ppd_path      = "/somepath"
  ppd_contents  = file("/Users/dafyddwatkins/localtesting/support_files/printerppd/HP_DesignJet_1050C_PS3.ppd")
}

// Generic printer ppd
resource "jamfpro_printer" "jamfpro_printers_002" {
  name          = "tf-example-printer-generic_ppd-01"
  category_name = "No category assigned"
  uri           = "lpd://10.1.20.204/"
  cups_name     = "HP_9th_Floor"
  location      = "string"
  model         = "HP LaserJet 500 color MFP M575"
  info          = "string"
  notes         = "string"
  make_default  = false
  use_generic   = true
  ppd           = ""
  ppd_path      = "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/PrintCore.framework/Resources/Generic.ppd"
  ppd_contents  = ""
}
