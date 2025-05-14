// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package decoder

import (
	"context"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	tfmod "github.com/opentofu/opentofu-schema/module"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/features/variables/state"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
)

type StateReader interface {
	List() ([]*state.VariableRecord, error)
	VariableRecordByPath(path string) (*state.VariableRecord, error)
}

type ModuleReader interface {
	ModuleInputs(modPath string) (map[string]tfmod.Variable, error)
	MetadataReady(dir document.DirHandle) (<-chan struct{}, bool, error)
}

type PathReader struct {
	StateReader  StateReader
	ModuleReader ModuleReader
	UseAnySchema bool
}

var _ decoder.PathReader = &PathReader{}

func (pr *PathReader) Paths(ctx context.Context) []lang.Path {
	paths := make([]lang.Path, 0)

	variableRecords, err := pr.StateReader.List()
	if err != nil {
		return paths
	}

	for _, record := range variableRecords {
		paths = append(paths, lang.Path{
			Path:       record.Path(),
			LanguageID: ilsp.Tfvars.String(),
		})
	}

	return paths
}

// PathContext returns a PathContext for the given path based on the language ID.
func (pr *PathReader) PathContext(path lang.Path) (*decoder.PathContext, error) {
	mod, err := pr.StateReader.VariableRecordByPath(path.Path)
	if err != nil {
		return nil, err
	}
	return variablePathContext(mod, pr.ModuleReader, pr.UseAnySchema)
}
