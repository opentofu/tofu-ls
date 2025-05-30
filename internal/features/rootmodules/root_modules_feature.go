// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rootmodules

import (
	"context"
	"io"
	"log"

	"github.com/hashicorp/go-version"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	"github.com/opentofu/tofu-ls/internal/features/rootmodules/jobs"
	"github.com/opentofu/tofu-ls/internal/features/rootmodules/state"
	globalState "github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
)

// RootModulesFeature groups everything related to root modules. Its internal
// state keeps track of all root modules in the workspace. A root module is
// usually the directory where you would run `terraform init` and where the
// `.terraform` directory and `.terraform.lock.hcl` are located.
//
// The feature listens to events from the EventBus to update its state and
// act on lockfile changes. It also provides methods to query root modules
// for the installed providers, modules, and Terraform version.
type RootModulesFeature struct {
	Store    *state.RootStore
	eventbus *eventbus.EventBus
	stopFunc context.CancelFunc
	logger   *log.Logger

	tfExecFactory exec.ExecutorFactory
	stateStore    *globalState.StateStore
	fs            jobs.ReadOnlyFS
}

func NewRootModulesFeature(eventbus *eventbus.EventBus, stateStore *globalState.StateStore, fs jobs.ReadOnlyFS, tfExecFactory exec.ExecutorFactory) (*RootModulesFeature, error) {
	store, err := state.NewRootStore(stateStore.ChangeStore, stateStore.ProviderSchemas)
	if err != nil {
		return nil, err
	}
	discardLogger := log.New(io.Discard, "", 0)

	return &RootModulesFeature{
		Store:         store,
		eventbus:      eventbus,
		stopFunc:      func() {},
		logger:        discardLogger,
		tfExecFactory: tfExecFactory,
		stateStore:    stateStore,
		fs:            fs,
	}, nil
}

func (f *RootModulesFeature) SetLogger(logger *log.Logger) {
	f.logger = logger
	f.Store.SetLogger(logger)
}

// Start starts the features separate goroutine.
// It listens to various events from the EventBus and performs corresponding actions.
func (f *RootModulesFeature) Start(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	f.stopFunc = cancelFunc

	discoverDone := make(chan struct{}, 10)
	discover := f.eventbus.OnDiscover("feature.rootmodules", discoverDone)

	didOpenDone := make(chan struct{}, 10)
	didOpen := f.eventbus.OnDidOpen("feature.rootmodules", didOpenDone)

	manifestChangeDone := make(chan struct{}, 10)
	manifestChange := f.eventbus.OnManifestChange("feature.rootmodules", manifestChangeDone)

	pluginLockChangeDone := make(chan struct{}, 10)
	pluginLockChange := f.eventbus.OnPluginLockChange("feature.rootmodules", pluginLockChangeDone)

	go func() {
		for {
			select {
			case discover := <-discover:
				// TODO? collect errors
				f.discover(discover.Path, discover.Files)
				discoverDone <- struct{}{}
			case didOpen := <-didOpen:
				// TODO? collect errors
				f.didOpen(didOpen.Context, didOpen.Dir)
				didOpenDone <- struct{}{}
			case manifestChange := <-manifestChange:
				// TODO? collect errors
				f.manifestChange(manifestChange.Context, manifestChange.Dir, manifestChange.ChangeType)
				manifestChangeDone <- struct{}{}
			case pluginLockChange := <-pluginLockChange:
				// TODO? collect errors
				f.pluginLockChange(pluginLockChange.Context, pluginLockChange.Dir)
				pluginLockChangeDone <- struct{}{}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (f *RootModulesFeature) Stop() {
	f.stopFunc()
	f.logger.Print("stopped root modules feature")
}

// InstalledModuleCalls returns the installed module based on the module manifest
func (f *RootModulesFeature) InstalledModuleCalls(modPath string) (map[string]tfmod.InstalledModuleCall, error) {
	return f.Store.InstalledModuleCalls(modPath)
}

// TofuVersion tries to find a modules Tofu version on a best effort basis.
// If a root module exists at the given path, it will return the Terraform
// version of that root module. If not, it will return the version of any
// of the other root modules.
func (f *RootModulesFeature) TofuVersion(modPath string) *version.Version {
	record, err := f.Store.RootRecordByPath(modPath)
	if err != nil {
		if globalState.IsRecordNotFound(err) {
			// TODO try a proximity search to find the closest root module
			record, err = f.Store.RecordWithVersion()
			if err != nil {
				return nil
			}

			return record.TofuVersion
		}

		return nil
	}

	return record.TofuVersion
}

// InstalledProviders returns the installed providers for the given module path
func (f *RootModulesFeature) InstalledProviders(modPath string) (map[tfaddr.Provider]*version.Version, error) {
	record, err := f.Store.RootRecordByPath(modPath)
	if err != nil {
		return nil, err
	}

	return record.InstalledProviders, nil
}

func (f *RootModulesFeature) CallersOfModule(modPath string) ([]string, error) {
	return f.Store.CallersOfModule(modPath)
}

// InstalledModulePath checks the installed modules in the given root module
// for the given normalized source address.
//
// If the module is installed, it returns the path to the module installation
// directory on disk.
func (f *RootModulesFeature) InstalledModulePath(rootPath string, normalizedSource string) (string, bool) {
	record, err := f.Store.RootRecordByPath(rootPath)
	if err != nil {
		return "", false
	}

	dir, ok := record.InstalledModules[normalizedSource]
	if !ok {
		return "", false
	}

	return dir, true
}
