run "apply_categories" {
  command = apply

  module {
    source = "./testing/categories"
  }
}