# Local variable to create a flattened list of all user-to-group assignments
# This makes it easier to iterate and create individual binding resources.
locals {
  user_group_assignments = flatten([
    # Iterate over each user configuration provided
    for user_config in var.user_to_add_to_access_group : [
      # For each user, iterate over the list of group names they should be assigned to
      for group_name in toset(user_config.groups) : {
        # Create a unique key for the for_each in the binding resource
        assignment_key = "user-${user_config.username}-group-${group_name}"
        user_username  = user_config.username
        group_name     = group_name
      }
    ]
  ])
}