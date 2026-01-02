// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	tfmod "github.com/opentofu/opentofu-schema/module"
	"github.com/opentofu/opentofu-schema/registry"
	tfaddr "github.com/opentofu/registry-address"
)

// ModuleReader provides access to module state for validation
type ModuleReader interface {
	DeclaredModuleCalls(modPath string) (map[string]tfmod.DeclaredModuleCall, error)
	RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error)
	LocalModuleMeta(modPath string) (*tfmod.Meta, error)
}

// RootModuleReader provides access to root module state
type RootModuleReader interface {
	InstalledModulePath(rootPath string, normalizedSource string) (string, bool)
}

// UnloadedModuleSchemaResult contains diagnostics for modules without schema
// and the ranges of those modules (for filtering false positive errors).
type UnloadedModuleSchemaResult struct {
	// Diagnostics contains "Module schema not loaded" warnings per file
	Diagnostics lang.DiagnosticsMap
	// UnloadedRanges contains the ranges of module blocks without schema per file
	UnloadedRanges map[string][]hcl.Range
}

// ModuleSchemaNotLoaded checks for module blocks that don't have their schema loaded
// and returns diagnostics for each one, plus the ranges for filtering.
func ModuleSchemaNotLoaded(
	modReader ModuleReader,
	rootReader RootModuleReader,
	modPath string,
) UnloadedModuleSchemaResult {
	result := UnloadedModuleSchemaResult{
		Diagnostics:    make(lang.DiagnosticsMap),
		UnloadedRanges: make(map[string][]hcl.Range),
	}

	// Get declared module calls
	declaredCalls, err := modReader.DeclaredModuleCalls(modPath)
	if err != nil {
		return result
	}

	for _, call := range declaredCalls {
		if hasModuleSchemaLoaded(modReader, rootReader, modPath, call) {
			continue
		}

		// Module schema not loaded - emit diagnostic and track range
		if call.RangePtr != nil {
			fileName := call.RangePtr.Filename

			headerRange := makeRangeForLine(call.RangePtr.Filename, call.RangePtr.Start.Line)

			d := &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Module schema not loaded",
				Detail:   "Run 'tofu init' to download module metadata and enable full validation.",
				Subject:  headerRange.Ptr(),
			}
			result.Diagnostics[fileName] = result.Diagnostics[fileName].Append(d)

			result.UnloadedRanges[fileName] = append(result.UnloadedRanges[fileName], *call.RangePtr)
		}
	}

	return result
}

// FilterDiagnosticsForUnloadedModules removes "Unexpected attribute" diagnostics
// that fall within module blocks that don't have schema loaded.
func FilterDiagnosticsForUnloadedModules(diags lang.DiagnosticsMap, unloadedRanges map[string][]hcl.Range) lang.DiagnosticsMap {
	filtered := make(lang.DiagnosticsMap)

	for filename, fileDiags := range diags {
		ranges, hasUnloaded := unloadedRanges[filename]
		if !hasUnloaded {
			// No unloaded modules in this file, keep all diagnostics
			filtered[filename] = fileDiags
			continue
		}

		// Filter out false positives
		var kept hcl.Diagnostics
		for _, diag := range fileDiags {
			if ShouldFilterDiagnostic(diag, ranges) {
				continue
			}
			kept = append(kept, diag)
		}
		if len(kept) > 0 {
			filtered[filename] = kept
		}
	}

	return filtered
}

// ShouldFilterDiagnostic returns true if the diagnostic should be filtered out
// because it's a false positive from an unloaded module schema.
// For now this should only be unexpected attribute errors within unloaded module blocks.
func ShouldFilterDiagnostic(diag *hcl.Diagnostic, unloadedRanges []hcl.Range) bool {
	// Only filter "Unexpected attribute" errors
	if !strings.HasPrefix(diag.Summary, "Unexpected attribute") {
		return false
	}

	// Check if the diagnostic falls within an unloaded module block
	if diag.Subject == nil {
		return false
	}

	for _, moduleRange := range unloadedRanges {
		if rangeContains(moduleRange, *diag.Subject) {
			return true
		}
	}

	return false
}

// rangeContains checks if outer contains inner (same file, inner starts after outer starts, inner ends before outer ends)
func rangeContains(outer, inner hcl.Range) bool {
	if outer.Filename != inner.Filename {
		return false
	}

	// Check if inner start is after or at outer start
	if inner.Start.Line < outer.Start.Line {
		return false
	}
	if inner.Start.Line == outer.Start.Line && inner.Start.Column < outer.Start.Column {
		return false
	}

	// Check if inner end is before or at outer end
	if inner.End.Line > outer.End.Line {
		return false
	}
	if inner.End.Line == outer.End.Line && inner.End.Column > outer.End.Column {
		return false
	}

	return true
}

func hasModuleSchemaLoaded(
	modReader ModuleReader,
	rootReader RootModuleReader,
	modPath string,
	call tfmod.DeclaredModuleCall,
) bool {
	switch sourceAddr := call.SourceAddr.(type) {
	case tfaddr.Module:
		// Registry module - first check if installed locally
		installedDir, ok := rootReader.InstalledModulePath(modPath, sourceAddr.String())
		if ok {
			path := filepath.Join(modPath, installedDir)
			_, err := modReader.LocalModuleMeta(path)
			return err == nil
		}

		// Otherwise check if we have cached registry data
		_, err := modReader.RegistryModuleMeta(sourceAddr, call.Version)
		return err == nil

	case tfmod.RemoteSourceAddr:
		// Remote module (git, etc.) - check if installed
		installedDir, ok := rootReader.InstalledModulePath(modPath, sourceAddr.String())
		if !ok {
			return false
		}
		path := filepath.Join(modPath, installedDir)
		_, err := modReader.LocalModuleMeta(path)
		return err == nil

	case tfmod.LocalSourceAddr:
		// Local module - check if we can get its metadata
		path := filepath.Join(modPath, sourceAddr.String())
		_, err := modReader.LocalModuleMeta(path)
		return err == nil

	default:
		// Unknown source type - assume not loaded
		return false
	}
}

// makeRangeForLine creates an hcl.Range spanning an entire line (column 0 to end of line).
// Editors will clamp the end column to the actual line length.
func makeRangeForLine(filename string, line int) hcl.Range {
	return hcl.Range{
		Filename: filename,
		Start: hcl.Pos{
			Line:   line,
			Column: 0,
			Byte:   0,
		},
		End: hcl.Pos{
			Line:   line,
			Column: 9999,
			Byte:   9999,
		},
	}
}
