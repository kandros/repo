#!/bin/bash

# Exit on error
set -e

# Get the latest tag, default to v0.0.0 if no tags exist
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
MAJOR=$(echo $LATEST_TAG | cut -d. -f1 | tr -d 'v')
MINOR=$(echo $LATEST_TAG | cut -d. -f2)
PATCH=$(echo $LATEST_TAG | cut -d. -f3)

case $1 in
    major)
        NEW_VER="v$((MAJOR+1)).0.0"
        ;;
    minor)
        NEW_VER="v${MAJOR}.$((MINOR+1)).0"
        ;;
    patch)
        NEW_VER="v${MAJOR}.${MINOR}.$((PATCH+1))"
        ;;
    *)
        echo "Usage: $0 {major|minor|patch}"
        exit 1
        ;;
esac

echo "Current version: $LATEST_TAG"
echo "Creating new $1 version: $NEW_VER"
git tag -a "$NEW_VER" -m "Release $NEW_VER"