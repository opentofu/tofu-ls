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
)

const testFilename = "test.tf"

func makeTarget(addrType, name string) reference.Target {
	return reference.Target{
		Addr: lang.Address{
			lang.RootStep{Name: addrType},
			lang.AttrStep{Name: name},
		},
		RangePtr: &hcl.Range{Filename: testFilename},
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

func TestUnusedTargets(t *testing.T) {
	tests := []struct {
		name      string
		targets   reference.Targets
		origins   reference.Origins
		wantCount int
	}{
		{
			name:      "unused variable",
			targets:   reference.Targets{makeTarget("var", "foo")},
			origins:   reference.Origins{},
			wantCount: 1,
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
			}

			diags := UnusedTargets(t.Context(), pathCtx)
			got := len(diags[testFilename])
			if got != tt.wantCount {
				t.Errorf("got %d diagnostics, want %d", got, tt.wantCount)
			}
		})
	}
}
