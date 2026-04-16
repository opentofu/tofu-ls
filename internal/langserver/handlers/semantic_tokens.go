// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/creachadair/jrpc2"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

func (svc *service) TextDocumentSemanticTokensFull(ctx context.Context, params lsp.SemanticTokensParams) (lsp.SemanticTokens, error) {
	tks := lsp.SemanticTokens{}

	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return tks, err
	}

	caps := ilsp.SemanticTokensClientCapabilities{
		SemanticTokensClientCapabilities: cc.TextDocument.SemanticTokens,
	}
	if !caps.FullRequest() {
		// This would indicate a buggy client which sent a request
		// it didn't claim to support, so we just strictly follow
		// the protocol here and avoid serving buggy clients.
		svc.logger.Printf("semantic tokens full request support not announced by client")
		return tks, jrpc2.MethodNotFound.Err()
	}

	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)
	doc, err := svc.stateStore.DocumentStore.GetDocument(dh)
	if err != nil {
		return tks, err
	}

	jobIds, err := svc.stateStore.JobStore.ListIncompleteJobsForDir(dh.Dir)
	if err != nil {
		return tks, err
	}
	svc.stateStore.JobStore.WaitForJobs(ctx, jobIds...)

	d, err := svc.decoderForDocument(ctx, doc)
	if err != nil {
		return tks, err
	}

	tokens, err := d.SemanticTokensInFile(ctx, doc.Filename)
	if err != nil {
		return tks, err
	}

	// TODO remove this once finished. Meant only for debugging
	var b strings.Builder
	for _, token := range tokens {
		var mods []string
		for _, modifier := range token.Modifiers {
			mods = append(mods, string(modifier))
		}
		b.WriteString(fmt.Sprintf("%s; mod: %s; rangeStart: %d-%d-%d; rangeEnd: %d-%d-%d ### ",
			token.Type,
			strings.Join(mods, " "),
			token.Range.Start.Line,
			token.Range.Start.Column,
			token.Range.Start.Byte,
			token.Range.End.Line,
			token.Range.End.Column,
			token.Range.End.Byte,
		))
	}
	svc.logger.Println(b.String())
	te := &ilsp.TokenEncoder{
		Lines:      doc.Lines,
		Tokens:     tokens,
		ClientCaps: cc.TextDocument.SemanticTokens,
	}
	tks.Data = te.Encode()

	return tks, nil
}
