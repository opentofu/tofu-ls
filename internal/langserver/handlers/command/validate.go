// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/job"
	"github.com/opentofu/tofu-ls/internal/langserver/cmd"
	"github.com/opentofu/tofu-ls/internal/langserver/progress"
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
	"github.com/opentofu/tofu-ls/internal/uri"
)

func (h *CmdHandler) TofuValidateHandler(ctx context.Context, args cmd.CommandArgs) (interface{}, error) {
	dirUri, ok := args.GetString("uri")
	if !ok || dirUri == "" {
		return nil, fmt.Errorf("%w: expected module uri argument to be set", jrpc2.InvalidParams.Err())
	}

	if !uri.IsURIValid(dirUri) {
		return nil, fmt.Errorf("URI %q is not valid", dirUri)
	}

	dirHandle := document.DirHandleFromURI(dirUri)

	progress.Begin(ctx, "Validating")
	defer func() {
		progress.End(ctx, "Finished")
	}()

	progress.Report(ctx, "Running tofu validate ...")
	id, err := h.StateStore.JobStore.EnqueueJob(ctx, job.Job{
		Dir: dirHandle,
		Func: func(ctx context.Context) error {
			return nil //module.TofuValidate(ctx, h.StateStore.Modules, dirHandle.Path())
		},
		Type:        op.OpTypeTofuValidate.String(),
		IgnoreState: true,
	})
	if err != nil {
		return nil, err
	}

	return nil, h.StateStore.JobStore.WaitForJobs(ctx, id)
}
