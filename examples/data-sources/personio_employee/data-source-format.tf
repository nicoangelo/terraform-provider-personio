data "personio_employee" "example_with_format" {
  id = 12345 # The Personio employee ID to load. Fails if it does not exist

  format {
    attribute = "dynamic_987654" # the dynamic attribute key to format

    phonenumber = {
      default_region = "AT"
      format = "INTERNATIONAL"
    }
  }
}
