// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package modules

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	tfmod "github.com/opentofu/opentofu-schema/module"
	"github.com/opentofu/tofu-ls/internal/document"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	fdecoder "github.com/opentofu/tofu-ls/internal/features/modules/decoder"
	"github.com/opentofu/tofu-ls/internal/features/modules/hooks"
	"github.com/opentofu/tofu-ls/internal/features/modules/jobs"
	"github.com/opentofu/tofu-ls/internal/features/modules/state"
	"github.com/opentofu/tofu-ls/internal/langserver/diagnostics"
	"github.com/opentofu/tofu-ls/internal/registry"
	globalState "github.com/opentofu/tofu-ls/internal/state"
)

// ModulesFeature groups everything related to modules. Its internal
// state keeps track of all modules in the workspace.
type ModulesFeature struct {
	Store    *state.ModuleStore
	eventbus *eventbus.EventBus
	stopFunc context.CancelFunc
	logger   *log.Logger

	rootFeature    fdecoder.RootReader
	stateStore     *globalState.StateStore
	registryClient registry.Client
	fs             jobs.ReadOnlyFS
}

func NewModulesFeature(eventbus *eventbus.EventBus, stateStore *globalState.StateStore, fs jobs.ReadOnlyFS, rootFeature fdecoder.RootReader, registryClient registry.Client) (*ModulesFeature, error) {
	store, err := state.NewModuleStore(stateStore.ProviderSchemas, stateStore.RegistryModules, stateStore.ChangeStore)
	if err != nil {
		return nil, err
	}
	discardLogger := log.New(io.Discard, "", 0)

	return &ModulesFeature{
		Store:          store,
		eventbus:       eventbus,
		stopFunc:       func() {},
		logger:         discardLogger,
		stateStore:     stateStore,
		rootFeature:    rootFeature,
		fs:             fs,
		registryClient: registryClient,
	}, nil
}

func (f *ModulesFeature) SetLogger(logger *log.Logger) {
	f.logger = logger
	f.Store.SetLogger(logger)
}

// Start starts the features separate goroutine.
// It listens to various events from the EventBus and performs corresponding actions.
func (f *ModulesFeature) Start(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	f.stopFunc = cancelFunc

	discoverDone := make(chan struct{}, 10)
	discover := f.eventbus.OnDiscover("feature.modules", discoverDone)

	didOpenDone := make(chan struct{}, 10)
	didOpen := f.eventbus.OnDidOpen("feature.modules", didOpenDone)

	didChangeDone := make(chan struct{}, 10)
	didChange := f.eventbus.OnDidChange("feature.modules", didChangeDone)

	didChangeWatchedDone := make(chan struct{}, 10)
	didChangeWatched := f.eventbus.OnDidChangeWatched("feature.modules", didChangeWatchedDone)

	go func() {
		for {
			select {
			case discover := <-discover:
				// TODO? collect errors
				f.discover(discover.Path, discover.Files)
				discoverDone <- struct{}{}
			case didOpen := <-didOpen:
				// TODO? collect errors
				f.didOpen(didOpen.Context, didOpen.Dir, didOpen.LanguageID)
				didOpenDone <- struct{}{}
			case didChange := <-didChange:
				// TODO? collect errors
				f.didChange(didChange.Context, didChange.Dir)
				didChangeDone <- struct{}{}
			case didChangeWatched := <-didChangeWatched:
				// TODO? collect errors
				f.didChangeWatched(didChangeWatched.Context, didChangeWatched.RawPath, didChangeWatched.ChangeType, didChangeWatched.IsDir)
				didChangeWatchedDone <- struct{}{}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (f *ModulesFeature) Stop() {
	f.stopFunc()
	f.logger.Print("stopped modules feature")
}

func (f *ModulesFeature) PathContext(path lang.Path) (*decoder.PathContext, error) {
	pathReader := &fdecoder.PathReader{
		StateReader: f.Store,
		RootReader:  f.rootFeature,
	}

	return pathReader.PathContext(path)
}

func (f *ModulesFeature) Paths(ctx context.Context) []lang.Path {
	pathReader := &fdecoder.PathReader{
		StateReader: f.Store,
		RootReader:  f.rootFeature,
	}

	return pathReader.Paths(ctx)
}

func (f *ModulesFeature) DeclaredModuleCalls(modPath string) (map[string]tfmod.DeclaredModuleCall, error) {
	return f.Store.DeclaredModuleCalls(modPath)
}

func (f *ModulesFeature) ProviderRequirements(modPath string) (tfmod.ProviderRequirements, error) {
	mod, err := f.Store.ModuleRecordByPath(modPath)
	if err != nil {
		return nil, err
	}

	return mod.Meta.ProviderRequirements, nil
}

func (f *ModulesFeature) CoreRequirements(modPath string) (version.Constraints, error) {
	mod, err := f.Store.ModuleRecordByPath(modPath)
	if err != nil {
		return nil, err
	}

	return mod.Meta.CoreRequirements, nil
}

func (f *ModulesFeature) ModuleInputs(modPath string) (map[string]tfmod.Variable, error) {
	mod, err := f.Store.ModuleRecordByPath(modPath)
	if err != nil {
		return nil, err
	}

	return mod.Meta.Variables, nil
}

func (f *ModulesFeature) AppendCompletionHooks(srvCtx context.Context, decoderContext decoder.DecoderContext) {
	h := hooks.Hooks{
		ModStore:       f.Store,
		RegistryClient: f.registryClient,
		Logger:         f.logger,
	}

	decoderContext.CompletionHooks["CompleteLocalModuleSources"] = h.LocalModuleSources
	decoderContext.CompletionHooks["CompleteRegistryModuleSources"] = h.RegistryModuleSources
	decoderContext.CompletionHooks["CompleteRegistryModuleVersions"] = h.RegistryModuleVersions
}

func (f *ModulesFeature) Diagnostics(path string) diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()

	mod, err := f.Store.ModuleRecordByPath(path)
	if err != nil {
		return diags
	}

	for source, dm := range mod.ModuleDiagnostics {
		diags.Append(source, dm.AutoloadedOnly().AsMap())
	}

	return diags
}

// MetadataReady checks if a given module exists and if it's metadata has been
// loaded. We need the metadata to enable other features like validation for
// variables.
func (f *ModulesFeature) MetadataReady(dir document.DirHandle) (<-chan struct{}, bool, error) {
	if !f.Store.Exists(dir.Path()) {
		return nil, false, fmt.Errorf("%s: record not found", dir.Path())
	}

	return f.Store.MetadataReady(dir)
}
