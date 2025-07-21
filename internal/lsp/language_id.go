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
	// Some editors do not support language ID overrides which makes it difficult to use this language server
	// We also need to accept language IDs of Terraform in order to circumvent this issue
	Terraform     LanguageID = "terraform"
	TerraformVars LanguageID = "terraform-vars"
)

// LanguageIDsMatch checks if both IDs are of configuration language or vars language
func LanguageIDsMatch(a, b string) bool {
	return a == b || (a == OpenTofu.String() && b == Terraform.String()) || (a == Terraform.String() && b == OpenTofu.String())
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
