#!/bin/bash

# Define the Terraform directories
TERRAFORM_CORE="/home/dccoder/home-net/terraform/core"
TERRAFORM_APPS="/home/dccoder/home-net/terraform/apps"

# Define the desired output filename for the generated documentation
OUTPUT_FILENAME="README.md"

# Check if terraform-docs is installed
if ! command -v terraform-docs &> /dev/null
then
    echo "Error: terraform-docs is not installed. Please install it first."
    echo "Refer to: https://terraform-docs.io/user-guide/installation/"
    exit 1
fi

echo "Generating ${OUTPUT_FILENAME} for Terraform directories..."

# Function to process a single directory
process_dir() {
    local dir="$1"
    local config="${dir}/.terraform-docs.yml"
    echo "Processing: ${dir}"

    # Check if the directory contains any .tf files
    if find "${dir}" -maxdepth 1 -type f -name "*.tf" | grep -q .; then
        # Run terraform-docs from within the directory
        # This ensures the output file is created in the correct location
        # and relative paths within the module are handled correctly.
        # We explicitly override the output file name from .terraform-docs.yml
        (cd "${dir}" && terraform-docs . --config "${config}" --output-file "${OUTPUT_FILENAME}")
        if [ $? -ne 0 ]; then
            echo "Error: terraform-docs failed for directory: ${dir}"
            return 1
        fi
    else
        echo "Skipping empty directory (no .tf files found): ${dir}"
    fi
    return 0
}

# Process core Terraform directory
if ! process_dir "${TERRAFORM_CORE}"; then
    exit 1
fi

# Process apps Terraform directory
if ! process_dir "${TERRAFORM_APPS}"; then
    exit 1
fi

echo "${OUTPUT_FILENAME} generation complete for all specified directories."
