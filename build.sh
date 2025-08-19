#!/bin/bash
cd web

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    yarn install
fi

# Compute checksum of web source files
NEW_HASH=$(find src -type f -exec md5sum {} + | sort -k 2 | md5sum | awk '{print $1}')
HASH_FILE=".build_hash"

# Compare with last build
if [ -f "$HASH_FILE" ]; then
    OLD_HASH=$(cat "$HASH_FILE")
else
    OLD_HASH=""
fi

if [ "$NEW_HASH" != "$OLD_HASH" ]; then
    rm -rf dist/
    yarn run build
    if [ $? -eq 0 ]; then
        echo "$NEW_HASH" > "$HASH_FILE"
    else
        exit 1
    fi
fi

cd ..

# Build Go server
go build -o kodo ./main.go
if [ $? -ne 0 ]; then
    exit 1
fi
