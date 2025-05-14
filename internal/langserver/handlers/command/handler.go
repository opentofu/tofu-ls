// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"log"

	fmodules "github.com/opentofu/tofu-ls/internal/features/modules"
	frootmodules "github.com/opentofu/tofu-ls/internal/features/rootmodules"
	"github.com/opentofu/tofu-ls/internal/state"
)

type CmdHandler struct {
	StateStore *state.StateStore
	Logger     *log.Logger
	// TODO? Can features contribute commands, so we don't have to import
	// the features here?
	ModulesFeature     *fmodules.ModulesFeature
	RootModulesFeature *frootmodules.RootModulesFeature
}
