#!/bin/bash

# captures the user's home directory
home_dir="$PWD"

# Search for XML files containing "xpath" in their filenames in the home directory
xml_files=$(find "$PWD" -type f -name '*xpath*.xml')

# Define the folder name
folderName="TFtemplates"

# Check if the folder already exists
if [ ! -d "$folderName" ]; then
    # Create the folder if it doesn't exist
    mkdir "$folderName"
    echo "Folder '$folderName' created successfully."
else
    echo "Folder '$folderName' already exists."
fi


# Check if any matching XML files were found
if [ -n "$xml_files" ]; then
    # Now Build the provider
    cd $home_dir/cmd/processProviders
    go build
    ./processProviders -config $home_dir/config.toml
    cd $home_dir/terraform_providers
    go mod tidy -go=1.22
    go build
else
    echo "No XML xpath file found. Try renaming the xpath file to include 'xpath' in the name."
fi
