// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jobs

import (
	"context"

	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/features/rootmodules/state"
	"github.com/opentofu/tofu-ls/internal/job"
	"github.com/opentofu/tofu-ls/internal/tofu/datadir"
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
)

// ParseProviderVersions is a job complimentary to [ObtainSchema]
// in that it obtains versions of providers/schemas from Terraform
// CLI's lock file.
func ParseProviderVersions(ctx context.Context, fs ReadOnlyFS, rootStore *state.RootStore, modPath string) error {
	mod, err := rootStore.RootRecordByPath(modPath)
	if err != nil {
		return err
	}

	// Avoid parsing if it is already in progress or already known
	if mod.InstalledProvidersState != op.OpStateUnknown && !job.IgnoreState(ctx) {
		return job.StateNotChangedErr{Dir: document.DirHandleFromPath(modPath)}
	}

	err = rootStore.SetInstalledProvidersState(modPath, op.OpStateLoading)
	if err != nil {
		return err
	}

	pvm, err := datadir.ParsePluginVersions(fs, modPath)

	sErr := rootStore.UpdateInstalledProviders(modPath, pvm, err)
	if sErr != nil {
		return sErr
	}

	return err
}
