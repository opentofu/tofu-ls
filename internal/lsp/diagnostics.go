// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"github.com/hashicorp/hcl/v2"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

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

// DiagnosticMetadata holds metadata stored in the Extra field of hcl.Diagnostic.
// It can be used to provide extra information, such as Hint severity and tags.
type DiagnosticMetadata struct {
	Hint        bool
	Unnecessary bool
}

func HCLDiagsToLSP(hclDiags hcl.Diagnostics, source string) []lsp.Diagnostic {
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

		if metadata, ok := hclDiag.Extra.(DiagnosticMetadata); ok {
			if metadata.Hint {
				severity = lsp.SeverityHint
			}
			if metadata.Unnecessary {
				tags = append(tags, lsp.Unnecessary)
			}
		}
		diags = append(diags, lsp.Diagnostic{
			Range:    rnge,
			Severity: severity,
			Source:   source,
			Message:  msg,
			Tags:     tags,
		})

	}
	return diags
}
