// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ast

// IsRootModuleFilename checks if the given filename is a root module file.
// TODO: extend this for OpenTofu specific files
func IsRootModuleFilename(name string) bool {
	return (name == ".terraform.lock.hcl" ||
		name == ".terraform-version")
}
