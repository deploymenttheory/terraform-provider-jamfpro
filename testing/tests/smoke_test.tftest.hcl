

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

run "apply_computer_extension_attributes" {
  command = apply
  
  module {
    source = "./testing/computer_extension_attributes"
  }

}