// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/schemacontext"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
)

func TestUnnecessaryProviderMetaBlock(t *testing.T) {
	tests := []struct {
		name string
		node string
		want hcl.Diagnostics
	}{
		{
			name: "one diag when provider_meta block is inside terraform block",
			node: `terraform {
				provider_meta "x" {
					name = "abc"
				}
			}`,
			want: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "provider_meta block is ignored since OpenTofu 1.14",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 5,
							Byte:   16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 6,
							Byte:   59,
						},
					},
					Extra: ilsp.DiagnosticMetadata{Hint: true, Unnecessary: true},
				},
			},
		},
		{
			name: "multiple diags when multiple provider_meta blocks are present, other blocks unaffected",
			node: `terraform {
				provider_meta "a" {
					name = "abc"
				}

				backend "local" {}

				provider_meta "x" {
					name = "xyz"
				}
			}`,
			want: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "provider_meta block is ignored since OpenTofu 1.14",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 5,
							Byte:   16,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 6,
							Byte:   59,
						},
					},
					Extra: ilsp.DiagnosticMetadata{Hint: true, Unnecessary: true},
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "provider_meta block is ignored since OpenTofu 1.14",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   8,
							Column: 5,
							Byte:   89,
						},
						End: hcl.Pos{
							Line:   10,
							Column: 6,
							Byte:   132,
						},
					},
					Extra: ilsp.DiagnosticMetadata{Hint: true, Unnecessary: true},
				},
			},
		},
		{
			name: "no diag when provider_meta block is outside terraform block",
			node: `resource "a" "b" {
				provider_meta "x" {
					name = "abc"
				}
			}`,
			want: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%2d-%s", i, tt.name), func(t *testing.T) {
			ctx := context.Background()
			ctx = schemacontext.WithBlockNestingLevel(ctx, 0)

			f, parseDiags := hclsyntax.ParseConfig([]byte(tt.node), "test.tf", hcl.InitialPos)
			if parseDiags.HasErrors() {
				t.Fatalf("failed to parse node config: %s", parseDiags)
			}

			node := f.Body.(*hclsyntax.Body).Blocks[0]
			upmb := UnnecessaryProviderMetaBlock{}
			_, diags := upmb.Visit(ctx, node, nil)
			if diff := cmp.Diff(tt.want, diags); diff != "" {
				t.Fatalf("unexpected diagnostics -want +got: %s", diff)
			}
		})
	}
}
