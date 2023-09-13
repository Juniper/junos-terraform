#!/bin/bash

# Function to check if a command is available
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# caputres the user's home directory
home_dir="$PWD"

# Check if Python is installed
if ! command_exists python3; then
  echo "Python is not installed. Please Install 'Python' before running the script"
fi

# Check if Python is installed
if ! command_exists go; then
  echo "Go is not installed. Please Install 'Go' before running the script."
fi

# Ask user what device they are working on
valid_options=("vsrx" "vmx" "vqfx" "vptx")
user_input=""

while [[ ! " ${valid_options[*]} " =~ " $user_input " ]]; do
    read -p "Enter a valid option (vsrx, vmx, vqfx, vptx): " user_input
    
    if [[ ! " ${valid_options[*]} " =~ " $user_input " ]]; then
        echo "Invalid input. Please enter one of the following options: vsrx, vmx, vqfx, vptx."
    fi
done

# Create config.toml file
if [ ! -f "config.toml" ]; then
  echo "Creating config.toml file..."
  cat << EOF > config.toml
yangDir = "$(pwd)/yang_files"
providerDir = "$(pwd)/terraform_providers"
xpathPath = "$(pwd)/xpath_example.xml"
providerName = "$device_name"
fileType = "both"
EOF
fi

# Check if yang_files folder exists
if [ -d "yang_files" ]; then
  # Find at least one .yang file in yang_files folder
  yang_files=$(find yang_files -name "*.yang")
  # If YANG files are found, generate YIN and Xpath Files
  if [ -n "$yang_files" ]; then
    echo "Found .yang files in yang_files folder."
    # Change directory to /cmd/processYang
    cd $home_dir/cmd/processYang || exit 1
    # Activate venv
    python3 -m venv venv
    source venv/bin/activate
    # Check and install pyang if needed
    if ! command_exists pyang; then
      echo "pyang is not installed. Installing pyang..."
      pip install pyang
      pyang -v
    fi
    # Run go build command and generate YIN and Xpath Files
    go build
    ./processYang -config $home_dir/config.toml
    deactivate
  else
    echo "No .yang files found in yang_files folder. Add files and re-run script"
  fi
else
  echo "yang_files folder does not exist. Create and add necessary files"
fi