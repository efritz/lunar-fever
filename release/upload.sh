#!/usr/bin/env bash

set -eu -o pipefail

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <file_to_upload>"
    exit 1
fi

if [ ! -f "$1" ]; then
    echo "Error: File '$1' not found."
    exit 1
fi

if [ -z "$AUTH_PASSWORD" ]; then
    echo "Error: AUTH_PASSWORD environment variable is not set."
    exit 1
fi

response=$(curl -s -w "%{http_code}" -o /dev/null -X POST \
    -H "Content-Type: multipart/form-data" \
    -u ":$AUTH_PASSWORD" \
    -F "file=@$1" \
    "https://lunarfever.com/upload")

if [ "$response" -eq 200 ]; then
    echo "File uploaded successfully."
else
    echo "Error occurred while uploading the file. HTTP status code: $response"
    exit 1
fi
