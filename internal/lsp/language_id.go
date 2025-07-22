// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

// LanguageID represents the coding language
// of a file
type LanguageID string

const (
	OpenTofu     LanguageID = "opentofu"
	OpenTofuVars LanguageID = "opentofu-vars"
	// Terraform - Some editors do not support language ID overrides which makes it difficult to use this language server
	// We also need to accept language IDs of Terraform to circumvent this issue
	Terraform     LanguageID = "terraform"
	TerraformVars LanguageID = "terraform-vars"
)

// ParseLanguageID parses a string into a LanguageID
// We also remap Terraform to OpenTofu and TerraformVars to OpenTofuVars
// We assume that the language ID is valid or the validation step has been done before parsing
func ParseLanguageID(id string) LanguageID {
	switch LanguageID(id) {
	case Terraform:
		return OpenTofu
	case TerraformVars:
		return OpenTofuVars
	default:
		return LanguageID(id)
	}
}

func IsValidConfigLanguage(id string) bool {
	switch LanguageID(id) {
	case OpenTofu, Terraform:
		return true
	default:
		return false
	}
}

func IsValidVarsLanguage(id string) bool {
	switch LanguageID(id) {
	case OpenTofuVars, TerraformVars:
		return true
	default:
		return false
	}
}

func (l LanguageID) String() string {
	return string(l)
}
