// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jobs

import (
	"context"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/features/rootmodules/state"
	"github.com/opentofu/tofu-ls/internal/job"
	"github.com/opentofu/tofu-ls/internal/tofu/module"
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
)

// GetTofuVersion obtains "installed" Tofu version
// which can inform what version of core schema to pick.
// Knowing the version is not required though as we can rely on
// the constraint in `required_version` (as parsed via
// [LoadModuleMetadata] and compare it against known released versions.
func GetTofuVersion(ctx context.Context, rootStore *state.RootStore, modPath string) error {
	mod, err := rootStore.RootRecordByPath(modPath)
	if err != nil {
		return err
	}

	// Avoid getting version if getting is already in progress or already known
	if mod.TofuVersionState != op.OpStateUnknown && !job.IgnoreState(ctx) {
		return job.StateNotChangedErr{Dir: document.DirHandleFromPath(modPath)}
	}

	err = rootStore.SetTofuVersionState(modPath, op.OpStateLoading)
	if err != nil {
		return err
	}
	defer rootStore.SetTofuVersionState(modPath, op.OpStateLoaded)

	tfExec, err := module.TofuExecutorForModule(ctx, mod.Path())
	if err != nil {
		sErr := rootStore.UpdateTofuAndProviderVersions(modPath, nil, nil, err)
		if sErr != nil {
			return sErr
		}
		return err
	}

	v, pv, err := tfExec.Version(ctx)

	// TODO: Remove and rely purely on ParseProviderVersions
	// In most cases we get the provider version from the datadir/lockfile
	// but there is an edge case with custom plugin location
	// when this may not be available, so leveraging versions
	// from "terraform version" accounts for this.
	// See https://github.com/hashicorp/terraform-ls/issues/24
	pVersions := providerVersionsFromTfVersion(pv)

	sErr := rootStore.UpdateTofuAndProviderVersions(modPath, v, pVersions, err)
	if sErr != nil {
		return sErr
	}

	return err
}

func providerVersionsFromTfVersion(pv map[string]*version.Version) map[tfaddr.Provider]*version.Version {
	m := make(map[tfaddr.Provider]*version.Version, 0)

	for rawAddr, v := range pv {
		pAddr, err := tfaddr.ParseProviderSource(rawAddr)
		if err != nil {
			// skip unparsable address
			continue
		}
		if pAddr.IsLegacy() {
			// TODO: check for migrations via Registry API?
		}
		m[pAddr] = v
	}

	return m
}
