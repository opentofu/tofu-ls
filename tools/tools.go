// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools
// +build tools

package tools

import (
	_ "github.com/vektra/mockery/v2"
	_ "golang.org/x/tools/cmd/stringer"
)
