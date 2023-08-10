#!/bin/bash

# caputres the user's home directory
home_dir="$PWD"

# Search for XML files containing "xpath" in their filenames in the home directory
xml_files=$(find "$PWD" -type f -name '*xpath*.xml')

# Check if any matching XML files were found
if [ -n "$xml_files" ]; then
    # Now Build the provider
    cd $home_dir/cmd/processProviders
    go build
    ./processProviders -config $home_dir/config.toml
    cd $home_dir/terraform_providers
    go mod tidy -go=1.16 && go mod tidy -go=1.17
    go build
else
    echo "No XML xpath file found. Try renaming the xpath file to include 'xpath' in the name."
fi
