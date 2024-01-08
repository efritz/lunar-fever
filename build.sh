#! /usr/bin/env bash

set -eux -o pipefail

echo 'Compiling for Apple Silicon...'
env GOOS='darwin' \
    GOARCH='arm64' \
    go build -o lunar-fever-macos-arm64 \
    ./cmd/lunar-fever
shasum -a 256 lunar-fever-macos-arm64 | cut -d' ' -f1 > lunar-fever-macos-arm64.sum256
echo 'Done.'

# Note: requires `brew install mingw-w64`
echo 'Compiling for Windows...'
env GOOS='windows' \
    GOARCH='amd64' \
    CGO_ENABLED='1' \
    CC='x86_64-w64-mingw32-gcc' \
    go build -o lunar-fever-windows-amd64.exe \
    ./cmd/lunar-fever
shasum -a 256 lunar-fever-windows-amd64.exe | cut -d' ' -f1 > lunar-fever-windows-amd64.exe.sum256
echo 'Done.'
