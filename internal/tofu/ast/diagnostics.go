// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ast

import (
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
)

// DiagnosticSource differentiates different sources of diagnostics.
type DiagnosticSource int

const (
	HCLParsingSource DiagnosticSource = iota
	SchemaValidationSource
	ReferenceValidationSource
	TofuValidateSource
	UnusedDeclarationSource
)

func (d DiagnosticSource) String() string {
	switch d {
	case UnusedDeclarationSource:
		// We want to be able to attach a friendly source here so that we can provide
		// Quick fixes specific to unused declarations.
		return "OpenTofu: Unused Declarations"
	default:
		return "OpenTofu"
	}
}

type DiagnosticSourceState map[DiagnosticSource]op.OpState

func (dss DiagnosticSourceState) Copy() DiagnosticSourceState {
	newDiagnosticSourceState := make(DiagnosticSourceState, len(dss))
	for source, state := range dss {
		newDiagnosticSourceState[source] = state
	}

	return newDiagnosticSourceState
}
