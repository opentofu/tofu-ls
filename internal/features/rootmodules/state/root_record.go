// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package state

import (
	"github.com/hashicorp/go-version"
	"github.com/opentofu/tofu-ls/internal/tofu/datadir"
	op "github.com/opentofu/tofu-ls/internal/tofu/module/operation"
)

// RootRecord contains all information about a module root path, like
// anything related to .terraform/ or .terraform.lock.hcl.
type RootRecord struct {
	path string

	// ProviderSchemaState tracks if we tried loading all provider schemas
	// that this module is using via Terraform CLI
	ProviderSchemaState op.OpState
	ProviderSchemaErr   error

	ModManifest      *datadir.ModuleManifest
	ModManifestErr   error
	ModManifestState op.OpState

	// InstalledModules is a map of normalized source addresses from the
	// manifest to the path of the local directory where the module is installed
	InstalledModules InstalledModules

	TofuVersion      *version.Version
	TofuVersionErr   error
	TofuVersionState op.OpState

	InstalledProviders      InstalledProviders
	InstalledProvidersErr   error
	InstalledProvidersState op.OpState
}

func (m *RootRecord) Copy() *RootRecord {
	if m == nil {
		return nil
	}
	newRecord := &RootRecord{
		path: m.path,

		ProviderSchemaErr:   m.ProviderSchemaErr,
		ProviderSchemaState: m.ProviderSchemaState,

		ModManifest:      m.ModManifest.Copy(),
		ModManifestErr:   m.ModManifestErr,
		ModManifestState: m.ModManifestState,

		// version.Version is practically immutable once parsed
		TofuVersion:      m.TofuVersion,
		TofuVersionErr:   m.TofuVersionErr,
		TofuVersionState: m.TofuVersionState,

		InstalledProvidersErr:   m.InstalledProvidersErr,
		InstalledProvidersState: m.InstalledProvidersState,
	}

	if m.InstalledProviders != nil {
		newRecord.InstalledProviders = make(InstalledProviders, len(m.InstalledProviders))
		for addr, pv := range m.InstalledProviders {
			// version.Version is practically immutable once parsed
			newRecord.InstalledProviders[addr] = pv
		}
	}

	if m.InstalledModules != nil {
		newRecord.InstalledModules = make(InstalledModules, len(m.InstalledModules))
		for source, dir := range m.InstalledModules {
			newRecord.InstalledModules[source] = dir
		}
	}

	return newRecord
}

func (m *RootRecord) Path() string {
	return m.path
}

func newRootRecord(path string) *RootRecord {
	return &RootRecord{
		path:                    path,
		ProviderSchemaState:     op.OpStateUnknown,
		ModManifestState:        op.OpStateUnknown,
		TofuVersionState:        op.OpStateUnknown,
		InstalledProvidersState: op.OpStateUnknown,
	}
}

// NewRootRecordTest is a test helper to create a new Module object
func NewRootRecordTest(path string) *RootRecord {
	return &RootRecord{
		path: path,
	}
}
