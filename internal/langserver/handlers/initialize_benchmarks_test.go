// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opentofu/tofu-exec/tfexec"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/langserver"
	"github.com/opentofu/tofu-ls/internal/langserver/session"
	"github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/discovery"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
	"github.com/opentofu/tofu-ls/internal/walker"
	"github.com/opentofu/tofudl"
)

func BenchmarkInitializeFolder_basic(b *testing.B) {
	modules := []struct {
		name       string
		sourceAddr string
	}{
		{
			name:       "local-single-module-no-provider",
			sourceAddr: mustAbs(b, filepath.Join("testdata", "single-module-no-provider")),
		},
		{
			name:       "local-single-submodule-no-provider",
			sourceAddr: mustAbs(b, filepath.Join("testdata", "single-submodule")),
		},
		{
			name:       "local-single-module-random",
			sourceAddr: mustAbs(b, filepath.Join("testdata", "single-module-random")),
		},
		{
			name:       "local-single-module-aws",
			sourceAddr: mustAbs(b, filepath.Join("testdata", "single-module-aws")),
		},
		// TODO: module version pinning - requires explicit git cloning
		{
			name:       "aws-consul",
			sourceAddr: "github.com/hashicorp/terraform-aws-consul?ref=v0.11.0",
		},
		{
			name:       "aws-eks",
			sourceAddr: "terraform-aws-modules/eks/aws",
		},
		{
			name:       "aws-vpc",
			sourceAddr: "terraform-aws-modules/vpc/aws",
		},
		{
			name:       "google-project",
			sourceAddr: "terraform-google-modules/project-factory/google",
		},
		{
			name:       "google-network",
			sourceAddr: "terraform-google-modules/network/google",
		},
		{
			name:       "google-gke",
			sourceAddr: "terraform-google-modules/kubernetes-engine/google",
		},
		{
			name:       "k8s-metrics-server",
			sourceAddr: "cookielab/metrics-server/kubernetes",
		},
		{
			name:       "k8s-dashboard",
			sourceAddr: "cookielab/dashboard/kubernetes",
		},
	}

	workDir := b.TempDir()

	dl, err := tofudl.New()
	if err != nil {
		log.Fatalf("error when instantiating tofudl %s", err)
	}

	binary, err := dl.Download(b.Context())
	if err != nil {
		log.Fatalf("error when downloading %s", err)
	}

	execPath := filepath.Join(workDir, "tofu")
	// Windows executable case
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, binary, 0755); err != nil {
		log.Fatalf("error when writing the file %s: %s", execPath, err)
	}

	b.Cleanup(func() {
		if err := os.Remove(execPath); err != nil {
			b.Fatal(err)
		}
	})

	for _, mod := range modules {
		b.Run(mod.name, func(b *testing.B) {
			rootDir := b.TempDir()

			tf, err := exec.NewExecutor(rootDir, execPath)
			if err != nil {
				b.Fatal(err)
			}
			err = tf.Init(b.Context(), tfexec.FromModule(mod.sourceAddr))
			if err != nil {
				b.Fatal(err)
			}

			b.Cleanup(func() {
				os.RemoveAll(rootDir)
			})
			b.StopTimer()
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				rootDir := document.DirHandleFromPath(rootDir)
				ss, err := state.NewStateStore()
				if err != nil {
					b.Fatal(err)
				}
				wc := walker.NewWalkerCollector()

				b.StartTimer()
				ls := langserver.NewLangServerMock(b, func(ctx context.Context) session.Session {
					d := &discovery.Discovery{}
					sessCtx, stopSession := context.WithCancel(ctx)
					return &service{
						logger:          discardLogs,
						srvCtx:          ctx,
						sessCtx:         sessCtx,
						stopSession:     stopSession,
						tfDiscoFunc:     d.LookPath,
						tfExecFactory:   exec.NewExecutor,
						walkerCollector: wc,
						stateStore:      ss,
					}
				})
				stop := ls.Start(b)

				ls.Call(b, &langserver.CallRequest{
					Method: "initialize",
					ReqParams: fmt.Sprintf(`{
						"capabilities": {
							"workspace": {
								"workspaceFolders": true
							}
						},
						"rootUri": %q,
						"processId": 12345,
						"workspaceFolders": [
							{
								"uri": %q,
								"name": "root"
							}
						]
					}`, rootDir.URI, rootDir.URI)})
				waitForWalkerPath(b, ss, wc, rootDir)
				b.StopTimer()

				stop()
			}
		})
	}
}

func mustAbs(b *testing.B, path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		b.Fatal(err)
	}
	return absPath
}
