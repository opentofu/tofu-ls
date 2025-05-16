// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/creachadair/jrpc2/handler"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	fmodules "github.com/opentofu/tofu-ls/internal/features/modules"
	frootmodules "github.com/opentofu/tofu-ls/internal/features/rootmodules"
	fvariables "github.com/opentofu/tofu-ls/internal/features/variables"
	"github.com/opentofu/tofu-ls/internal/filesystem"
	"github.com/opentofu/tofu-ls/internal/langserver/session"
	"github.com/opentofu/tofu-ls/internal/registry"
	"github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/terraform/discovery"
	"github.com/opentofu/tofu-ls/internal/terraform/exec"
	"github.com/opentofu/tofu-ls/internal/walker"
)

type MockSessionInput struct {
	TerraformCalls     *exec.TerraformMockCalls
	AdditionalHandlers map[string]handler.Func
	StateStore         *state.StateStore
	WalkerCollector    *walker.WalkerCollector
	RegistryServer     *httptest.Server
	Features           *Features
	FileSystem         *filesystem.Filesystem
	EventBus           *eventbus.EventBus
}

type mockSession struct {
	mockInput      *MockSessionInput
	registryServer *httptest.Server

	stopFunc     func()
	stopCalled   bool
	stopCalledMu *sync.RWMutex
}

func (ms *mockSession) new(srvCtx context.Context) session.Session {
	sessCtx, stopSession := context.WithCancel(srvCtx)
	ms.stopFunc = stopSession

	var handlers map[string]handler.Func
	var stateStore *state.StateStore
	var features *Features
	var walkerCollector *walker.WalkerCollector
	var fileSystem *filesystem.Filesystem
	var eventBus *eventbus.EventBus
	if ms.mockInput != nil {
		stateStore = ms.mockInput.StateStore
		walkerCollector = ms.mockInput.WalkerCollector
		handlers = ms.mockInput.AdditionalHandlers
		ms.registryServer = ms.mockInput.RegistryServer
		features = ms.mockInput.Features
		fileSystem = ms.mockInput.FileSystem
		eventBus = ms.mockInput.EventBus
	}

	var tfCalls *exec.TerraformMockCalls
	if ms.mockInput != nil && ms.mockInput.TerraformCalls != nil {
		tfCalls = ms.mockInput.TerraformCalls
	}

	d := &discovery.MockDiscovery{
		Path: "tf-mock",
	}

	regClient := registry.NewClient()
	if ms.registryServer == nil {
		ms.registryServer = defaultRegistryServer()
	}
	ms.registryServer.Start()

	regClient.BaseAPIURL = ms.registryServer.URL

	svc := &service{
		logger:             testLogger(),
		srvCtx:             srvCtx,
		sessCtx:            sessCtx,
		stopSession:        ms.stop,
		tfDiscoFunc:        d.LookPath,
		tfExecFactory:      exec.NewMockExecutor(tfCalls),
		additionalHandlers: handlers,
		stateStore:         stateStore,
		walkerCollector:    walkerCollector,
		registryClient:     regClient,
		features:           features,
		fs:                 fileSystem,
		eventBus:           eventBus,
	}

	return svc
}

func defaultRegistryServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unexpected Registry API request", 500)
	}))
}

func testLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	return log.New(ioutil.Discard, "", 0)
}

func (ms *mockSession) stop() {
	ms.registryServer.Close()

	ms.stopCalledMu.Lock()
	defer ms.stopCalledMu.Unlock()

	ms.stopFunc()
	ms.stopCalled = true
}

func (ms *mockSession) StopFuncCalled() bool {
	ms.stopCalledMu.RLock()
	defer ms.stopCalledMu.RUnlock()

	return ms.stopCalled
}

func newMockSession(input *MockSessionInput) *mockSession {
	return &mockSession{
		mockInput:    input,
		stopCalledMu: &sync.RWMutex{},
	}
}

func NewMockSession(input *MockSessionInput) session.SessionFactory {
	return newMockSession(input).new
}

func NewTestFeatures(eventBus *eventbus.EventBus, s *state.StateStore, fs *filesystem.Filesystem, tfCalls *exec.TerraformMockCalls) (*Features, error) {
	rootModulesFeature, err := frootmodules.NewRootModulesFeature(eventBus, s, fs, exec.NewMockExecutor(tfCalls))
	if err != nil {
		return nil, err
	}

	modulesFeature, err := fmodules.NewModulesFeature(eventBus, s, fs, rootModulesFeature, registry.Client{})
	if err != nil {
		return nil, err
	}

	variablesFeature, err := fvariables.NewVariablesFeature(eventBus, s, fs, modulesFeature)
	if err != nil {
		return nil, err
	}

	return &Features{
		Modules:     modulesFeature,
		RootModules: rootModulesFeature,
		Variables:   variablesFeature,
	}, nil
}
