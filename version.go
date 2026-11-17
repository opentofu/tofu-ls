// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	_ "embed"
	"strings"

	goversion "github.com/hashicorp/go-version"
)

var (
	// The next version number that will be released. This will be updated after every release
	// Version must conform to the format expected by github.com/hashicorp/go-version
	// for tests to work.
	// A pre-release marker for the version can also be specified (e.g -dev). If this is omitted
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	//go:embed version/VERSION
	rawVersion string

	// Short commit hash, set by goreleaser via -ldflags (see FORK.md); empty in dev builds.
	commit string

	version = withCommit(goversion.Must(goversion.NewVersion(strings.TrimSpace(rawVersion))).String())
)

// withCommit appends the commit as semver build metadata: 0.5.3+abc1234.
func withCommit(v string) string {
	if commit == "" {
		return v
	}
	return v + "+" + commit
}
