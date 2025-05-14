// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ast

func IsRootModuleFilename(name string) bool {
	return (name == ".terraform.lock.hcl" ||
		name == ".terraform-version")
}
