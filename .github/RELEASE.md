## Releasing

Releases are made on a reasonably regular basis by the maintainers, using the goreleaser tool. The following notes are only relevant to maintainers.

Release process:

1. Update [`version/VERSION`](https://github.com/opentofu/tofu-ls/blob/main/version/VERSION) and set it to the intended version to be released
2. Wait for [`build` workflow](https://github.com/opentofu/tofu-ls/actions/workflows/build.yml) workflow to finish
3. Run the [Release workflow](https://github.com/opentofu/tofu-ls/actions/workflows/release.yml) with the appropriate version (matching the one in `version/VERSION`) & SHA (long one).
4. Wait for action release to finish.
