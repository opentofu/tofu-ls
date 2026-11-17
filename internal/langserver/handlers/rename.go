// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
	"github.com/opentofu/tofu-ls/internal/uri"
)

var errNotRenamable = errors.New("only resource, data, variable, local and module names can be renamed")

// renameSymbol is the definition of a renamable name.
type renameSymbol struct {
	name string
	// nameIdx is the index of the name in reference traversals, e.g. 1 for var.NAME, 2 for data.TYPE.NAME.
	nameIdx int
	// nameRange is the name at its definition, without quotes.
	nameRange hcl.Range
	// outerRange is nameRange including quotes; the cursor may sit anywhere in it.
	outerRange hcl.Range
	// blockStart is a position inside the definition, used to look up its references.
	blockStart hcl.Pos
}

// findDefinition returns the renamable name defined at pos, if any.
func findDefinition(src []byte, filename string, pos hcl.Pos) (*renameSymbol, bool) {
	file, diags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if file == nil || file.Body == nil {
		_ = diags
		return nil, false
	}
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, false
	}

	for _, b := range body.Blocks {
		if !b.Range().ContainsPos(pos) {
			continue
		}
		var labelIdx, nameIdx int
		switch b.Type {
		case "variable", "module":
			labelIdx, nameIdx = 0, 1
		case "resource":
			labelIdx, nameIdx = 1, 1
		case "data":
			labelIdx, nameIdx = 1, 2
		case "locals":
			for _, attr := range b.Body.Attributes {
				if attr.Range().ContainsPos(pos) {
					return &renameSymbol{
						name: attr.Name, nameIdx: 1,
						nameRange: attr.NameRange, outerRange: attr.NameRange,
						blockStart: attr.NameRange.Start,
					}, true
				}
			}
			return nil, false
		default:
			return nil, false
		}
		if len(b.Labels) <= labelIdx {
			return nil, false
		}
		outer := b.LabelRanges[labelIdx]
		inner := outer
		if src[outer.Start.Byte] == '"' && outer.End.Byte-outer.Start.Byte >= 2 {
			inner.Start.Column++
			inner.Start.Byte++
			inner.End.Column--
			inner.End.Byte--
		}
		return &renameSymbol{
			name: b.Labels[labelIdx], nameIdx: nameIdx,
			nameRange: inner, outerRange: outer,
			blockStart: b.TypeRange.Start,
		}, true
	}
	return nil, false
}

// nameInTraversal locates the name step of the traversal covering origin.
func nameInTraversal(src []byte, origin hcl.Range, sym *renameSymbol) (hcl.Range, bool) {
	if origin.End.Byte > len(src) {
		return hcl.Range{}, false
	}
	trav, diags := hclsyntax.ParseTraversalAbs(src[origin.Start.Byte:origin.End.Byte], origin.Filename, origin.Start)
	if diags.HasErrors() || len(trav) <= sym.nameIdx {
		return hcl.Range{}, false
	}
	step, ok := trav[sym.nameIdx].(hcl.TraverseAttr)
	if !ok || step.Name != sym.name {
		return hcl.Range{}, false
	}
	rng := step.SrcRange // includes the leading dot
	rng.Start.Column++
	rng.Start.Byte++
	return rng, true
}

// resolveRename finds the definition behind the symbol at the cursor, plus
// the path it lives in and the range of the name under the cursor.
func (svc *service) resolveRename(ctx context.Context, params lsp.TextDocumentPositionParams) (*renameSymbol, lang.Path, string, lsp.Range, error) {
	var none lang.Path

	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)
	doc, err := svc.stateStore.DocumentStore.GetDocument(dh)
	if err != nil {
		return nil, none, "", lsp.Range{}, err
	}

	jobIds, err := svc.stateStore.JobStore.ListIncompleteJobsForDir(dh.Dir)
	if err != nil {
		return nil, none, "", lsp.Range{}, err
	}
	svc.stateStore.JobStore.WaitForJobs(ctx, jobIds...)

	pos, err := ilsp.HCLPositionFromLspPosition(params.Position, doc)
	if err != nil {
		return nil, none, "", lsp.Range{}, err
	}
	path := lang.Path{
		Path:       doc.Dir.Path(),
		LanguageID: ilsp.ParseLanguageID(doc.LanguageID).String(),
	}

	// Cursor on the definition's name.
	if sym, ok := findDefinition(doc.Text, doc.Filename, pos); ok && sym.outerRange.ContainsPos(pos) {
		return sym, path, doc.Filename, ilsp.HCLRangeToLSP(sym.nameRange), nil
	}

	// Cursor on a reference: follow it to the definition.
	targets, err := svc.decoder.ReferenceTargetsForOriginAtPos(path, doc.Filename, pos)
	if err != nil || len(targets) == 0 {
		return nil, none, "", lsp.Range{}, errNotRenamable
	}
	t := targets[0]
	src, err := svc.fs.ReadFile(filepath.Join(t.Path.Path, t.Range.Filename))
	if err != nil {
		return nil, none, "", lsp.Range{}, err
	}
	sym, ok := findDefinition(src, t.Range.Filename, t.Range.Start)
	if !ok {
		return nil, none, "", lsp.Range{}, errNotRenamable
	}
	rng, ok := nameInTraversal(doc.Text, t.OriginRange, sym)
	if !ok {
		return nil, none, "", lsp.Range{}, errNotRenamable
	}
	return sym, t.Path, t.Range.Filename, ilsp.HCLRangeToLSP(rng), nil
}

func (svc *service) PrepareRename(ctx context.Context, params lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	sym, _, _, rng, err := svc.resolveRename(ctx, params.TextDocumentPositionParams)
	if err != nil {
		return nil, err
	}
	return &lsp.PrepareRenameResult{Range: rng, Placeholder: sym.name}, nil
}

func (svc *service) Rename(ctx context.Context, params lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	if !hclsyntax.ValidIdentifier(params.NewName) {
		return nil, fmt.Errorf("%q is not a valid identifier", params.NewName)
	}

	sym, defPath, defFile, _, err := svc.resolveRename(ctx, lsp.TextDocumentPositionParams{
		TextDocument: params.TextDocument,
		Position:     params.Position,
	})
	if err != nil {
		return nil, err
	}

	edit := &lsp.WorkspaceEdit{Changes: map[lsp.DocumentURI][]lsp.TextEdit{}}
	if params.NewName == sym.name {
		return edit, nil
	}
	add := func(dir, file string, rng hcl.Range) {
		u := lsp.DocumentURI(uri.FromPath(filepath.Join(dir, file)))
		edit.Changes[u] = append(edit.Changes[u], lsp.TextEdit{Range: ilsp.HCLRangeToLSP(rng), NewText: params.NewName})
	}

	add(defPath.Path, defFile, sym.nameRange)

	origins := svc.decoder.ReferenceOriginsTargetingPos(defPath, defFile, sym.blockStart)
	for _, o := range origins {
		src, err := svc.fs.ReadFile(filepath.Join(o.Path.Path, o.Range.Filename))
		if err != nil {
			return nil, err
		}
		if rng, ok := nameInTraversal(src, o.Range, sym); ok {
			add(o.Path.Path, o.Range.Filename, rng)
		}
	}
	return edit, nil
}
