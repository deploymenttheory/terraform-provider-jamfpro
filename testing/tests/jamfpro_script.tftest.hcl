run "apply_scripts" {
  command = apply

  module {
    source = "./testing/scripts"
  }
}