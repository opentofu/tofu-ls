// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"fmt"
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	"github.com/opentofu/tofu-ls/internal/filesystem"
	"github.com/opentofu/tofu-ls/internal/langserver"
	"github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/ast"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
	"github.com/opentofu/tofu-ls/internal/tofu/module/operation"
	"github.com/opentofu/tofu-ls/internal/uri"
	"github.com/opentofu/tofu-ls/internal/walker"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLangServer_textDocumentDidSave_runsTofuValidate(t *testing.T) {
	modDir := t.TempDir()
	modURI := uri.FromPath(modDir)
	fileURI := fmt.Sprintf("%s/main.tf", modURI)

	stateStore, err := state.NewStateStore()
	require.NoError(t, err)

	mockCalls := &exec.TofuMockCalls{
		PerWorkDir: map[string][]*mock.Call{
			modDir: append(validTfMockCalls(), &mock.Call{
				Method:        "Validate",
				Repeatability: 1,
				Arguments:     []interface{}{mock.Anything},
				ReturnArguments: []interface{}{
					[]tfjson.Diagnostic{},
					nil,
				},
			}),
		},
	}
	eventBus := eventbus.NewEventBus()
	fs := filesystem.NewFilesystem(stateStore.DocumentStore)
	features, err := NewTestFeatures(eventBus, stateStore, fs, mockCalls)
	require.NoError(t, err)
	require.NoError(t, features.Modules.Store.Add(modDir))

	wc := walker.NewWalkerCollector()
	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TofuCalls:       mockCalls,
		StateStore:      stateStore,
		WalkerCollector: wc,
		Features:        features,
		EventBus:        eventBus,
		FileSystem:      fs,
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
			"capabilities": {},
			"rootUri": %q,
			"processId": 12345,
			"initializationOptions": {
				"experimentalFeatures": {
					"validateOnSave": true
				}
			}
		}`, modURI),
	})
	waitForWalkerPath(t, stateStore, wc, document.DirHandleFromURI(modURI))
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Notify(t, &langserver.CallRequest{
		Method: "textDocument/didSave",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"uri": %q
			}
		}`, fileURI),
	})

	require.Eventually(t, func() bool {
		mod, err := features.Modules.Store.ModuleRecordByPath(modDir)
		if err != nil {
			return false
		}
		return mod.ModuleDiagnosticsState[ast.TofuValidateSource] == operation.OpStateLoaded
	}, time.Second, 10*time.Millisecond, "OpenTofu validation did not finish after save")
}
