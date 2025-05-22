run "apply_static_computer_groups" {
  command = apply

  module {
    source = "./testing/static_computer_groups"
  }
}