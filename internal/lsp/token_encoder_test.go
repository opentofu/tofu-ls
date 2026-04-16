// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/tofu-ls/internal/lsp/semtok"
	"github.com/opentofu/tofu-ls/internal/protocol"
	"github.com/opentofu/tofu-ls/internal/source"
)

func TestTokenEncoder_singleLineTokens(t *testing.T) {
	bytes := []byte(`myblock "mytype" {
  str_attr = "something"
  num_attr = 42
  bool_attr = true
}`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type: lang.TokenBlockType,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 7},
				},
			},
			{
				Type: lang.TokenBlockLabel,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 9, Byte: 8},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 16},
				},
			},
			{
				Type: lang.TokenAttrName,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 21},
					End:      hcl.Pos{Line: 2, Column: 11, Byte: 29},
				},
			},
			{
				Type: lang.TokenAttrName,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 3, Column: 3, Byte: 46},
					End:      hcl.Pos{Line: 3, Column: 11, Byte: 54},
				},
			},
			{
				Type: lang.TokenAttrName,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 4, Column: 3, Byte: 62},
					End:      hcl.Pos{Line: 4, Column: 12, Byte: 71},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		0, 0, 7, 10, 0,
		0, 8, 8, 11, 0,
		1, 2, 8, 9, 0,
		1, 2, 8, 9, 0,
		1, 2, 9, 9, 0,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_unknownTokenType(t *testing.T) {
	bytes := []byte(`variable "test" {
  type = string
  default = "foo"
}
`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type:      lang.SemanticTokenType("unknown"),
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 9, Byte: 8},
				},
			},
			{
				Type:      lang.SemanticTokenType("another-unknown"),
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 20},
					End:      hcl.Pos{Line: 2, Column: 7, Byte: 24},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 3, Column: 3, Byte: 36},
					End:      hcl.Pos{Line: 3, Column: 10, Byte: 43},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		2, 2, 7, 9, 0,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_multiLineTokens(t *testing.T) {
	bytes := []byte(`myblock "mytype" {
  str_attr = "something"
  num_attr = 42
  bool_attr = true
}`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type: lang.TokenAttrName,
				Range: hcl.Range{
					Filename: "test.tf",
					// Attribute name would actually never span
					// multiple lines, but we don't have any token
					// type that would *yet*
					Start: hcl.Pos{Line: 2, Column: 3, Byte: 21},
					End:   hcl.Pos{Line: 4, Column: 12, Byte: 71},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		1, 2, 24, 9, 0,
		1, 0, 15, 9, 0,
		1, 0, 11, 9, 0,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_deltaStartCharBug(t *testing.T) {
	bytes := []byte(`resource "aws_iam_role_policy" "firehose_s3_access" {
}
`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type:      lang.TokenBlockType,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 9, Byte: 8},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: []lang.SemanticTokenModifier{lang.TokenModifierDependent},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 10, Byte: 9},
					End:      hcl.Pos{Line: 1, Column: 31, Byte: 30},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 32, Byte: 31},
					End:      hcl.Pos{Line: 1, Column: 52, Byte: 51},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		0, 0, 8, 10, 0,
		0, 9, 21, 11, 2,
		0, 22, 20, 11, 0,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_tokenModifiers(t *testing.T) {
	bytes := []byte(`myblock "mytype" {
  str_attr = "something"
  num_attr = 42
  bool_attr = true
}`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type: lang.TokenBlockType,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 7},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 9, Byte: 8},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 16},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 21},
					End:      hcl.Pos{Line: 2, Column: 11, Byte: 29},
				},
			},
			{
				Type: lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{
					lang.TokenModifierDependent,
				},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 3, Column: 3, Byte: 46},
					End:      hcl.Pos{Line: 3, Column: 11, Byte: 54},
				},
			},
			{
				Type: lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{
					lang.TokenModifierDependent,
				},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 4, Column: 3, Byte: 62},
					End:      hcl.Pos{Line: 4, Column: 12, Byte: 71},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		0, 0, 7, 10, 0,
		0, 8, 8, 11, 0,
		1, 2, 8, 9, 0,
		1, 2, 8, 9, 2,
		1, 2, 9, 9, 2,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_unsupported(t *testing.T) {
	bytes := []byte(`myblock "mytype" {
  str_attr = "something"
  num_attr = 42
  bool_attr = true
}`)
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type: lang.TokenBlockType,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 7},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 9, Byte: 8},
					End:      hcl.Pos{Line: 1, Column: 8, Byte: 16},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 21},
					End:      hcl.Pos{Line: 2, Column: 11, Byte: 29},
				},
			},
			{
				Type: lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{
					lang.TokenModifierDependent,
				},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 3, Column: 3, Byte: 46},
					End:      hcl.Pos{Line: 3, Column: 11, Byte: 54},
				},
			},
			{
				Type: lang.TokenAttrName,
				Modifiers: []lang.SemanticTokenModifier{
					lang.TokenModifierDependent,
				},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 4, Column: 3, Byte: 62},
					End:      hcl.Pos{Line: 4, Column: 12, Byte: 71},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     []string{"hcl-blockType", "hcl-attrName"},
			TokenModifiers: []string{},
		},
	}
	data := te.Encode()
	expectedData := []uint32{
		0, 0, 7, 1, 0,
		1, 2, 8, 0, 0,
		1, 2, 8, 0, 0,
		1, 2, 9, 0, 0,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v",
			expectedData, data)
	}
}

