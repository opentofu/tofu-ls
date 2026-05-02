// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"github.com/hashicorp/hcl/v2"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

// FullBlockRanger is implemented by diagnostic Extra data that provides
// the full block range for quick fix code actions
type FullBlockRanger interface {
	FullBlockRange() *hcl.Range
}

func HCLSeverityToLSP(severity hcl.DiagnosticSeverity) lsp.DiagnosticSeverity {
	var sev lsp.DiagnosticSeverity
	switch severity {
	case hcl.DiagError:
		sev = lsp.SeverityError
	case hcl.DiagWarning:
		sev = lsp.SeverityWarning
	case hcl.DiagInvalid:
		panic("invalid diagnostic")
	}
	return sev
}

type DiagnosticOptions struct {
	Severity *lsp.DiagnosticSeverity
	Tags     []lsp.DiagnosticTag
}

func HCLDiagsToLSP(hclDiags hcl.Diagnostics, source string, opts ...DiagnosticOptions) []lsp.Diagnostic {
	diags := []lsp.Diagnostic{}

	for _, hclDiag := range hclDiags {
		msg := hclDiag.Summary
		if hclDiag.Detail != "" {
			msg += ": " + hclDiag.Detail
		}
		var rnge lsp.Range
		if hclDiag.Subject != nil {
			rnge = HCLRangeToLSP(*hclDiag.Subject)
		}

		severity := HCLSeverityToLSP(hclDiag.Severity)
		var tags []lsp.DiagnosticTag

		// Apply options if provided
		if len(opts) > 0 {
			if opts[0].Severity != nil {
				severity = *opts[0].Severity
			}
			if opts[0].Tags != nil {
				tags = opts[0].Tags
			}
		}

		diag := lsp.Diagnostic{
			Range:    rnge,
			Severity: severity,
			Source:   source,
			Message:  msg,
			Tags:     tags,
		}

		// Check if Extra provides full block range for quick fixes
		if extra, ok := hclDiag.Extra.(FullBlockRanger); ok {
			if fullRange := extra.FullBlockRange(); fullRange != nil {
				diag.Data = map[string]interface{}{
					"fullRange": HCLRangeToLSP(*fullRange),
				}
			}
		}

		diags = append(diags, diag)
	}
	return diags
}
