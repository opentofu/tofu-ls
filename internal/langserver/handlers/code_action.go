// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/opentofu/tofu-ls/internal/langserver/errors"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
	"github.com/opentofu/tofu-ls/internal/tofu/module"
)

func (svc *service) TextDocumentCodeAction(ctx context.Context, params lsp.CodeActionParams) []lsp.CodeAction {
	ca, err := svc.textDocumentCodeAction(ctx, params)
	if err != nil {
		svc.logger.Printf("code action failed: %s", err)
	}

	return ca
}

func (svc *service) textDocumentCodeAction(ctx context.Context, params lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	var ca []lsp.CodeAction

	// For action definitions, refer to https://code.visualstudio.com/api/references/vscode-api#CodeActionKind
	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)

	// Always check for quickfix actions when there are diagnostics, even if no
	// explicit request via Only
	if len(params.Context.Diagnostics) > 0 {
		for _, diag := range params.Context.Diagnostics {
			// Check for "Module schema not loaded" warning to suggest "tofu init"
			if diag.Severity == lsp.SeverityWarning &&
				strings.HasPrefix(diag.Message, "Module schema not loaded") {
				// Arguments must be in "key=value" string format
				arg, err := json.Marshal("uri=" + dh.Dir.URI)
				if err != nil {
					continue
				}
				ca = append(ca, lsp.CodeAction{
					Title:       "Run 'tofu init'",
					Kind:        ilsp.QuickFixTofuInit,
					Diagnostics: []lsp.Diagnostic{diag},
					Command: &lsp.Command{
						Title:     "Run 'tofu init'",
						Command:   "tofu-ls.tofu.init",
						Arguments: []json.RawMessage{arg},
					},
				})
				break // Only add one action
			}
		}
	}

	// For source actions, require explicit request via Only
	if len(params.Context.Only) == 0 {
		return ca, nil
	}

	for _, o := range params.Context.Only {
		svc.logger.Printf("Code actions requested: %q", o)
	}

	wantedCodeActions := ilsp.SupportedCodeActions.Only(params.Context.Only)
	if len(wantedCodeActions) == 0 {
		// No matching source actions, but we may have quickfix actions already
		return ca, nil
	}

	svc.logger.Printf("Code actions supported: %v", wantedCodeActions)

	doc, err := svc.stateStore.DocumentStore.GetDocument(dh)
	if err != nil {
		return ca, err
	}

	for action := range wantedCodeActions {
		switch action {
		case ilsp.SourceFormatAllTofu:
			tfExec, err := module.TofuExecutorForModule(ctx, dh.Dir.Path())
			if err != nil {
				return ca, errors.EnrichTfExecError(err)
			}

			edits, err := svc.formatDocument(ctx, tfExec, doc.Text, dh)
			if err != nil {
				return ca, err
			}

			ca = append(ca, lsp.CodeAction{
				Title: "Format Document",
				Kind:  action,
				Edit: lsp.WorkspaceEdit{
					Changes: map[lsp.DocumentURI][]lsp.TextEdit{
						lsp.DocumentURI(dh.FullURI()): edits,
					},
				},
			})
		}
	}

	return ca, nil
}
