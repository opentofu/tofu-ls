// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

func TestHCLDiagsToLSP_NeverReturnsNil(t *testing.T) {
	diags := HCLDiagsToLSP(nil, "test")
	if diags == nil {
		t.Fatal("diags should not be nil")
	}

	diags = HCLDiagsToLSP(hcl.Diagnostics{}, "test")
	if diags == nil {
		t.Fatal("diags should not be nil")
	}

	diags = HCLDiagsToLSP(hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
		},
	}, "source")
	if diags == nil {
		t.Fatal("diags should not be nil")
	}
}

func TestHCLDiagsToLSP_Metadata(t *testing.T) {
	tests := []struct {
		name     string
		hclDiags hcl.Diagnostics
		want     []lsp.Diagnostic
	}{
		{
			name: "diag warning and no metadata",
			hclDiags: []*hcl.Diagnostic{
				{
					Severity: hcl.DiagWarning,
					Summary:  "warning",
				},
			},
			want: []lsp.Diagnostic{
				{
					Severity: lsp.SeverityWarning,
					Source:   "source",
					Message:  "warning",
				},
			},
		},
		{
			name: "diag hint and tag present",
			hclDiags: []*hcl.Diagnostic{
				{
					Severity: hcl.DiagWarning,
					Summary:  "warning",
					Extra:    DiagnosticMetadata{Hint: true, Unnecessary: true},
				},
			},
			want: []lsp.Diagnostic{
				{
					Severity: lsp.SeverityHint,
					Source:   "source",
					Message:  "warning",
					Tags:     []lsp.DiagnosticTag{lsp.Unnecessary},
				},
			},
		},
		{
			name: "diag hint and no tag present",
			hclDiags: []*hcl.Diagnostic{
				{
					Severity: hcl.DiagWarning,
					Summary:  "warning",
					Extra:    DiagnosticMetadata{Hint: true},
				},
			},
			want: []lsp.Diagnostic{
				{
					Severity: lsp.SeverityHint,
					Source:   "source",
					Message:  "warning",
				},
			},
		},
		{
			name: "diag warning and tag present",
			hclDiags: []*hcl.Diagnostic{
				{
					Severity: hcl.DiagWarning,
					Summary:  "warning",
					Extra:    DiagnosticMetadata{Unnecessary: true},
				},
			},
			want: []lsp.Diagnostic{
				{
					Severity: lsp.SeverityWarning,
					Source:   "source",
					Message:  "warning",
					Tags:     []lsp.DiagnosticTag{lsp.Unnecessary},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%2d-%s", i, tt.name), func(t *testing.T) {
			got := HCLDiagsToLSP(tt.hclDiags, "source")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("HCLDiagsToLSP() -want +got: %s", diff)
			}
		})
	}
}
