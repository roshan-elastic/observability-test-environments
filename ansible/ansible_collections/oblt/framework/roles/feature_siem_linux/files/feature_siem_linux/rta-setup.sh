#!/usr/bin/env bash

# Set default values
username="elastic"
repo="detection-rules"
tag_prefix="integration-v"
version="8.12.5"
repo_location="$HOME/detection-rules"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    -username) username="$2"; shift 2;;
    -repo) repo="$2"; shift 2;;
    -tag_prefix) tag_prefix="$2"; shift 2;;
    -version) version="$2"; shift 2;;
    -repo_location) repo_location="$2"; shift 2;;
    *) echo "Unknown option: $1"; exit 1;;
  esac
done

# Output the params
echo "username: $username"
echo "repo: $repo"
echo "tag_prefix: $tag_prefix"
echo "version: $version"
echo "repo_location: $repo_location"

# Get the list of tags in the repository
# tags=$(curl -s "https://api.github.com/repos/${username}/${repo}/tags" | jq -c '.[] | select(.name | startswith("'"${tag_prefix}"'"))')

# Loop through each tag and print its name and archives
# echo "${tags}" | while read -r tag; do
#     name=$(echo "${tag}" | jq -r '.name')
#     tarball_url=$(echo "${tag}" | jq -r '.tarball_url')
#     zipball_url=$(echo "${tag}" | jq -r '.zipball_url')
#     echo "Tag name: ${name}"
#     echo "Tarball URL: ${tarball_url}"
#     echo "Zipball URL: ${zipball_url}"
#     echo
# done

tag="${tag_prefix}${version}"

# make sure the target directory for the repo exists, and if not, create it
if [ ! -d "$repo_location" ]; then
    mkdir -p "$repo_location"
    echo "Directory '$repo_location' created."
fi

# Download the tarball
curl -L "https://github.com/${username}/${repo}/archive/${tag}.tar.gz" -o "${tag}.tar.gz"

# Extract the tarball
tar -xzf "${tag}.tar.gz" -C "$repo_location" --strip-components 1

# Clean up - remove the downloaded tarball
rm "${tag}.tar.gz"
