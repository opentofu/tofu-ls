// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jobs

import (
	"context"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	idecoder "github.com/opentofu/tofu-ls/internal/decoder"
	"github.com/opentofu/tofu-ls/internal/document"
	fdecoder "github.com/opentofu/tofu-ls/internal/features/variables/decoder"
	"github.com/opentofu/tofu-ls/internal/features/variables/state"
	"github.com/opentofu/tofu-ls/internal/job"
	ilsp "github.com/opentofu/tofu-ls/internal/lsp"
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
)

// DecodeVarsReferences collects reference origins within
// variable files (*.tfvars) where each valid attribute
// (as informed by schema provided via [LoadModuleMetadata])
// is considered an origin.
//
// This is useful in hovering over those variable names,
// go-to-definition and go-to-references.
func DecodeVarsReferences(ctx context.Context, varStore *state.VariableStore, moduleFeature fdecoder.ModuleReader, modPath string) error {
	mod, err := varStore.VariableRecordByPath(modPath)
	if err != nil {
		return err
	}

	// TODO: Avoid collection if upstream (parsing) job reported no changes

	// Avoid collection if it is already in progress or already done
	if mod.VarsRefOriginsState != op.OpStateUnknown && !job.IgnoreState(ctx) {
		return job.StateNotChangedErr{Dir: document.DirHandleFromPath(modPath)}
	}

	err = varStore.SetVarsReferenceOriginsState(modPath, op.OpStateLoading)
	if err != nil {
		return err
	}

	d := decoder.NewDecoder(&fdecoder.PathReader{
		StateReader:  varStore,
		ModuleReader: moduleFeature,
		UseAnySchema: true,
	})
	d.SetContext(idecoder.DecoderContext(ctx))

	varsDecoder, err := d.Path(lang.Path{
		Path:       modPath,
		LanguageID: ilsp.Tfvars.String(),
	})
	if err != nil {
		return err
	}

	origins, rErr := varsDecoder.CollectReferenceOrigins()
	sErr := varStore.UpdateVarsReferenceOrigins(modPath, origins, rErr)
	if sErr != nil {
		return sErr
	}

	return rErr
}
