# Contributing to OpenTofu Language Server

Welcome, and thank you for wanting to contribute!

## Get started

- Have a question? Post it in [GitHub Discussions ➡️](https://github.com/orgs/opentofu/discussions) or on the [OpenTofu Slack ➡️](https://opentofu.org/slack/)!
- Found a bug? [Report it here ➡️](https://github.com/opentofu/opentofu-ls/issues/new?assignees=&labels=bug%2Cpending-decision&projects=&template=bug_report.yml)
- Have a feature idea? [Submit it here ➡️](https://github.com/opentofu/opentofu-ls/issues/new?assignees=&labels=enhancement%2Cpending-decision&projects=&template=feature_request.yml)
- Want to provide a proof-of-concept for an issue? Please [submit a draft PR here ➡️](https://github.com/opentofu/opentofu-ls/compare)
- Want to add a feature, fix a linter error, refactor something, or add CI tooling?
  1. Check if there is an [open issues with the `accepted` label](https://github.com/opentofu/opentofu-ls/issues?q=is%3Aopen+is%3Aissue+label%3Aaccepted),
  2. Comment on the issue that you want to work on it,
  3. Wait for a maintainer to assign it to you,
  4. Then [submit your code here ➡️](https://github.com/opentofu/opentofu-ls/compare)
- Want to fix a bug? [Submit a PR here ➡️](https://github.com/opentofu/opentofu-ls/compare)
- Want to know what's going on? Read the [weekly updates ➡️](https://github.com/opentofu/opentofu/blob/main/WEEKLY_UPDATES.md), the [TSC summary ➡️](https://github.com/opentofu/opentofu/blob/main/TSC_SUMMARY.md) or join the [community meetings ➡️](https://meet.google.com/xfm-cgms-has) on Wednesdays at 14:30 CET / 8:30 AM Eastern / 5:30 AM Western / 19:00 India time on this link: https://meet.google.com/xfm-cgms-has ([📅 calendar link](https://calendar.google.com/calendar/event?eid=NDg0aWl2Y3U1aHFva3N0bGhyMHBhNzdpZmsgY18zZjJkZDNjMWZlMGVmNGU5M2VmM2ZjNDU2Y2EyZGQyMTlhMmU4ZmQ4NWY2YjQwNzUwYWYxNmMzZGYzNzBiZjkzQGc))

> [!TIP]
> For more OpenTofu events, subscribe to the [OpenTofu Events Calendar](https://calendar.google.com/calendar/embed?src=c_3f2dd3c1fe0ef4e93ef3fc456ca2dd219a2e8fd85f6b40750af16c3df370bf93%40group.calendar.google.com)!

**⚠️ Important:** Please avoid working on features or refactoring without [an
`accepted` issue](https://github.com/opentofu/opentofu-ls/issues?q=is%3Aopen+is%3Aissue+label%3Aaccepted). Every change
needs careful consideration. We cannot merge non-bug pull requests without discussing them, no matter how trivial the
issue may seem.

We specifically do not merge PRs **without prior issues** that:

- Reformat code
- Rename things
- Move code around
- Fix linter warnings for tools not currently in the CI pipeline
- Add new CI tooling

## Scope

This repository contains the source code only for the OpenTofu language server,
which in turn relies on other projects that have their own repositories.

[OpenTofu CLI/core has its own repository.](https://github.com/opentofu/opentofu)

OpenTofu providers are **not** maintained in this repository; you can find relevant
repository and relevant issue tracker for each provider within the
Registry index.

This repository also does **not** include the source code for some other parts of
OpenTofu. If you have feedback about them (including bug reports) please do feel free to
[open a GitHub issue in the core repository](https://github.com/opentofu/opentofu/issues/new/choose).

## Writing code for OpenTofu-ls

Eager to get started on coding? Here's the short version:

1. Set up a Go development environment with git.
2. Pay attention to copyright: [please read the DCO](https://developercertificate.org/), write the code yourself, avoid copy/paste. **Disable your AI coding assistant.**
3. Run the tests with `go test` in the package you are working on.
4. Build the Language Server by running either:
  - `go install` in the root of the repository. This will create a `opentofu-ls` executable in your `$GOBIN` (or
    `$GOPATH/bin`) directory.
  - `go build` in the root of the repository. This will create a `opentofu-ls` executable in the current directory.
5. Update [the changelog](CHANGELOG.md).
6. When you commit, use `git commit -s` to sign off your commits.
7. Complete the checklist below before you submit your PR (or submit a draft PR).
8. Your PR will be reviewed by the core team once it is marked as ready to review.

## PR checklist

<!-- Make sure to keep this in sync with the PR template. -->

Please make sure you complete the following checklist before you mark your PR ready for review. If you cannot complete the checklist but want to submit a PR, please submit it as a draft PR. Please note, the core team will only review your PR if you have completed the checklist and marked your PR as ready to review.

- [ ] I have read the contribution guidelines.
- [ ] I have not used an AI coding assistant to create this PR.
- [ ] I have written all code in this PR myself OR I have marked all code I have not written myself (including modified code, e.g. copied from other places and then modified) with a comment indicating where it came from.
- [ ] I (and other contributors to this PR) have not looked at the Terraform-ls source code while implementing this PR.

### Go checklist

If your PR contains a Go code, please make sure you check off all items on this list:

- [ ] I have run golangci-lint on my change and received no errors relevant to my code.
- [ ] I have run existing tests to ensure my code doesn't break anything.
- [ ] I have added tests for all relevant use cases of my code, which are passing.
- [ ] I have only exported functions, variables and structs that should be used from the other packages.
- [ ] I have added meaningful comments to all exported functions, variables, and structs.

## Development

If you wish to work on the source code, you'll first need to install
the [Go](https://golang.org/) compiler and the version control system
[Git](https://git-scm.com/).

Refer to the file [`.go-version`](.go-version) to see which version of Go
the language server is currently built with. Other versions will often work,
but if you run into any build or testing problems, please try the specific
Go version indicated. You can optionally simplify the installation of multiple
specific versions of Go on your system by installing
[`goenv`](https://github.com/syndbg/goenv), which reads `.go-version` and
automatically selects the correct Go version.

Use git to clone this repository into a location of your choice. Dependencies
are tracked via [Go Modules](https://blog.golang.org/using-go-modules),
and so you should _not_ clone it inside your `GOPATH`.

Switch into the root directory of the cloned repository and build
the language server

```
cd opentofu-ls
go install
```

Once the compilation succeeds, you can find a `opentofu-ls` executable in
the Go executable directory. If you haven't overridden it with the `GOBIN`
environment variable, the executable directory is the `bin` directory inside
the directory returned by the following command:

```
go env GOPATH
```

If you are planning to make changes to the source code, you should run the
unit test suite before you start to make sure everything is initially passing:

```
go test ./...
```

As you make your changes, you can re-run the above command to ensure that the
tests are _still_ passing. If you are working only on a specific Go package,
you can speed up your testing cycle by testing only that single package or
packages under a particular package prefix:

```
go test ./internal/terraform/exec/...
go test ./langserver
```

### Bundled schemas

The language server ships with bundled provider schemas of all official and partner providers to allow completions for common providers without the need to initialize an OpenTofu project first.
These schemas are not checked into version control and are automatically fetched and bundled when the language server is built in CI.

To generate these locally, run
```bash
go generate ./internal/schemas
```

The generated schemas can then be found in `internal/schemas/data`.

## External Dependencies

OpenTofu uses [Go Modules]((https://blog.golang.org/using-go-modules))
for dependency management.

If you need to add a new dependency to OpenTofu or update the selected version
for an existing one, use `go get` from the root of the OpenTofu repository
as follows:

```
go get github.com/hashicorp/hcl/v2@2.0.0
```

This command will download the requested version (2.0.0 in the above example)
and record that version selection in the `go.mod` file. It will also record
checksums for the module in the `go.sum`.

To complete the dependency change, clean up any redundancy in the module
metadata files by running the following command:

```
go mod tidy
```

Because dependency changes affect a shared, top-level file, they are more likely
than some other change types to become conflicted with other proposed changes
during the code review process. For that reason, and to make dependency changes
more visible in the change history, we prefer to record dependency changes as
separate commits that include only the results of the above commands and the
minimal set of changes to the language server's code for compatibility
with the new version:

```
git add go.mod go.sum
git commit -m "deps: go get github.com/hashicorp/hcl/v2@2.0.0"
```

You can then make use of the new or updated dependency in the new code added in
subsequent commits.

### Updating the changelog

We are keeping track of the changes to opentofu-ls in the [CHANGELOG.md](CHANGELOG.md) file. Please update it when you add features or fix bugs.

---

### Signing off your commits

When you contribute code to OpenTofu, we require you to add
a [Developer Certificate of Origin](https://developercertificate.org/) sign-off. Please read the DCO carefully before
you proceed and only contribute the code you have written yourself. Please do not add code that you have not written (
from scratch) yourself without first discussing it in the related issue.

The simplest way to add a sign-off is to use the `-s` command when you commit:

```
git commit -s -m "My commit message"
```

> [!IMPORTANT]
> Make sure your `user.name` and `user.email` settings in git match your GitHub settings. This will allow the automated
> DCO check to pass and avoid delays when merging your PR.

> [!TIP]
> Have you forgotten your sign-off? Click the "details" button on the failing DCO check for a guide on how to fix it!

---

### A note on copyright (WIP)

We take copyright and intellectual property very seriously. A few quick rules should help you:

1. When you submit a PR, you are responsible for the code in that pull request. You signal your acceptance of the [DCO](https://developercertificate.org/) with your sign-off.
2. If you include code in your PR that you didn't write yourself, make sure you have permission from the author. If you have permission, always add the `Co-authored-by` sign-off to your commits to indicate the author of the code you are adding.
3. Be careful about AI coding assistants! Coding assistants based on large language models (LLMs), such as ChatGPT or GitHub Copilot, are awesome tools to help. However, in the specific case of OpenTofu the training data may include the BSL-licensed Terraform. Since the OpenTofu/Terraform codebase is very specific and LLMs don't have any other training sources, they may emit copyrighted code. Please avoid using LLM-based coding assistants.
4. When you copy/paste code from within the OpenTofu code, always make it explicit where you copied from. This helps us resolve issues later on.
5. Before you copy code from external sources, make sure that the license allows this. Also make sure that any licensing requirements, such as attribution, are met. When in doubt, ask first!
6. Specifically, do not copy from the Terraform repository, or any PRs others have filed against that repository. This code is licensed under the BSL, a license which is not compatible with OpenTofu. (You may submit the same PR to both Terraform and OpenTofu as long as you are the author of both.)

> [!WARNING]
> To protect the OpenTofu project from legal issues violating these rules will immediately disqualify your PR from being merged and you from working on that area of the OpenTofu code base in the future. Repeat violations may get you barred from contributing to OpenTofu.

---

## Debugging

When launched as such, [PacketSender](https://packetsender.com) enables you to open a TCP socket with a server.
Approximate steps of debugging follow.

 - Install PacketSender (e.g. on MacOS via `brew cask install packet-sender`)
 - Launch LS in TCP mode: `opentofu-ls serve -port=8080`
 - Send any requests via PacketSender
   - Set `Address` to `127.0.0.1`
   - Set `Port` to `8080`
   - Tick `Persistent TCP`
   - Hit the `Send` button (which opens the TCP connection)
   - Paste or type the request in LSP format (see below) & hit `Send`

Examples of formatted requests follow.

```
Content-Length: 164\n\n{"jsonrpc":"2.0","params":{"textDocument":{"uri":"file:///var/path/to/file/main.tf"},"position":{"line":1,"character":0}},"method":"textDocument/completion","id":2}
```
```
Content-Length: 72\n\n{"jsonrpc":"2.0","params":{"id":2},"method":"$/cancelRequest","id":null}
```
```
Content-Length: 47\n\n{"jsonrpc":"2.0","method":"shutdown","id":null}
```

Keep in mind that each TCP session receives an isolated context,
so you cannot cancel requests you didn't start yourself
