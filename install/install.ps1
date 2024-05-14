# PowerShell

$VERSION="0.3.1"

# Define the URL for the Windows binary
$WINDOWS_BINARY_URL="https://github.com/tome-gg/librarian/releases/download/$VERSION/tome-win.exe"

# Check if gh CLI tool is installed
if (!(Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Output "gh CLI tool is not installed. Installing now..."

    # Download and install the GitHub CLI
    Invoke-WebRequest -Uri "https://github.com/cli/cli/releases/download/v2.0.0/gh_2.0.0_windows_amd64.msi" -OutFile gh.msi
    Start-Process -Wait -FilePath msiexec -ArgumentList /i, (Resolve-Path .\gh.msi)
    Remove-Item .\gh.msi
}

# Download the binary
Invoke-WebRequest -Uri $WINDOWS_BINARY_URL -OutFile tome.exe

# Move the binary to a directory in the PATH
# Here we use C:\Windows\System32, which is generally in the PATH, but you should 
# adjust this to your preferred location.
Move-Item -Path .\tome.exe -Destination C:\Windows\System32\tome.exe
