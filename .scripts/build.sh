#!/bin/bash

cd "$(dirname "$0")/../pbf-bridge"

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/pbf-bridge-linux ./cmd/server

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/pbf-bridge-windows.exe ./cmd/server

echo "Complete"