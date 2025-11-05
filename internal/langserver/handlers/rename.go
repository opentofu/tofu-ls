// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/opentofu/tofu-ls/internal/document"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	"github.com/hashicorp/hcl/v2"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

// TODO If a client neither supports `documentChanges` nor `workspace.workspaceEdit.resourceOperations` then only plain `TextEdit`s using the `changes` property are supported.
// For now this is a naive implementation without versioning using Changes property
// Also, the current implementation won't work if there are no references for the origin
func (svc *service) Rename(ctx context.Context, params lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	//return nil, nil
	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)
	doc, err := svc.stateStore.DocumentStore.GetDocument(dh)
	if err != nil {
		return nil, err
	}

	pos, err := ilsp.HCLPositionFromLspPosition(params.Position, doc)
	if err != nil {
		return nil, err
	}

	jobIds, err := svc.stateStore.JobStore.ListIncompleteJobsForDir(dh.Dir)
	if err != nil {
		return nil, err
	}
	svc.stateStore.JobStore.WaitForJobs(ctx, jobIds...)

	path := lang.Path{
		Path:       doc.Dir.Path(),
		LanguageID: doc.LanguageID,
	}

	localCtx, err := svc.pathReader.PathContext(path)
	if err != nil {
		return nil, err
	}

	var declaration *reference.Target
	var referenceOrigin *reference.Origin
	// Used to adjust the range start for renames, to only rename the attribute name part of the reference
	var originalNameLength int

	// Check if this position contains reference targets first, if it does, we assume the declaration we find is the target
	// Small caveat here, this might return multiple declarations and that needs to be handled in the flushed out version //TODO
	if refTargets, ok := localCtx.ReferenceTargets.InnermostAtPos(doc.Filename, pos); ok {
		for _, ref := range refTargets {
			if !ref.DefRangePtr.ContainsPos(pos) {
				continue
			}
			declaration = &ref
			originalNameLength = len(ref.Addr[len(ref.Addr)-1].String()[1:])
			break
			// TODO Is there something we need to do with other matches?
		}
	}

	if refOrigins, ok := localCtx.ReferenceOrigins.AtPos(doc.Filename, pos); ok {
		for _, ref := range refOrigins {
			if morig, ok := ref.(reference.MatchableOrigin); ok {
				referenceOrigin = &ref
				addr := morig.Address()
				originalNameLength = len(addr[len(addr)-1].String()[1:])
			}
			// TODO Is there something we need to do with other matches?
		}
	}
	//TODO handle no declaration or referenceOrigin found, likely a singular declaration with no reference

	if declaration == nil && referenceOrigin == nil {
		return nil, fmt.Errorf("could not resolve reference")
	}

	rangesToEditInDocs := make(map[string]map[hcl.Range]string)
	// Build ranges to edit from declaration
	if declaration != nil && declaration.DefRangePtr != nil {
		rng := *declaration.DefRangePtr
		err := svc.buildReferencesFromTarget(&decoder.ReferenceTarget{
			OriginRange: rng,
			Path: lang.Path{
				Path:       doc.Dir.Path(),
				LanguageID: doc.LanguageID,
			},
			Range:       rng,
			DefRangePtr: &rng,
		}, params.NewName, rangesToEditInDocs)
		if err != nil {
			return nil, err
		}
	}

	// In this case we have a reference as the starting point for rename
	// We will work backwards to find all targets and then all origins for those targets
	if referenceOrigin != nil {
		// Get all targets for this origin
		targets, err := svc.decoder.ReferenceTargetsForOriginAtPos(path, doc.Filename, pos)
		if err != nil {
			return nil, err
		}
		for _, target := range targets {
			err := svc.buildReferencesFromTarget(target, params.NewName, rangesToEditInDocs)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(rangesToEditInDocs) == 0 {
		return nil, fmt.Errorf("no references found for rename")
	}

	workspaceEdit := &lsp.WorkspaceEdit{
		Changes: make(map[lsp.DocumentURI][]lsp.TextEdit),
	}
	for docPath, ranges := range rangesToEditInDocs {
		dh := document.HandleFromPath(docPath)
		docURI := lsp.DocumentURI(dh.FullURI())
		edits := make([]lsp.TextEdit, 0, len(ranges))
		for r, s := range ranges {
			if r.Start.Line != r.End.Line {
				// Sanity check to avoid bad ranges
				panic("unexpected range for rename, multi-line ranges not supported")
			}
			r.Start.Column = r.End.Column - originalNameLength
			edits = append(edits, lsp.TextEdit{
				Range:   ilsp.HCLRangeToLSP(r),
				NewText: s,
			})
		}
		workspaceEdit.Changes[docURI] = edits
	}

	return workspaceEdit, nil
}

func (svc *service) buildReferencesFromTarget(refTarget *decoder.ReferenceTarget, newText string, toEdit map[string]map[hcl.Range]string) error {
	//TODO add handling for the target too
	path := lang.Path{
		Path:       refTarget.Path.Path,
		LanguageID: refTarget.Path.LanguageID,
	}
	// Get all origins for this declaration
	origins := svc.decoder.ReferenceOriginsTargetingPos(path, refTarget.Range.Filename, refTarget.Range.Start)
	for _, origin := range origins {
		docURI := origin.Path.Path + "/" + origin.Range.Filename
		if _, ok := toEdit[docURI]; !ok {
			toEdit[docURI] = make(map[hcl.Range]string)
		}
		toEdit[docURI][origin.Range] = newText
	}

	targetURI := refTarget.Path.Path + "/" + refTarget.Range.Filename
	defRange := *refTarget.DefRangePtr
	localCtx, err := svc.pathReader.PathContext(path)
	if err != nil {
		return err
	}
	var detailedTarget *reference.Target
	// Get reference target for the declaration we have
	if refTargets, ok := localCtx.ReferenceTargets.InnermostAtPos(defRange.Filename, defRange.Start); ok {
		for _, ref := range refTargets {
			//if !ref.DefRangePtr.ContainsPos(pos) {
			//	continue
			//}
			detailedTarget = &ref
			break
			// TODO Is there something we need to do with other matches?
		}
	}
	if detailedTarget == nil {
		return fmt.Errorf("could not resolve reference")
	}
	// specific case here for variable-like declaration - shifting the end of the range to dodge the last "
	if detailedTarget.ScopeId != "local" {
		defRange.End.Column = defRange.End.Column - 1
	}
	if _, ok := toEdit[targetURI]; !ok {
		toEdit[targetURI] = make(map[hcl.Range]string)
	}
	toEdit[targetURI][defRange] = newText
	return nil
}

func (svc *service) PrepareRename(ctx context.Context, params lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	dh := ilsp.HandleFromDocumentURI(params.TextDocument.URI)
	doc, err := svc.stateStore.DocumentStore.GetDocument(dh)
	if err != nil {
		return nil, err
	}

	pos, err := ilsp.HCLPositionFromLspPosition(params.Position, doc)
	if err != nil {
		return nil, err
	}

	jobIds, err := svc.stateStore.JobStore.ListIncompleteJobsForDir(dh.Dir)
	if err != nil {
		return nil, err
	}
	svc.stateStore.JobStore.WaitForJobs(ctx, jobIds...)

	path := lang.Path{
		Path:       doc.Dir.Path(),
		LanguageID: doc.LanguageID,
	}

	localCtx, err := svc.pathReader.PathContext(path)
	if err != nil {
		return nil, err
	}

	var symbolRange hcl.Range
	var symbolAddr lang.Address

	// TODO validations to omit references that can't be renamed (resource props, etc)
	if refOrigins, ok := localCtx.ReferenceOrigins.AtPos(doc.Filename, pos); ok {
		for _, ref := range refOrigins {
			if localRef, ok := ref.(reference.MatchableOrigin); ok {
				symbolRange = localRef.OriginRange()
				symbolAddr = localRef.Address()
				break
				// TODO Is there something we need to do with other matches?
			}
		}
	} else if refTargets, ok := localCtx.ReferenceTargets.InnermostAtPos(doc.Filename, pos); ok {
		for _, ref := range refTargets {
			symbolRange = *ref.RangePtr // This should always be non-nil since we just found this with pos
			symbolAddr = ref.Addr
			break
			// TODO Is there something we need to do with other matches?
		}
	} else {
		return nil, fmt.Errorf("could not find symbol reference")
	}

	symbolName := symbolAddr[len(symbolAddr)-1].String()[1:]

	return &lsp.PrepareRenameResult{
		Range:       ilsp.HCLRangeToLSP(symbolRange),
		Placeholder: symbolName,
	}, nil
}

func (svc *service) RenameLinkedEditingRange(ctx context.Context, params lsp.LinkedEditingRangeParams) (*lsp.LinkedEditingRanges, error) {
	return nil, nil
}
