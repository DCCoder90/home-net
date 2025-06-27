locals {
  # Create a map of secret names to their generated password values.
  # This is the ideal format for a for_each loop, ensuring each resource
  # instance has a unique, string-based key.
  generated_secrets_map = {
    for name, password_obj in random_password.generated_secret : name => password_obj.result
  }

  # Create a list of objects from the map.
  # This is useful for outputs that require a list format.
  generated_secret_list = [
    for name, value in local.generated_secrets_map : {
      name  = name
      value = value
    }
  ]
}