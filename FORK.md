# About this fork

Personal fork of [opentofu/tofu-ls](https://github.com/opentofu/tofu-ls). It adds features that upstream
cannot take because of its [AI usage policy](./AI-USAGE-POLICY.md). Current additions:

- `textDocument/rename` and `textDocument/prepareRename` for resource, data, variable, local and module names.

Everything else is upstream. Don't report issues from this fork's builds upstream.

## Versions

Releases track upstream's semver tags and add the commit as build metadata: `v<upstream version>+<short commit>`,
e.g. `v0.5.3+158d88e`. `tofu-ls version` prints the same string.

## Syncing and releasing

```bash
git fetch upstream --tags
git merge upstream/main          # or the upstream tag you want to follow
# upstream's version/VERSION must equal an upstream tag (e.g. 0.5.3 -> v0.5.3)
git push origin main
scripts/release.sh               # tags v$(cat version/VERSION)+<short hash> and pushes it
```

The tag triggers `.github/workflows/release.yml`. It fails unless the tag matches `version/VERSION`,
the tagged commit's short hash, and a descendant of the upstream tag. GoReleaser then builds
linux/windows/darwin archives and publishes the GitHub release.

Local build with the commit embedded:

```bash
go build -ldflags "-X main.commit=$(git rev-parse --short=7 HEAD)" -o tofu-ls .
```
