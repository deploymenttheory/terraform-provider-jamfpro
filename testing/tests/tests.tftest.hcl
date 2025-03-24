

run "apply_buildings" {
  command = apply
  
  module {
    source = "./testing/buildings"
  }
  
}

run "apply_scripts" {
  command = apply
  
  module {
    source = "./testing/scripts"
  }
  
}
