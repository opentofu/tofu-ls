// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	"github.com/opentofu/tofu-ls/internal/filesystem"
	"github.com/opentofu/tofu-ls/internal/langserver"
	"github.com/opentofu/tofu-ls/internal/langserver/cmd"
	"github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
	"github.com/opentofu/tofu-ls/internal/uri"
	"github.com/opentofu/tofu-ls/internal/walker"
	"github.com/stretchr/testify/mock"
)

func TestLangServer_workspaceExecuteCommand_tofuVersion_basic(t *testing.T) {
	modDir := t.TempDir()
	modUri := uri.FromPath(modDir)

	s, err := state.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}

	eventBus := eventbus.NewEventBus()
	mockCalls := &exec.TofuMockCalls{
		PerWorkDir: map[string][]*mock.Call{
			modDir: validTfMockCalls(),
		},
	}
	fs := filesystem.NewFilesystem(s.DocumentStore)
	features, err := NewTestFeatures(eventBus, s, fs, mockCalls)
	if err != nil {
		t.Fatal(err)
	}

	err = features.Modules.Store.Add(modDir)
	if err != nil {
		t.Fatal(err)
	}
	err = features.RootModules.Store.Add(modDir)
	if err != nil {
		t.Fatal(err)
	}

	metadata := &tfmod.Meta{
		Path:             modDir,
		CoreRequirements: testConstraint(t, "~> 0.15"),
	}

	err = features.Modules.Store.UpdateMetadata(modDir, metadata, nil)
	if err != nil {
		t.Fatal(err)
	}

	ver, err := version.NewVersion("1.1.0")
	if err != nil {
		t.Fatal(err)
	}

	err = features.RootModules.Store.UpdateTerraformAndProviderVersions(modDir, ver, map[tfaddr.Provider]*version.Version{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	wc := walker.NewWalkerCollector()

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls:  mockCalls,
		StateStore:      s,
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
		"processId": 12345
	}`, modUri)})
	waitForWalkerPath(t, s, wc, document.DirHandleFromURI(modUri))
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "workspace/executeCommand",
		ReqParams: fmt.Sprintf(`{
		"command": %q,
		"arguments": ["uri=%s"]
	}`, cmd.Name("module.terraform"), modUri)}, `{
		"jsonrpc": "2.0",
		"id": 2,
		"result": {
			"v": 0,
			"required_version": "~\u003e 0.15",
			"discovered_version": "1.1.0"
		}
	}`)
}
