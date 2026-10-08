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

func deprecatedVariableKey(filename, name string) string {
	return filename + "\x00" + name
}

func deprecatedVariableTargets(files map[string]*hcl.File) map[string]bool {
	targets := make(map[string]bool)
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{{
			Type:       "variable",
			LabelNames: []string{"name"},
		}},
	}

	for filename, file := range files {
		content, _, diags := file.Body.PartialContent(schema)
		if diags.HasErrors() {
			continue
		}

		for _, block := range content.Blocks {
			if len(block.Labels) != 1 {
				continue
			}

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				continue
			}
			// If we have the `deprecated` attribute then keep a note of this to return at the end
			if _, ok := attrs["deprecated"]; ok {
				targets[deprecatedVariableKey(filename, block.Labels[0])] = true
			}
		}
	}

	return targets
}

func UnusedTargets(ctx context.Context, pathCtx *decoder.PathContext) lang.DiagnosticsMap {
	diagsMap := make(lang.DiagnosticsMap)

	// getKey builds a lookup key from the first two address parts (e.g., "var.foo" or "local.bar")
	getKey := func(addr lang.Address) string {
		return addr[0].String() + "." + addr[1].String()
	}

	// Build a map of used origins for looking up later
	usedOrigins := make(map[string]bool)
	for _, origin := range pathCtx.ReferenceOrigins {
		localOrigin, ok := origin.(reference.LocalOrigin)
		if !ok {
			continue
		}

		originAddr := localOrigin.Address()
		if len(originAddr) >= 2 {
			usedOrigins[getKey(originAddr)] = true
		}
	}

	deprecatedVariables := deprecatedVariableTargets(pathCtx.Files)

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

		deprecated := false
		if targetAddrType == "var" {
			if name, ok := target.Addr[1].(lang.AttrStep); ok {
				deprecated = deprecatedVariables[deprecatedVariableKey(target.RangePtr.Filename, name.Name)]
			}
		}

		if !usedOrigins[getKey(target.Addr)] && !deprecated {
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
