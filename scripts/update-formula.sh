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

# Create a temporary file for the updated formula
TEMP_FILE=$(mktemp)

# Process each platform and update the formula
awk -v version="$VERSION" '
BEGIN {
    # Define platform patterns
    platforms["darwin_arm64"] = "darwin_arm64"
    platforms["darwin_amd64"] = "darwin_amd64" 
    platforms["linux_arm64"] = "linux_arm64"
    platforms["linux_amd64"] = "linux_amd64"
}

{
    line = $0
    
    # Update version
    if (line ~ /version "/) {
        gsub(/version "[0-9.]*"/, "version \"" substr(version, 2) "\"")
    }
    
    # Check if this line contains a platform URL
    for (platform in platforms) {
        if (line ~ platform && line ~ /url/) {
            # Download and calculate SHA256 for this platform
            url = "https://github.com/bhanurp/jfrm/releases/download/" version "/jfrm_" platform ".tar.gz"
            cmd = "curl -sL \"" url "\" | shasum -a 256 | cut -d\" \" -f1"
            cmd | getline sha256
            close(cmd)
            
            print "  Processing " platform "..."
            print "  URL: " url
            print "  SHA256: " sha256 > "/dev/stderr"
            
            # Store the SHA256 for this platform
            sha256s[platform] = sha256
        }
    }
    
    # Update SHA256 lines for each platform
    for (platform in platforms) {
        if (line ~ platform && line ~ /sha256/) {
            if (sha256s[platform] != "") {
                gsub(/sha256 "[a-f0-9]{64}"/, "sha256 \"" sha256s[platform] "\"")
            }
        }
    }
    
    print line
}
' Formula/jfrm.rb > "$TEMP_FILE"

# Replace the original file
mv "$TEMP_FILE" Formula/jfrm.rb

echo ""
echo "Formula updated successfully!"
echo "Don't forget to commit and push the changes to your tap repository." 