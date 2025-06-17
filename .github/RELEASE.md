## Releasing

Releases are made on a reasonably regular basis by the maintainers, using the goreleaser tool. The following notes are only relevant to maintainers.

Release process:

1. Update [`version/VERSION`](https://github.com/opentofu/tofu-ls/blob/main/version/VERSION) and set it to the intended version to be released;
2. Create a Pull Request against main with the changes from this branch.
3. After the branch is merged, on your computer, make sure you have checked out the correct branch:
   * `main` for `alpha` and `beta` releases;
   * `vX.Y` for any other releases (assuming you are releasing version `X.Y.Z`)
2. Make sure the branch is up-to-date by running `git pull`;
3. Create the correct tag: `git tag -m "X.Y.Z" vX.Y.Z` (assuming you are releasing version `X.Y.Z`)
   * If you have a GPG key, consider adding the `-s` option to create a GPG-signed tag
4. Push the tag: `git push origin vX.Y.Z`;
5. This will trigger the `release.yml` workflow to create the release;
