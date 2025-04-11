run "apply_buildings" {
  command = apply

  module {
    source = "./testing/buildings"
  }
}

run "apply_departments" {
  command = apply

  module {
    source = "./testing/departments"
  }
}

run "apply_categories" {
  command = apply

  module {
    source = "./testing/categories"
  }
}

run "apply_computer_extension_attributes" {
  command = apply

  module {
    source = "./testing/computer_extension_attributes"
  }
}

run "apply_scripts" {
  command = apply

  module {
    source = "./testing/scripts"
  }
}

run "apply_static_computer_groups" {
  command = apply

  module {
    source = "./testing/static_computer_groups"
  }
}