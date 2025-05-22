run "apply_buildings" {
  command = apply

  module {
    source = "./testing/buildings"
  }
}