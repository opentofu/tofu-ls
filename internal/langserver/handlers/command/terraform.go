// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2"
	"github.com/opentofu/tofu-ls/internal/langserver/cmd"
	"github.com/opentofu/tofu-ls/internal/langserver/progress"
	"github.com/opentofu/tofu-ls/internal/uri"
)

const tofuVersionRequestVersion = 0

type tofuInfoResponse struct {
	FormatVersion     int    `json:"v"`
	RequiredVersion   string `json:"required_version,omitempty"`
	DiscoveredVersion string `json:"discovered_version,omitempty"`
}

func (h *CmdHandler) TofuVersionRequestHandler(ctx context.Context, args cmd.CommandArgs) (interface{}, error) {
	progress.Begin(ctx, "Initializing")
	defer func() {
		progress.End(ctx, "Finished")
	}()

	response := tofuInfoResponse{
		FormatVersion: tofuVersionRequestVersion,
	}

	progress.Report(ctx, "Finding current module info ...")
	modUri, ok := args.GetString("uri")
	if !ok || modUri == "" {
		return response, fmt.Errorf("%w: expected module uri argument to be set", jrpc2.InvalidParams.Err())
	}

	if !uri.IsURIValid(modUri) {
		return response, fmt.Errorf("URI %q is not valid", modUri)
	}

	modPath, err := uri.PathFromURI(modUri)
	if err != nil {
		return response, err
	}

	progress.Report(ctx, "Recording terraform version info ...")

	tofuVersion := h.RootModulesFeature.TofuVersion(modPath)
	if tofuVersion != nil {
		response.DiscoveredVersion = tofuVersion.String()
	}

	coreRequirements, err := h.ModulesFeature.CoreRequirements(modPath)
	if err != nil {
		return response, err
	}
	if coreRequirements != nil {
		response.RequiredVersion = coreRequirements.String()
	}

	progress.Report(ctx, "Sending response ...")

	return response, nil
}
