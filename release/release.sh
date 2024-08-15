#!/usr/bin/env bash

set -eu -o pipefail
cd "$(dirname "$0")"

SCRATCH_DIR=$(mktemp -d)
echo "Using scratch directory: $SCRATCH_DIR"
mkdir -p "$SCRATCH_DIR"
trap 'rm -rf "$SCRATCH_DIR"' EXIT

echo "Running build script..."
./build.sh "$SCRATCH_DIR"

echo "Generating manifest.json..."
TEMP_MANIFEST=$(mktemp)
ls -1 "$SCRATCH_DIR" | jq -R -s 'split("\n") | map(select(length > 0)) | {files: .}' > "$TEMP_MANIFEST"
mv "$TEMP_MANIFEST" "$SCRATCH_DIR/manifest.json"

for file in "$SCRATCH_DIR"/*; do
    if [[ -f "$file" ]]; then
        ./upload.sh "$file"
    fi
done

echo "Release process completed successfully!"
