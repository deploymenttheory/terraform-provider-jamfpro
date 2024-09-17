resource "jamfpro_network_segment" "jamfpro_network_segment_001" {
  name                 = "Example Network Segment"
  starting_address     = "192.168.1.1"
  ending_address       = "192.168.1.254"
  distribution_server  = "Example Distribution Server"
  distribution_point   = "Example Distribution Point"
  url                  = "http://example.com"
  swu_server           = "Example SWU Server"
  building             = "Main Building"
  department           = "IT Department"
  override_buildings   = true
  override_departments = false
}
