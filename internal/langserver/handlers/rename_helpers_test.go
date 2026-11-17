// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
)

func TestFindDefinition(t *testing.T) {
	src := []byte("resource \"aws_s3_bucket\" \"old\" {\n}\n\ndata \"aws_iam_policy\" \"p\" {\n}\n\nlocals {\n  x = 1\n}\n\nmodule \"m\" {\n}\n")
	for _, tc := range []struct {
		pos     hcl.Pos
		name    string
		nameIdx int
		text    string
	}{
		{hcl.Pos{Line: 1, Column: 30, Byte: 29}, "old", 1, "old"},
		{hcl.Pos{Line: 4, Column: 25, Byte: 62}, "p", 2, "p"},
		{hcl.Pos{Line: 8, Column: 4, Byte: 0}, "x", 1, "x"},
		{hcl.Pos{Line: 11, Column: 10, Byte: 0}, "m", 1, "m"},
	} {
		// derive Byte from line/column so the table stays readable
		tc.pos.Byte = bytePos(src, tc.pos.Line, tc.pos.Column)
		sym, ok := findDefinition(src, "main.tf", tc.pos)
		if !ok {
			t.Fatalf("%v: no definition found", tc.pos)
		}
		if sym.name != tc.name || sym.nameIdx != tc.nameIdx {
			t.Errorf("%v: got %q/%d, want %q/%d", tc.pos, sym.name, sym.nameIdx, tc.name, tc.nameIdx)
		}
		if got := string(src[sym.nameRange.Start.Byte:sym.nameRange.End.Byte]); got != tc.text {
			t.Errorf("%v: nameRange covers %q, want %q", tc.pos, got, tc.text)
		}
	}

	if _, ok := findDefinition([]byte("output \"o\" {\n}\n"), "main.tf", hcl.Pos{Line: 1, Column: 3, Byte: 2}); ok {
		t.Error("output must not be renamable")
	}
}

func TestNameInTraversal(t *testing.T) {
	for _, tc := range []struct {
		expr string
		sym  renameSymbol
		want string
	}{
		{"aws_s3_bucket.old.id", renameSymbol{name: "old", nameIdx: 1}, "old"},
		{"aws_s3_bucket.old[0].id", renameSymbol{name: "old", nameIdx: 1}, "old"},
		{"data.aws_iam_policy.p.json", renameSymbol{name: "p", nameIdx: 2}, "p"},
		{"var.old", renameSymbol{name: "old", nameIdx: 1}, "old"},
	} {
		src := []byte("x = " + tc.expr)
		origin := hcl.Range{
			Filename: "main.tf",
			Start:    hcl.Pos{Line: 1, Column: 5, Byte: 4},
			End:      hcl.Pos{Line: 1, Column: 5 + len(tc.expr), Byte: 4 + len(tc.expr)},
		}
		rng, ok := nameInTraversal(src, origin, &tc.sym)
		if !ok {
			t.Fatalf("%s: not found", tc.expr)
		}
		if got := string(src[rng.Start.Byte:rng.End.Byte]); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.expr, got, tc.want)
		}
	}

	// wrong name at the expected step: refuse rather than edit something else
	if _, ok := nameInTraversal([]byte("var.other"), hcl.Range{Start: hcl.InitialPos, End: hcl.Pos{Line: 1, Column: 10, Byte: 9}},
		&renameSymbol{name: "old", nameIdx: 1}); ok {
		t.Error("mismatching name must not match")
	}
}

func bytePos(src []byte, line, col int) int {
	l, c := 1, 1
	for i := range src {
		if l == line && c == col {
			return i
		}
		if src[i] == '\n' {
			l++
			c = 1
		} else {
			c++
		}
	}
	return len(src)
}
