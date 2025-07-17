#!/bin/bash

# Script to update Homebrew formula with correct SHA256 hashes
# Usage: ./scripts/update-formula.sh <version>

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.1.0"
    exit 1
fi

echo "Updating Homebrew formula for version: $VERSION"
echo ""

# Download and calculate SHA256 for each platform
PLATFORMS=(
    "darwin_amd64"
    "darwin_arm64"
    "linux_amd64"
    "linux_arm64"
)

for platform in "${PLATFORMS[@]}"; do
    echo "Processing $platform..."
    
    # Download the release asset
    URL="https://github.com/bhanurp/jfrm/releases/download/$VERSION/jfrm_${platform}.tar.gz"
    
    # Calculate SHA256
    SHA256=$(curl -sL "$URL" | shasum -a 256 | cut -d' ' -f1)
    
    echo "  URL: $URL"
    echo "  SHA256: $SHA256"
    echo ""
    
    # Update the formula file
    sed -i.bak "s|url \"https://github.com/bhanurp/jfrm/releases/download/$VERSION/jfrm_${platform}.tar.gz\"|url \"$URL\"|g" Formula/jfrm.rb
    sed -i.bak "s|sha256 \"YOUR_SHA256_HERE\"|sha256 \"$SHA256\"|g" Formula/jfrm.rb
done

# Update version
sed -i.bak "s|version \"0.1.0\"|version \"${VERSION#v}\"|g" Formula/jfrm.rb

echo "Formula updated successfully!"
echo "Don't forget to commit and push the changes to your tap repository." 