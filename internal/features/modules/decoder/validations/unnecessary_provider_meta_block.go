// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"context"

	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl-lang/schemacontext"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
)

type UnnecessaryProviderMetaBlock struct{}

func (upmb UnnecessaryProviderMetaBlock) Visit(ctx context.Context, node hclsyntax.Node, nodeSchema schema.Schema) (context.Context, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	block, ok := node.(*hclsyntax.Block)
	if !ok || block.Type != "terraform" {
		return ctx, diags
	}

	if lvl, ok := schemacontext.BlockNestingLevel(ctx); !ok || lvl != 0 {
		return ctx, diags
	}

	for _, inner := range block.Body.Blocks {
		if inner.Type != "provider_meta" {
			continue
		}

		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "provider_meta block is ignored by OpenTofu",
			Subject:  inner.Range().Ptr(),
			Extra:    ilsp.DiagnosticMetadata{Hint: true, Unnecessary: true},
		})
	}

	return ctx, diags
}
