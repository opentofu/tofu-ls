// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/langserver/cmd"
	"github.com/opentofu/tofu-ls/internal/langserver/errors"
	"github.com/opentofu/tofu-ls/internal/langserver/progress"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
	"github.com/opentofu/tofu-ls/internal/tofu/module"
	"github.com/opentofu/tofu-ls/internal/uri"
)

func (h *CmdHandler) TofuInitHandler(ctx context.Context, args cmd.CommandArgs) (any, error) {
	dirURI, err := validateURI(args)
	if err != nil {
		return nil, err
	}

	dirHandle := document.DirHandleFromURI(dirURI)
	tfExec, err := module.TofuExecutorForModule(ctx, dirHandle.Path())
	if err != nil {
		return nil, errors.EnrichTfExecError(err)
	}

	progress.Begin(ctx, "Initializing")
	defer func() {
		progress.End(ctx, "Finished")
	}()

	progress.Report(ctx, "Running tofu init ...")
	err = tfExec.Init(ctx)
	if err != nil {
		return nil, err
	}

	// Clear diagnostics for open documents in this module so stale warnings disappear, the next time the file is edited
	// or saved, diagnostics will be re-computed.
	if h.Server != nil && h.StateStore != nil {
		docs, _ := h.StateStore.DocumentStore.OpenDocumentsForDir(dirHandle)
		for _, doc := range docs {
			docURI := doc.Dir.URI + "/" + doc.Filename
			_ = h.Server.Notify(ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
				URI:         lsp.DocumentURI(docURI),
				Diagnostics: []lsp.Diagnostic{},
			})
		}
	}

	return nil, nil
}

func validateURI(args cmd.CommandArgs) (string, error) {
	dirURI, ok := args.GetString("uri")
	if !ok || dirURI == "" {
		return "", fmt.Errorf("%w: expected module uri argument to be set", jrpc2.InvalidParams.Err())
	}

	if !uri.IsURIValid(dirURI) {
		return "", fmt.Errorf("URI %q is not valid", dirURI)
	}

	return dirURI, nil
}
