#!/usr/bin/env bash

# Set default values
repo_location="$HOME/detection-rules"
os=$(uname | tr '[:upper:]' '[:lower:]')

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    -repo_location) repo_location="$2"; shift 2;;
    -os) os="$2"; shift 2;;
    *) echo "Unknown option: $1"; exit 1;;
  esac
done

# Output the params
echo "repo_location: $repo_location"
echo "os: $os"

# Get the current directory
original_dir=$(pwd)

# Go to the repo directory
cd "$repo_location" || exit 1

# check which python is available
python_cmd=""
if command -v python3 &> /dev/null; then
  python_cmd="python3"
elif command -v python &> /dev/null; then
  python_cmd="python"
else
  echo "No Python found!"
  exit
fi

# Run the command and store its output
output=$($python_cmd -m rta -l -o "$os")

# Extract lines starting from the third line, skip rows that start with spaces, and select the first word
attacks=$(echo "$output" | sed -n '3,$p' | grep -v '^[[:space:]]' | awk '{print $1}')

# Select a random entry from the attacks
random_attack=$(echo "$attacks" | shuf -n 1)

echo "Running random attack: ${random_attack%???}"
$python_cmd -m rta -n "${random_attack%???}"

# Return to the original directory
cd "$original_dir" || exit 1
