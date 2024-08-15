#!/usr/bin/env bash

set -eu -o pipefail
cd "$(dirname "$0")"

SCRATCH_DIR=$(mktemp -d)
echo "Using scratch directory: $SCRATCH_DIR"
mkdir -p "$SCRATCH_DIR"
trap 'rm -rf "$SCRATCH_DIR"' EXIT

echo "Running build script..."
./build.sh "$SCRATCH_DIR"

for file in "$SCRATCH_DIR"/*; do
    if [[ -f "$file" ]]; then
        ./upload.sh "$file"
    fi
done

echo "Release process completed successfully!"
