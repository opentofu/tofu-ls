// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validations

import (
	"testing"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
)

const (
	testFilename     = "test.tf"
	testJSONFilename = "test.tf.json"
)

func makeTarget(addrType, name string) reference.Target {
	return makeTargetInFile(addrType, name, testFilename)
}

func makeTargetInFile(addrType, name, filename string) reference.Target {
	return reference.Target{
		Addr: lang.Address{
			lang.RootStep{Name: addrType},
			lang.AttrStep{Name: name},
		},
		RangePtr: &hcl.Range{Filename: filename},
	}
}

func makeOrigin(addrType, name string) reference.LocalOrigin {
	return reference.LocalOrigin{
		Addr: lang.Address{
			lang.RootStep{Name: addrType},
			lang.AttrStep{Name: name},
		},
	}
}

func parseHCLTestFile(t *testing.T, filename, src string) *hcl.File {
	t.Helper()

	file, diags := hclsyntax.ParseConfig([]byte(src), filename, hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("failed to parse test configuration: %s", diags)
	}

	return file
}

func parseTestFiles(t *testing.T, src string) map[string]*hcl.File {
	return map[string]*hcl.File{testFilename: parseHCLTestFile(t, testFilename, src)}
}

func parseJSONTestFiles(t *testing.T, src string) map[string]*hcl.File {
	t.Helper()

	file, diags := hcljson.Parse([]byte(src), testJSONFilename)
	if diags.HasErrors() {
		t.Fatalf("failed to parse test JSON configuration: %s", diags)
	}

	return map[string]*hcl.File{testJSONFilename: file}
}

func TestDeprecatedVariableTargets(t *testing.T) {
	targets := deprecatedVariableTargets(parseTestFiles(t, `variable "foo" {
		deprecated = "Use another input instead."
	}`))

	if !targets[deprecatedVariableKey(testFilename, "foo")] {
		t.Fatalf("deprecated variable target was not detected: %#v", targets)
	}
}

func TestUnusedTargets_DeprecatedVariableIsScopedToFile(t *testing.T) {
	const deprecatedFilename = "deprecated.tf"
	const activeFilename = "active.tf"

	pathCtx := &decoder.PathContext{
		ReferenceTargets: reference.Targets{
			makeTargetInFile("var", "deprecated", deprecatedFilename),
			makeTargetInFile("var", "active", activeFilename),
		},
		Files: map[string]*hcl.File{
			deprecatedFilename: parseHCLTestFile(t, deprecatedFilename, `variable "deprecated" {
				deprecated = "Use another input instead."
			}`),
			activeFilename: parseHCLTestFile(t, activeFilename, `variable "active" {
			}`),
		},
	}

	diags := UnusedTargets(t.Context(), pathCtx)
	if got := len(diags[deprecatedFilename]); got != 0 {
		t.Fatalf("deprecated file had %d diagnostics, want 0", got)
	}
	if got := len(diags[activeFilename]); got != 1 {
		t.Fatalf("active file had %d diagnostics, want 1", got)
	}
}

func TestUnusedTargets(t *testing.T) {
	tests := []struct {
		name      string
		targets   reference.Targets
		origins   reference.Origins
		files     map[string]*hcl.File
		wantCount int
	}{
		{
			name:      "unused variable",
			targets:   reference.Targets{makeTarget("var", "foo")},
			origins:   reference.Origins{},
			wantCount: 1,
		},
		{
			name:    "unused non-deprecated variable",
			targets: reference.Targets{makeTarget("var", "foo")},
			files: parseTestFiles(t, `variable "foo" {
			}`),
			origins:   reference.Origins{},
			wantCount: 1,
		},
		{
			name:    "unused local matching deprecated variable name",
			targets: reference.Targets{makeTarget("local", "foo")},
			files: parseTestFiles(t, `variable "foo" {
				deprecated = "Use another input instead."
			}`),
			origins:   reference.Origins{},
			wantCount: 1,
		},
		{
			name:    "unused deprecated variable",
			targets: reference.Targets{makeTarget("var", "foo")},
			files: parseTestFiles(t, `variable "foo" {
				deprecated = "Use another input instead."
			}`),
			origins:   reference.Origins{},
			wantCount: 0,
		},
		{
			name:    "unused deprecated JSON variable",
			targets: reference.Targets{makeTargetInFile("var", "foo", testJSONFilename)},
			files: parseJSONTestFiles(t, `{
				"variable": {
					"foo": {
						"deprecated": "Use another input instead."
					}
				}
			}`),
			origins:   reference.Origins{},
			wantCount: 0,
		},
		{
			name:      "unused local",
			targets:   reference.Targets{makeTarget("local", "bar")},
			origins:   reference.Origins{},
			wantCount: 1,
		},
		{
			name:      "used variable",
			targets:   reference.Targets{makeTarget("var", "foo")},
			origins:   reference.Origins{makeOrigin("var", "foo")},
			wantCount: 0,
		},
		{
			name:      "used local",
			targets:   reference.Targets{makeTarget("local", "bar")},
			origins:   reference.Origins{makeOrigin("local", "bar")},
			wantCount: 0,
		},
		{
			name:      "multiple unused",
			targets:   reference.Targets{makeTarget("var", "a"), makeTarget("local", "b")},
			origins:   reference.Origins{},
			wantCount: 2,
		},
		{
			name:      "one used one unused",
			targets:   reference.Targets{makeTarget("var", "used"), makeTarget("var", "unused")},
			origins:   reference.Origins{makeOrigin("var", "used")},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathCtx := &decoder.PathContext{
				ReferenceTargets: tt.targets,
				ReferenceOrigins: tt.origins,
				Files:            tt.files,
			}

			diags := UnusedTargets(t.Context(), pathCtx)
			got := len(diags[testFilename])
			if got != tt.wantCount {
				t.Errorf("got %d diagnostics, want %d", got, tt.wantCount)
			}
		})
	}
}
