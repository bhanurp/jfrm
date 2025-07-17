#!/bin/bash

# Script to create a new release
# Usage: ./scripts/create-release.sh <version>
# Example: ./scripts/create-release.sh v0.1.0

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.1.0"
    exit 1
fi

echo "Creating release: $VERSION"
echo ""

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Warning: You're not on the main branch. Current branch: $CURRENT_BRANCH"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    echo "Error: Tag $VERSION already exists!"
    exit 1
fi

# Create and push tag
echo "Creating tag: $VERSION"
git tag $VERSION

echo "Pushing tag to remote..."
git push origin $VERSION

echo ""
echo "✅ Release $VERSION created!"
echo "📋 GitHub Actions will automatically:"
echo "   - Build binaries for all platforms"
echo "   - Create a GitHub release"
echo "   - Update the Homebrew tap"
echo ""
echo "🔗 Check the progress at: https://github.com/bhanurp/jfrm/actions" 