
run "apply_computer_extension_attributes" {
  command = apply

  module {
    source = "./testing/computer_extension_attributes"
  }
}