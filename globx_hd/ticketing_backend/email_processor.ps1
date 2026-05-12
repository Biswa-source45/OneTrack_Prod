# PowerShell script to run the email processor

# Navigate to the project directory
$projectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectDir

# Ensure the environment variables are loaded
if (Test-Path .env.email) {
    # Read and set environment variables from .env.email file
    Get-Content .env.email | ForEach-Object {
        if ($_ -match "^\s*([^#].*?)=(.*)$") {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            # Remove quotes if present
            $value = $value -replace '^["'']|["'']$'
            [Environment]::SetEnvironmentVariable($name, $value)
        }
    }
}

# Build the email processor if needed
go build -o bin/email_processor.exe ./cmd/email_processor

# Run the email processor
./bin/email_processor.exe
