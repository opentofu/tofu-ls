#!/usr/bin/env bash
# Tag HEAD as v<upstream version>+<short commit> and push it; CI (release.yml) builds the release.
set -euo pipefail

base=$(tr -d '[:space:]' < version/VERSION)
tag_base="v$base"

[ -z "$(git status --porcelain)" ] || { echo "working tree is dirty" >&2; exit 1; }

git fetch --quiet --no-tags upstream "refs/tags/$tag_base:refs/tags/$tag_base" 2>/dev/null || true
git rev-parse -q --verify "refs/tags/$tag_base" >/dev/null ||
  { echo "upstream has no tag $tag_base; merge an upstream release first" >&2; exit 1; }
git merge-base --is-ancestor "$tag_base" HEAD ||
  { echo "HEAD does not contain upstream $tag_base" >&2; exit 1; }

git fetch --quiet origin
git branch -r --contains HEAD | grep -q '^ *origin/' ||
  { echo "HEAD is not pushed to origin; push it first" >&2; exit 1; }

tag="$tag_base+$(git rev-parse --short=7 HEAD)"
git tag -a "$tag" -m "$tag"
git push origin "$tag"
echo "pushed $tag; watch the Release workflow"
