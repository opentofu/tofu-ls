// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ast

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

func TestModuleDiags_autoloadedOnly(t *testing.T) {
	md := ModDiagsFromMap(map[string]hcl.Diagnostics{
		"alpha.tf": {},
		"beta.tf": {
			{
				Severity: hcl.DiagError,
				Summary:  "Test error",
				Detail:   "Test description",
			},
		},
		".hidden.tf": {},
	})
	diags := md.AutoloadedOnly().AsMap()
	expectedDiags := map[string]hcl.Diagnostics{
		"alpha.tf": {},
		"beta.tf": {
			{
				Severity: hcl.DiagError,
				Summary:  "Test error",
				Detail:   "Test description",
			},
		},
	}

	if diff := cmp.Diff(expectedDiags, diags, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("unexpected diagnostics: %s", diff)
	}
}

func TestIsModuleFilename(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"main.tf", true},
		{"main.tofu", true},
		{"main.tf.json", true},
		{"main.tofu.json", true},
		{"main.txt", false},
		{"main.json", false},
		{".hidden.tf", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := IsModuleFilename(tt.filename)
			if got != tt.expected {
				t.Errorf("IsModuleFilename(%q) = %v, want %v", tt.filename, got, tt.expected)
			}
		})
	}
}
