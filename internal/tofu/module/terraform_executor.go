// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"fmt"

	"github.com/opentofu/tofu-ls/internal/tofu/exec"
)

func TofuExecutorForModule(ctx context.Context, modPath string) (exec.TofuExecutor, error) {
	newExecutor, ok := exec.ExecutorFactoryFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no tofu executor provided")
	}

	execPath, err := TofuExecPath(ctx)
	if err != nil {
		return nil, err
	}

	tfExec, err := newExecutor(modPath, execPath)
	if err != nil {
		return nil, err
	}

	opts, ok := exec.ExecutorOptsFromContext(ctx)
	if ok && opts.ExecLogPath != "" {
		tfExec.SetExecLogPath(opts.ExecLogPath)
	}
	if ok && opts.Timeout != 0 {
		tfExec.SetTimeout(opts.Timeout)
	}

	return tfExec, nil
}

func TofuExecPath(ctx context.Context) (string, error) {
	opts, ok := exec.ExecutorOptsFromContext(ctx)
	if ok && opts.ExecPath != "" {
		return opts.ExecPath, nil
	} else {
		return "", NoTofuExecPathErr{}
	}
}
