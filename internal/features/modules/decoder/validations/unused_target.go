// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"context"
	"slices"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	"github.com/hashicorp/hcl/v2"
)

// This is based slightly on the unused_origin alongside this file

func UnusedTargets(ctx context.Context, pathCtx *decoder.PathContext) lang.DiagnosticsMap {
	diagsMap := make(lang.DiagnosticsMap)

	for _, target := range pathCtx.ReferenceTargets {
		// Skip anything that doesn't have a declaration range
		if target.RangePtr == nil {
			continue
		}

		targetAddrType := target.Addr[0].String()
		if !slices.Contains([]string{"var", "local"}, targetAddrType) {
			continue
		}

		isUsed := false
		for _, origin := range pathCtx.ReferenceOrigins {
			localOrigin, ok := origin.(reference.LocalOrigin)
			if !ok {
				continue
			}

			originAddr := localOrigin.Address()
			if len(originAddr) != len(target.Addr) {
				if originAddr[0].String() != target.Addr[0].String() && originAddr[1].String() != target.Addr[1].String() {
					continue
				}
				isUsed = true
				break
			}
		}

		if !isUsed {
			file := target.RangePtr.Filename
			diagRange := target.RangePtr
			if target.DefRangePtr != nil {
				diagRange = target.DefRangePtr
			}

			d := &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Unused",
				Detail:   "This may not referenced anywhere in the configuration.",
				Subject:  diagRange,
			}
			diagsMap[file] = append(diagsMap[file], d)
		}
	}
	return diagsMap
}