func TestTokenEncoder_multiLineTokensWithChildTokens(t *testing.T) {
	// Reproduces https://github.com/opentofu/tofu-ls/issues/156
	bytes := []byte(`variable "a" {
  default = "hello"
}

variable "b" {
  default = "world"
}

locals {
  test = <<-EOT
    a: ${var.a}
    b: ${var.b}
  EOT
}`)
	// The TokenEncoder configured here simulates the result that the handlers/semantic_token.go would generate.
	te := &TokenEncoder{
		Lines: source.MakeSourceLines("test.tf", bytes),
		Tokens: []lang.SemanticToken{
			{
				Type:      lang.TokenBlockType,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:      hcl.Pos{Line: 1, Column: 9, Byte: 8},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 1, Column: 10, Byte: 9},
					End:      hcl.Pos{Line: 1, Column: 13, Byte: 12},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 2, Column: 3, Byte: 17},
					End:      hcl.Pos{Line: 2, Column: 10, Byte: 24},
				},
			},
			{
				Type: lang.TokenString,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 2, Column: 13, Byte: 27},
					End:      hcl.Pos{Line: 2, Column: 20, Byte: 34},
				},
			},
			{
				Type:      lang.TokenBlockType,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 5, Column: 1, Byte: 38},
					End:      hcl.Pos{Line: 5, Column: 9, Byte: 46},
				},
			},
			{
				Type:      lang.TokenBlockLabel,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 5, Column: 10, Byte: 47},
					End:      hcl.Pos{Line: 5, Column: 13, Byte: 50},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-variable")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 6, Column: 3, Byte: 55},
					End:      hcl.Pos{Line: 6, Column: 10, Byte: 62},
				},
			},
			{
				Type: lang.TokenString,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 6, Column: 13, Byte: 65},
					End:      hcl.Pos{Line: 6, Column: 20, Byte: 72},
				},
			},
			{
				Type:      lang.TokenBlockType,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-locals")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 9, Column: 1, Byte: 76},
					End:      hcl.Pos{Line: 9, Column: 7, Byte: 82},
				},
			},
			{
				Type:      lang.TokenAttrName,
				Modifiers: lang.SemanticTokenModifiers{lang.SemanticTokenModifier("opentofu-locals")},
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 10, Column: 3, Byte: 87},
					End:      hcl.Pos{Line: 10, Column: 7, Byte: 91},
				},
			},
			{
				Type: lang.TokenString,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 11, Column: 5, Byte: 105},
					End:      hcl.Pos{Line: 11, Column: 8, Byte: 108},
				},
			},
			{
				Type: lang.TokenReferenceStep,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 11, Column: 10, Byte: 110},
					End:      hcl.Pos{Line: 11, Column: 13, Byte: 113},
				},
			},
			{
				Type: lang.TokenReferenceStep,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 11, Column: 14, Byte: 114},
					End:      hcl.Pos{Line: 11, Column: 15, Byte: 115},
				},
			},
			{
				Type: lang.TokenString,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 11, Column: 16, Byte: 116},
					End:      hcl.Pos{Line: 12, Column: 8, Byte: 124},
				},
			},
			{
				Type: lang.TokenReferenceStep,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 12, Column: 10, Byte: 126},
					End:      hcl.Pos{Line: 12, Column: 13, Byte: 129},
				},
			},
			{
				Type: lang.TokenReferenceStep,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 12, Column: 14, Byte: 130},
					End:      hcl.Pos{Line: 12, Column: 15, Byte: 131},
				},
			},
			{
				Type: lang.TokenString,
				Range: hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 12, Column: 16, Byte: 132},
					End:      hcl.Pos{Line: 13, Column: 1, Byte: 133},
				},
			},
		},
		ClientCaps: protocol.SemanticTokensClientCapabilities{
			TokenTypes:     serverTokenTypes.AsStrings(),
			TokenModifiers: serverTokenModifiers.AsStrings(),
		},
	}
	data := te.Encode()

	noModMask := uint32(TokenModifiersLegend(serverTokenModifiers.AsStrings()).BitMask(semtok.TokenModifiers{}))
	varModMask := uint32(TokenModifiersLegend(serverTokenModifiers.AsStrings()).BitMask(semtok.TokenModifiers{"opentofu-variable"}))
	localsModMask := uint32(TokenModifiersLegend(serverTokenModifiers.AsStrings()).BitMask(semtok.TokenModifiers{"opentofu-locals"}))
	blockTypeIdx := uint32(TokenTypesLegend(serverTokenTypes.AsStrings()).Index(semtok.TokenType(lang.TokenBlockType)))
	blockLabelIdx := uint32(TokenTypesLegend(serverTokenTypes.AsStrings()).Index(semtok.TokenType(lang.TokenBlockLabel)))
	attrNameIdx := uint32(TokenTypesLegend(serverTokenTypes.AsStrings()).Index(semtok.TokenType(lang.TokenAttrName)))
	strIdx := uint32(TokenTypesLegend(serverTokenTypes.AsStrings()).Index(semtok.TokenType(lang.TokenString)))
	refIdx := uint32(TokenTypesLegend(serverTokenTypes.AsStrings()).Index(semtok.TokenType(lang.TokenReferenceStep)))

	expectedData := []uint32{
		0, 0, 8, blockTypeIdx, varModMask,
		0, 9, 3, blockLabelIdx, varModMask,
		1, 2, 7, attrNameIdx, varModMask,
		0, 10, 7, strIdx, noModMask,
		3, 0, 8, blockTypeIdx, varModMask,
		0, 9, 3, blockLabelIdx, varModMask,
		1, 2, 7, attrNameIdx, varModMask,
		0, 10, 7, strIdx, noModMask,
		3, 0, 6, blockTypeIdx, localsModMask,
		1, 2, 4, attrNameIdx, localsModMask,
		1, 4, 3, strIdx, noModMask,
		0, 5, 3, refIdx, noModMask,
		0, 4, 1, refIdx, noModMask,
		0, 15, 15, strIdx, noModMask,
		1, 0, 7, strIdx, noModMask,
		0, 4294967290, 3, refIdx, noModMask,
		0, 4, 1, refIdx, noModMask,
		0, 15, 15, strIdx, noModMask,
		1, 0, 0, strIdx, noModMask,
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("unexpected encoded data.\nexpected: %#v\ngiven:    %#v\ndiff: %s",
			expectedData, data, diff)
	}
}
