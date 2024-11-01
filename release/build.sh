#!/usr/bin/env bash

set -eu -o pipefail
SCRATCH_DIR="${1:-.}"
cd "$(dirname "$0")/.."

# Add timestamp to filenames
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo 'Compiling for Apple Silicon...'
env GOOS='darwin' \
    GOARCH='arm64' \
    go build \
        -ldflags -linkmode=external \
        -o "$SCRATCH_DIR/lunar-fever-darwin-arm64_$TIMESTAMP" \
        ./cmd/lunar-fever
shasum -a 256 "$SCRATCH_DIR/lunar-fever-darwin-arm64_$TIMESTAMP" | cut -d' ' -f1 > "$SCRATCH_DIR/lunar-fever-darwin-arm64_$TIMESTAMP.sum256"
echo 'Done.'

# Note: requires `brew install mingw-w64`
#
# echo 'Compiling for Windows...'
# env GOOS='windows' \
#     GOARCH='amd64' \
#     CGO_ENABLED='1' \
#     CC='x86_64-w64-mingw32-gcc' \
#     go build \
#         -ldflags -linkmode=external \
#         -o "$SCRATCH_DIR/lunar-fever-windows-amd64_$TIMESTAMP.exe" \
#         ./cmd/lunar-fever
# shasum -a 256 "$SCRATCH_DIR/lunar-fever-windows-amd64_$TIMESTAMP.exe" | cut -d' ' -f1 > "$SCRATCH_DIR/lunar-fever-windows-amd64_$TIMESTAMP.exe.sum256"
# echo 'Done.'
