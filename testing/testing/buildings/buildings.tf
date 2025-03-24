// ========================================================================== //
// Buildings
// ========================================================================== //


resource "jamfpro_building" "building" {
  name = "tf-testing-local-bw"
}

resource "jamfpro_building" "building_multiple" {
  count = 100
  name = "tf-testing-local-bw-${count.index}"
}
