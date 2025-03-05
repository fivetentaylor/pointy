#!/bin/bash

# Run both commands in the background
go run -mod=mod github.com/a-h/templ/cmd/templ generate -path pkg/ &
PID1=$!

GOOS=js GOARCH=wasm go build -ldflags "-X main.ImageTag=$IMAGE_TAG" -o pkg/assets/static/rogueV3.wasm rogue/v3/wasm/main.go &
PID2=$!

(cd pkg/admin/src && npx tailwindcss -i admin.css -o ../static/admin.css --minify) &
PID3=$!

# Wait for allcommands to finish
wait $PID1
EXIT_STATUS1=$?

wait $PID2
EXIT_STATUS2=$?

wait $PID3
EXIT_STATUS3=$?

# Initialize an error flag
ERROR_OCCURRED=0

# Check if the first command failed
if [ $EXIT_STATUS1 -ne 0 ]; then
    echo "Error: The 'templ generate' command failed."
    ERROR_OCCURRED=1
fi

# Check if the second command failed
if [ $EXIT_STATUS2 -ne 0 ]; then
    echo "Error: The 'go build for wasm' command failed."
    ERROR_OCCURRED=1
fi

# Check if the third command failed
if [ $EXIT_STATUS3 -ne 0 ]; then
    echo "Error: The 'tailwind' command failed."
    ERROR_OCCURRED=1
fi

# Exit if any command failed
if [ $ERROR_OCCURRED -ne 0 ]; then
    exit 1
fi

# Proceed with the rest of the script if both commands succeeded
cp pkg/assets/static/rogueV3.wasm pkg/admin/static/rogueV3.wasm
go build -v -o ./tmp/main ./cmd/reviso/main.go
