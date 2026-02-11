#!/usr/bin/env bash
set -euo pipefail

TAG=v0.1.17
BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ -n "$(git status --porcelain)" ]; then
  echo "Working tree is dirty. Please commit or stash changes first."
  git status --porcelain
  exit 1
fi

# Commit CHANGELOG if not already staged
if git diff --name-only --cached | grep -q "CHANGELOG.md"; then
  echo "CHANGELOG.md already staged for commit"
else
  git add CHANGELOG.md
  git commit -m "chore(release): ${TAG}"
fi

# Create annotated tag
if git rev-parse "${TAG}" >/dev/null 2>&1; then
  echo "Tag ${TAG} already exists. Aborting."
  exit 1
fi

git tag -a "${TAG}" -m "Release ${TAG}"

echo "Pushing branch ${BRANCH} and tag ${TAG} to origin"

git push origin "${BRANCH}"
git push origin "${TAG}"

# Optional: run goreleaser (requires configuration and secrets)
# goreleaser release --rm-dist

echo "Release ${TAG} created and pushed."

echo "Note: CI workflows or goreleaser will handle artifact creation and registry publication."
