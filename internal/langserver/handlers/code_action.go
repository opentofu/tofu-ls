// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"fmt"

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
	// We only support format type code actions at the moment, and do not want to format without the client asking for
	// them, so exit early here if nothing is requested.
	if len(params.Context.Only) == 0 {
		svc.logger.Printf("No code action requested, exiting")
		return ca, nil
	}

	for _, o := range params.Context.Only {
		svc.logger.Printf("Code actions requested: %q", o)
	}

	wantedCodeActions := ilsp.SupportedCodeActions.Only(params.Context.Only)
	if len(wantedCodeActions) == 0 {
		return nil, fmt.Errorf("could not find a supported code action to execute for %s, wanted %v",
			params.TextDocument.URI, params.Context.Only)
	}

	svc.logger.Printf("Code actions supported: %v", wantedCodeActions)

	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)

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
