<!--

** Thank you for your contribution! Please read this carefully! **

Please make sure you go through the checklist below. If your PR does not meet all requirements, please file it
as a draft PR. Core team members will only review your PR once it meets all the requirements below (unless your
change is something as trivial as a typo fix).

-->

<!-- If your PR resolves an issue, please add it here. -->
Resolves #

## Checklist

<!-- Please check of ALL items in this list for all PRs: -->

- [ ] I have read the [contribution guide](https://github.com/opentofu/tofu-ls/blob/main/.github/CONTRIBUTING.md).
- [ ] I have not used an AI coding assistant to create this PR.
- [ ] I have written all code in this PR myself OR I have marked all code I have not written myself (including modified code, e.g. copied from other places and then modified) with a comment indicating where it came from.
- [ ] I (and other contributors to this PR) have not looked at the Terraform source code while implementing this PR.
- [ ] If I'm releasing, I have read the [releasing guide](https://github.com/opentofu/tofu-ls/blob/main/.github/RELEASE.md).

### Go checklist

<!-- If your PR contains Go code, please make sure you check off all items on this list: -->

- [ ] I have run golangci-lint on my change and receive no errors relevant to my code.
- [ ] I have run existing tests to ensure my code doesn't break anything.
- [ ] I have added tests for all relevant use cases of my code, and those tests are passing.
- [ ] I have only exported functions, variables and structs that should be used from other packages.
- [ ] I have added meaningful comments to all exported functions, variables, and structs.
