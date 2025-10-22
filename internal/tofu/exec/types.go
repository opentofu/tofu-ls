// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/opentofu/tofu-exec/tfexec"
)

// ExecutorFactory can be used in external consumers of exec pkg
// to enable easy swapping with MockExecutor
type ExecutorFactory func(workDir, execPath string) (TofuExecutor, error)

type Formatter func(ctx context.Context, input []byte) ([]byte, error)

//go:generate go tool github.com/vektra/mockery/v2 --name TofuExecutor --structname Executor --filename executor.go --outpkg mock --output ./mock

type TofuExecutor interface {
	SetLogger(logger *log.Logger)
	SetExecLogPath(path string) error
	SetTimeout(duration time.Duration)
	GetExecPath() string
	Init(ctx context.Context, opts ...tfexec.InitOption) error
	Get(ctx context.Context, opts ...tfexec.GetCmdOption) error
	Format(ctx context.Context, input []byte) ([]byte, error)
	Version(ctx context.Context) (*version.Version, map[string]*version.Version, error)
	Validate(ctx context.Context) ([]tfjson.Diagnostic, error)
	ProviderSchemas(ctx context.Context) (*tfjson.ProviderSchemas, error)
}
