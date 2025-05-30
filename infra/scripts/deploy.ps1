#!/usr/bin/env pwsh

# Check if required commands are available
$commands = @("azd", "az", "func")

foreach ($cmd in $commands) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Error "Error: $cmd command is not available, check pre-requisites in README.md"
        exit 1
    }
}

# Fetch Azure environment variables using azd env get-values
$output = azd env get-values

# Parse the output to get the endpoint values
foreach ($line in $output) {
    if ($line -match "SERVICE_API_NAME"){
        $AZURE_FUNCTION_NAME = ($line -split "=")[1] -replace '"',''
    }
}

# Publish the function app after building for Windows
Set-Location ./app
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -o main .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
func azure functionapp publish $AZURE_FUNCTION_NAME 

Write-Host "Deployment completed."
go build -o main .
Set-Location ../