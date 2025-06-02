#!/bin/bash
set -e 

# Function to check if a command is available
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

retry_script() {
  while true; do
    ./generateFiles.sh
    
    # Check the exit status of the script
    if [ $? -eq 0 ]; then
        # Script ran successfully, exit the loop
        break
    else
        # Script failed, display an error message
        echo "Script failed. Retrying in 5 seconds..."
        sleep 5
    fi
  done
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

# Prompt the user for their choice
echo "Do you want to:"
echo "1. Build a provider from scratch"
echo "2. Provide a configuration"
read -p "Enter your choice (1/2): " choice

# Check the user's choice
if [ "$choice" == "1" ]; then
    echo "You chose to build a provider from scratch."
    # Create config.toml file
    if [ ! -f "config.toml" ]; then
      echo "Creating config.toml file..."
      cat << EOF > config.toml
    yangDir = "$(pwd)/yang_files"
    providerDir = "$(pwd)/terraform_providers"
    xpathPath = "$(pwd)/xpath_inputs.xml"
    providerName = "vsrx"
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
  
elif [ "$choice" == "2" ]; then
    echo "You chose to provide a configuration."

    script_directory="$home_dir/user_config_files"

    # Create the directory if it doesn't exist
    mkdir -p "$script_directory"

    # Initialize the config_file variable
    config_file=""

    echo ""
    echo ""

    # Prompt the user for the configuration file name and validate it
    while [ ! -f "$config_file" ]; do
      read -p "Enter the configuration file name: " config_file
      full_path="$script_directory/$config_file"
      
      if [ -f "$full_path" ]; then
        config_file="$full_path"
      else
        echo "The provided file does not exist in your script's directory."
        config_file=""
      fi
    done

    # Ask user what device they are working on
    valid_options=("vsrx" "vmx" "vqfx" "vptx")
    user_input=""

    while [[ ! " ${valid_options[*]} " =~ " $user_input " ]]; do
        read -p "Enter a valid device option (vsrx, vmx, vqfx, vptx): " user_input
        
        if [[ ! " ${valid_options[*]} " =~ " $user_input " ]]; then
            echo "Invalid input. Please enter one of the following options: vsrx, vmx, vqfx, vptx."
        fi
    done

    # Prompt the user for the Junos version
    read -p "Enter the Junos version: " junos_version

    echo ""
    echo ""
    # Display the user's selections
    echo "Configuration file path: $config_file"
    echo "Device name: $user_input"
    echo "Junos version: $junos_version"


    # Ask for confirmation
    while true; do
      read -p "Are these selections correct? (yes/no): " confirmation
      case "$confirmation" in
        [Yy]* ) break;;
        [Nn]* ) 
          # If the user says no, allow them to start over
          exec "$0";;
        * ) echo "Please answer 'yes' or 'no'.";;
      esac
    done

    echo ""
    echo ""
    # Create or overwrite config.toml file
    echo "Creating or overwriting config.toml file..."
    cat > config.toml << EOF
    yangDir = "$(pwd)/yang_files"
    providerDir = "$(pwd)/terraform_providers"
    xpathPath = "$(pwd)/xpath_inputs.xml"
    providerName = "$user_input"
    fileType = "both"
