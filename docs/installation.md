# Installation

## Automatic Installation

Some editors have built-in logic to install and update the language server automatically, so you typically shouldn't need to worry about installation or updating of the server in these editors, as long as you use the linked extension.
<!-- TODO: Update this link when we get a better display and itemName. See https://github.com/opentofu/vscode-opentofu/issues/30 -->
 - OpenTofu VS Code extension [stable](https://marketplace.visualstudio.com/items?itemName=opentofu.vscode-opentofu)
 - OpenTofu Zed extension [stable](https://zed.dev/extensions?query=OpenTofu)
<!-- We don't have a Sublime Text version yet [Sublime Text LSP-terraform](https://packagecontrol.io/packages/LSP-terraform) -->

## Manual Installation

You can install the language server manually using one of the many package managers available or download an archive from the release page. After installation, follow the [install instructions for your IDE](./USAGE.md)

### Homebrew (macOS / Linux)

You can install via [Homebrew](https://brew.sh)

```shell
brew install tofu-ls
```

This tap only contains stable releases (i.e. no pre-releases). -->

### All platforms

1. [Download for the latest version](https://github.com/opentofu/tofu-ls/releases)
  of the language server relevant for your operating system and architecture.
2. The language server is distributed as a single binary.
  Install it by unzipping it and moving it to a directory
  included in your system's `PATH`.
3. You can verify integrity by comparing the SHA256 checksums
  which are part of the release (called `tofu-ls_<VERSION>_SHA256SUMS`).
4. Check that you have installed the server correctly via `tofu-ls -v`.
  You should see the latest version printed to your terminal. -->
