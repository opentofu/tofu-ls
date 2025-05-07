## Releasing (WIP)

!!TODO This test is inherited from before the fork, this will be updated after we setup the final version of the release process.

Releases are made on a reasonably regular basis by the maintainers, using our internal tooling. The following notes are only relevant to maintainers.

Release process:

1. Update [`version/VERSION`](https://github.com/opentofu/tofu-ls/blob/main/version/VERSION) to remove `-dev` suffix and
   set it to the intended version to be released
2. Wait for [`build` workflow](https://github.com/opentofu/tofu-ls/actions/workflows/build.yml) and dependent `prepare`
   workflow to finish
3. Run the [Release workflow](https://github.com/opentofu/tofu-ls/actions/workflows/release.yml) with the appropriate
   version (matching the one in `version/VERSION`) & SHA (long one).
4. Wait for `staging` release to finish.