EOF

    echo ""
    echo ""

    junos_version_combined="${junos_version}R1"

    device_names=("vmx" "vptx" "vsrx" "vqfx")
    device_mappings=("junos" "junos" "junos-es" "junos-qfx")

    # Find the index of the input device name in the device_names array
    index=-1
    for ((i=0; i<${#device_names[@]}; i++)); do
        if [ "${device_names[$i]}" = "$user_input" ]; then
            index=$i
            break
        fi
    done

    # Output the corresponding mapping if found
    if [ $index -ne -1 ]; then
        supported_devices="${device_mappings[$index]}"
    else
        echo "Device mapping not found."
    fi

    # Use the user's input to look up the supported devices
    supported_devices="${device_mappings[${user_input}]}"

    # Check if the device name exists in the mapping
    if [ -z "$supported_devices" ]; then
      echo "Device name not recognized."
      exit 1
    fi

    # # Output the supported devices
    # echo "Supported devices for $selected_device: $supported_devices"
    year_prefix=${junos_version%%.*}
    # common_path="yang/$junos_version/$junos_version_combined/common/junos-common-types@2023-01-01.yang"
    common_path="yang/$junos_version/$junos_version_combined/common/junos-common-types@20${year_prefix}-01-01.yang"
    path_to="yang/$junos_version/$junos_version_combined/$supported_devices/conf/"

    # Define the target directory in the home directory
    target_dir="$home_dir/yang_files"

    # Check if the target directory already exists, and create it if not
    if [ ! -d "$target_dir" ]; then
      echo "yang_files folder does not exist. Creating and adding necessary files"
      mkdir -p "$target_dir"
      target_dir="$home_dir/yang"

      # Check if the target directory already exists, and create it if not
      if [ ! -d "$target_dir" ]; then
        mkdir -p "$target_dir"
        # Git clone the Juniper YANG repository into the target directory
        echo "Cloning the Juniper YANG repository for Junos $junos_version into $target_dir..."
        git clone https://github.com/Juniper/yang.git "$target_dir"
      fi

      # Check if the git clone was successful
      if [ $? -eq 0 ]; then
        echo "Cloning successful. YANG files are in $target_dir."
      else
        echo "Folder already created or clone failed. If clone failed, please check your internet connection or repository URL."
      fi

      target_dir="$home_dir/yang_files"

      # Copy all files from the source directory to the target directory
      cp -r "$path_to"* "$target_dir"

      # Copy the common file to the target directory
      cp "$common_path" "$target_dir"

      # Check if the copy operation was successful
      if [ $? -eq 0 ]; then
        echo "Files copied successfully to $target_dir."
      else
        echo "Copy operation failed."
      fi

      target_dir="$home_dir/yang"

      rm -rf $target_dir
    fi

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
      # go run $home_dir/Internal/processYang/createXpathInputs.go
    else
      echo "No .yang files found in yang_files folder. Add files and re-run script"
    fi

    cd "$home_dir"

    go run createXpathInputs.go

    echo ""
    echo ""

    # Search for XML files containing "xpath" in their filenames in the home directory
    xml_files=$(find "$PWD" -type f -name '*xpath*.xml')

    # Define the folder name
    folderName="TFtemplates"

    # Check if the folder already exists
    if [ ! -d "$folderName" ]; then
        # Create the folder if it doesn't exist
        mkdir "$folderName"
        echo "Folder '$folderName' created successfully. Folder stores tf templates for testing"
    else
        echo "Folder '$folderName' already exists to store .tf templates"
    fi

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

    # Define the target directory in the home directory
    target_dir="$home_dir/testbed"

    # Check if the target directory already exists
    if [ ! -d "$target_dir" ]; then
        # If it doesn't exist, create the directory
        mkdir -p "$target_dir"
        echo "Created target directory: $target_dir"
    fi

    # Check if main.tf does not exist in the testbed directory
    if [ ! -f "$target_dir/main.tf" ]; then
        # Create main.tf
        touch "$target_dir/main.tf"
        echo "Created main.tf in the testbed directory"
    fi

    # Create a folder named after device_name appended with "_1"
    device_folder="$target_dir/${user_input}_1"
    if [ ! -d "$device_folder" ]; then
        # If it doesn't exist, create the folder
        mkdir -p "$device_folder"
        echo "Created folder: $device_folder"
    fi

    # Check if main.tf does not exist in the device folder
    if [ ! -f "$device_folder/main.tf" ]; then
        # Create main.tf in the device folder
        touch "$device_folder/main.tf"
        echo "Created main.tf in the device folder: $device_folder"
    fi

else
    echo "Invalid choice. Please enter 1 or 2."
fi

