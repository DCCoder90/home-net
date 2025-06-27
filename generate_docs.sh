#!/bin/bash

# Define the root Terraform directory
TERRAFORM_ROOT="/home/dccoder/home-net/terraform"

# Define the path to the central terraform-docs configuration file
TERRAFORM_DOCS_CONFIG="${TERRAFORM_ROOT}/.terraform-docs.yml"

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
    echo "Processing: ${dir}"

    # Check if the directory contains any .tf files
    if find "${dir}" -maxdepth 1 -type f -name "*.tf" | grep -q .; then
        # Run terraform-docs from within the directory
        # This ensures the output file is created in the correct location
        # and relative paths within the module are handled correctly.
        # We explicitly override the output file name from .terraform-docs.yml
        (cd "${dir}" && terraform-docs . --config "${TERRAFORM_DOCS_CONFIG}" --output-file "${OUTPUT_FILENAME}")
        if [ $? -ne 0 ]; then
            echo "Error: terraform-docs failed for directory: ${dir}"
            return 1
        fi
    else
        echo "Skipping empty directory (no .tf files found): ${dir}"
    fi
    return 0
}

# Process the root Terraform directory first
if ! process_dir "${TERRAFORM_ROOT}"; then
    exit 1
fi

# Process all module directories
# Using -print0 and read -r -d $'\0' for robust handling of spaces in directory names
find "${TERRAFORM_ROOT}/modules" -type d -print0 | while IFS= read -r -d $'\0' dir; do
    if ! process_dir "${dir}"; then
        exit 1
    fi
done

echo "${OUTPUT_FILENAME} generation complete for all specified directories."