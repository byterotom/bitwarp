#!/bin/bash

set -e

GO_VERSION="1.24.2"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"

# install curl if it's missing
if ! command -v curl >/dev/null; then
    echo "[*] curl not found. Installing..."
    sudo apt update && sudo apt install -y curl
fi

# download and install Go
echo "[*] Downloading Go ${GO_VERSION}..."
curl -LO "$GO_URL"

echo "[*] Installing Go..."
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "$GO_TAR"
rm "$GO_TAR"

# setting PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# source .bashrc to apply changes immediately
source ~/.bashrc

# version check
go version