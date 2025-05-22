run "apply_departments" {
  command = apply

  module {
    source = "./testing/departments"
  }
}