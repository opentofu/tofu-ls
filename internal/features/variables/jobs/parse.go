// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jobs

import (
	"context"
	"path/filepath"

	lsctx "github.com/opentofu/tofu-ls/internal/context"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/features/variables/ast"
	"github.com/opentofu/tofu-ls/internal/features/variables/parser"
	"github.com/opentofu/tofu-ls/internal/features/variables/state"
	"github.com/opentofu/tofu-ls/internal/job"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	globalAst "github.com/opentofu/tofu-ls/internal/terraform/ast"
	op "github.com/opentofu/tofu-ls/internal/terraform/module/operation"
	"github.com/opentofu/tofu-ls/internal/uri"
)

// ParseVariables parses the variables configuration,
// i.e. turns bytes of `*.tfvars` files into AST ([*hcl.File]).
func ParseVariables(ctx context.Context, fs ReadOnlyFS, varStore *state.VariableStore, modPath string) error {
	mod, err := varStore.VariableRecordByPath(modPath)
	if err != nil {
		return err
	}

	// TODO: Avoid parsing if the content matches existing AST

	// Avoid parsing if it is already in progress or already known
	if mod.VarsDiagnosticsState[globalAst.HCLParsingSource] != op.OpStateUnknown && !job.IgnoreState(ctx) {
		return job.StateNotChangedErr{Dir: document.DirHandleFromPath(modPath)}
	}

	var files ast.VarsFiles
	var diags ast.VarsDiags
	rpcContext := lsctx.DocumentContext(ctx)
	// Only parse the file that's being changed/opened, unless this is 1st-time parsing
	if mod.VarsDiagnosticsState[globalAst.HCLParsingSource] == op.OpStateLoaded && rpcContext.IsDidChangeRequest() && rpcContext.LanguageID == ilsp.Tfvars.String() {
		// the file has already been parsed, so only examine this file and not the whole module
		err = varStore.SetVarsDiagnosticsState(modPath, globalAst.HCLParsingSource, op.OpStateLoading)
		if err != nil {
			return err
		}
		filePath, err := uri.PathFromURI(rpcContext.URI)
		if err != nil {
			return err
		}
		fileName := filepath.Base(filePath)

		f, vDiags, err := parser.ParseVariableFile(fs, filePath)
		if err != nil {
			return err
		}

		existingFiles := mod.ParsedVarsFiles.Copy()
		existingFiles[ast.VarsFilename(fileName)] = f
		files = existingFiles

		existingDiags, ok := mod.VarsDiagnostics[globalAst.HCLParsingSource]
		if !ok {
			existingDiags = make(ast.VarsDiags)
		} else {
			existingDiags = existingDiags.Copy()
		}
		existingDiags[ast.VarsFilename(fileName)] = vDiags
		diags = existingDiags
	} else {
		// this is the first time file is opened so parse the whole module
		err = varStore.SetVarsDiagnosticsState(modPath, globalAst.HCLParsingSource, op.OpStateLoading)
		if err != nil {
			return err
		}

		files, diags, err = parser.ParseVariableFiles(fs, modPath)
	}

	if err != nil {
		return err
	}

	sErr := varStore.UpdateParsedVarsFiles(modPath, files, err)
	if sErr != nil {
		return sErr
	}

	sErr = varStore.UpdateVarsDiagnostics(modPath, globalAst.HCLParsingSource, diags)
	if sErr != nil {
		return sErr
	}

	return err
}
