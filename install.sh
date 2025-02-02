#!/bin/bash

# Exit immediately if any command fails
set -e

BINARY_NAME="repo"

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert machine hardware names to match GoReleaser's conventions
# GoReleaser uses standardized names: amd64, 386, arm64, arm
case "${ARCH}" in
    x86_64)  ARCH="amd64" ;;
    i386|i686) ARCH="386" ;;
    aarch64) ARCH="arm64" ;;
    armv7l|armv6l) ARCH="arm" ;;
esac

# Fetch the latest release and extract the correct download URL using grep and sed
LATEST_VERSION=$(curl -s "https://api.github.com/repos/kandros/repo/releases/latest" | grep -o '"tag_name": "v[^"]*"' | sed 's/"tag_name": "v\(.*\)"/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not determine latest version"
    exit 1
fi

DOWNLOAD_URL="https://github.com/kandros/repo/releases/download/v${LATEST_VERSION}/repo_${LATEST_VERSION}_${OS}_${ARCH}.tar.gz"

# Verify the URL exists
if ! curl --output /dev/null --silent --head --fail "$DOWNLOAD_URL"; then
    echo "Error: Release not found for ${OS}_${ARCH}"
    exit 1
fi

TMP_DIR=$(mktemp -d)
TMP_FILE="${TMP_DIR}/${BINARY_NAME}.tar.gz"

curl -sL "$DOWNLOAD_URL" -o "$TMP_FILE"
tar xzf "$TMP_FILE" -C "$TMP_DIR"

# Install to /usr/local/bin which is in most users' PATH
# Using sudo since this directory typically requires root access
sudo mv "${TMP_DIR}/${BINARY_NAME}" "/usr/local/bin/${BINARY_NAME}"
sudo chmod +x "/usr/local/bin/${BINARY_NAME}"

rm -rf "$TMP_DIR"

echo "${BINARY_NAME} has been installed successfully!" 