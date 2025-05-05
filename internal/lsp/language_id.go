// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

// LanguageID represents the coding language
// of a file
type LanguageID string

const (
	Terraform LanguageID = "opentofu"
	Tfvars    LanguageID = "opentofu-vars"
)

func (l LanguageID) String() string {
	return string(l)
}
