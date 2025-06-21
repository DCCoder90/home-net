#!/bin/bash

# This script generates a dynamic inventory for Ansible by querying the Terraform state.

set -e

cd "$(dirname "$0")/../../terraform"

terraform output -json ansible_inventory