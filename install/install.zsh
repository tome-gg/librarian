#!/bin/zsh

$VERSION="0.3.1"

# Define the URLs for the binaries
MAC_BINARY_URL="https://github.com/tome-gg/librarian/releases/download/$VERSION/tome-darwin-arm-osx-m1"
LINUX_BINARY_URL="https://github.com/tome-gg/librarian/releases/download/$VERSION/tome-linux-amd64"

# Check if gh CLI tool is installed
if ! command -v gh &> /dev/null; then
    echo "gh CLI tool is not installed. Installing now..."

    # Detect the operating system
    OS="$(uname)"

    # Install gh CLI tool
    if [ "$OS" == "Darwin" ]; then
        # macOS
        brew install gh
    elif [ "$OS" == "Linux" ]; then
        # Linux
        sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-key C99B11DEB97541F0
        sudo apt-add-repository https://cli.github.com/packages
        sudo apt update
        sudo apt install gh
    else
        echo "Unsupported operating system: $OS"
        exit 1
    fi
fi

# Detect the operating system
OS="$(uname)"

# Download and configure the appropriate binary
if [ "$OS" == "Darwin" ]; then
    # macOS
    curl -L $MAC_BINARY_URL -o tome
elif [ "$OS" == "Linux" ]; then
    # Linux
    curl -L $LINUX_BINARY_URL -o tome
else
    echo "Unsupported operating system: $OS"
    exit 1
fi

# Make the binary executable
chmod +x tome

# Move the binary to a directory in the PATH
sudo mv tome /usr/local/bin/
