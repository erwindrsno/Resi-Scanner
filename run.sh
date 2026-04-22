#!/bin/bash
# 1. Name of your binary
APP_NAME="resi-scanner"

echo "Building..."

# 2. Run go build for the whole directory
# Use -o to specify the exact output name
go build -o $APP_NAME .

# 3. Check if build was successful ($? is the exit code)
if [ $? -eq 0 ]; then
  echo "Build successful. Starting app..."
  echo "---------------------------------"
  ./$APP_NAME
else
  echo "Build failed! Fix the errors above."
fi
