#!/usr/bin/env bash

# Check if required commands are available
commands=("azd" "az" "func")

for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &> /dev/null; then
    echo "Error: $cmd command is not available, check pre-requisites in README.md" >&2
    exit 1
  fi
done

# Fetch Azure environment variables using azd env get-values
output=$(azd env get-values)

# Parse the output to get the endpoint values
while IFS= read -r line; do
  if [[ $line == SERVICE_API_NAME=* ]]; then
    AZURE_FUNCTION_NAME=$(echo "$line" | cut -d'=' -f2 | tr -d '"')
  fi
done <<< "$output"

# Publish the function app after building for Windows
cd ./app || exit
env GOOS=windows GOARCH=386 go build -o main .
func azure functionapp publish $AZURE_FUNCTION_NAME

echo "Deployment completed."
go build -o main .
cd ../ || exit