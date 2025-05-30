// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"context"
	"fmt"

	exec_mock "github.com/opentofu/tofu-ls/internal/tofu/exec/mock"
	"github.com/stretchr/testify/mock"
)

type TofuMockCalls struct {
	PerWorkDir map[string][]*mock.Call
	AnyWorkDir []*mock.Call
}

func NewMockExecutor(calls *TofuMockCalls) ExecutorFactory {
	return func(workDir string, execPath string) (TofuExecutor, error) {
		if calls == nil {
			return nil, fmt.Errorf("%s: no mock calls defined", workDir)
		}
		mockCalls := calls.AnyWorkDir
		if len(calls.PerWorkDir) > 0 {
			mc, ok := calls.PerWorkDir[workDir]
			if ok {
				mockCalls = mc
			}
		}
		if len(mockCalls) == 0 {
			return nil, fmt.Errorf("%s: no mock calls available for this workdir", workDir)
		}

		me := &exec_mock.Executor{}
		firstCalls := []*mock.Call{
			{
				Method:        "SetLogger",
				Arguments:     []interface{}{mock.Anything},
				Repeatability: 1,
			},
		}

		me.ExpectedCalls = append(firstCalls, mockCalls...)
		return me, nil
	}
}

var ctxExecutorFactory = ctxKey("executor factory")

func ExecutorFactoryFromContext(ctx context.Context) (ExecutorFactory, bool) {
	f, ok := ctx.Value(ctxExecutorFactory).(ExecutorFactory)
	return f, ok
}

func WithExecutorFactory(ctx context.Context, f ExecutorFactory) context.Context {
	return context.WithValue(ctx, ctxExecutorFactory, f)
}
