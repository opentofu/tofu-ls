// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	"github.com/hashicorp/hcl/v2"
)

// This is based slightly on the unused_origin alongside this file

// UnusedDeclarationExtra carries the full block range for quick fix code actions
type UnusedDeclarationExtra struct {
	FullRange hcl.Range
}

// FullBlockRange returns the full block range for removal
func (u UnusedDeclarationExtra) FullBlockRange() *hcl.Range {
	return &u.FullRange
}

func UnusedTargets(ctx context.Context, pathCtx *decoder.PathContext) lang.DiagnosticsMap {
	diagsMap := make(lang.DiagnosticsMap)

	// Track seen targets to avoid duplicates (PathContext may contain duplicate targets)
	seenTargets := make(map[string]bool)

	for _, target := range pathCtx.ReferenceTargets {
		// Skip anything that doesn't have a declaration range
		if target.RangePtr == nil {
			continue
		}

		targetAddrType := target.Addr[0].String()
		if !slices.Contains([]string{"var", "local"}, targetAddrType) {
			continue
		}

		// Build a unique key for this target (e.g., "var.foo" or "local.bar")
		targetKey := target.Addr.String()
		if seenTargets[targetKey] {
			continue
		}
		seenTargets[targetKey] = true

		isUsed := false
		for _, origin := range pathCtx.ReferenceOrigins {
			localOrigin, ok := origin.(reference.LocalOrigin)
			if !ok {
				continue
			}

			originAddr := localOrigin.Address()
			if len(originAddr) >= 2 && len(target.Addr) >= 2 &&
				originAddr[0].String() == target.Addr[0].String() &&
				originAddr[1].String() == target.Addr[1].String() {
				isUsed = true
				break
			}
		}

		if !isUsed {
			file := target.RangePtr.Filename

			d := &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  fmt.Sprintf("unused: %q", target.Addr.String()),
				Subject:  target.RangePtr,
				Extra:    UnusedDeclarationExtra{FullRange: *target.RangePtr},
			}
			diagsMap[file] = append(diagsMap[file], d)
		}
	}
	return diagsMap
}
