// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package decoder

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfmodule "github.com/opentofu/opentofu-schema/module"
	tfschema "github.com/opentofu/opentofu-schema/schema"
	"github.com/opentofu/tofu-ls/internal/features/modules/state"
)

func schemaForModule(mod *state.ModuleRecord, stateReader CombinedReader) (*schema.BodySchema, error) {
	resolvedVersion := tfschema.ResolveVersion(stateReader.TofuVersion(mod.Path()), mod.Meta.CoreRequirements)
	sm := tfschema.NewSchemaMerger(mustCoreSchemaForVersion(resolvedVersion))
	sm.SetTofuVersion(resolvedVersion)
	sm.SetStateReader(stateReader)

	meta := &tfmodule.Meta{
		Path:                 mod.Path(),
		CoreRequirements:     mod.Meta.CoreRequirements,
		ProviderRequirements: mod.Meta.ProviderRequirements,
		ProviderReferences:   mod.Meta.ProviderReferences,
		Variables:            mod.Meta.Variables,
		Filenames:            mod.Meta.Filenames,
		ModuleCalls:          mod.Meta.ModuleCalls,
	}

	return sm.SchemaForModule(meta)
}

func mustCoreSchemaForVersion(v *version.Version) *schema.BodySchema {
	s, err := tfschema.CoreModuleSchemaForVersion(v)
	if err != nil {
		// this should never happen
		panic(err)
	}
	return s
}
