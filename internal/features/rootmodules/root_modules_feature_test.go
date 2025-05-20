// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rootmodules

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/opentofu/tofu-ls/internal/eventbus"
	"github.com/opentofu/tofu-ls/internal/filesystem"
	globalState "github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
)

func TestRootModulesFeature_TofuVersion(t *testing.T) {
	ss, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	eventBus := eventbus.NewEventBus()
	fs := filesystem.NewFilesystem(ss.DocumentStore)

	type records struct {
		path    string
		version *version.Version
	}

	testCases := []struct {
		name    string
		records []records
		path    string
		version *version.Version
	}{
		{
			"no records",
			[]records{},
			"path/to/module",
			nil,
		},
		{
			"matching record exists",
			[]records{
				{
					"path/to/module",
					version.Must(version.NewVersion("0.12.0")),
				},
			},
			"path/to/module",
			version.Must(version.NewVersion("0.12.0")),
		},
		{
			"no exact match",
			[]records{
				{
					"path/to/module",
					version.Must(version.NewVersion("0.12.0")),
				},
			},
			"path/another/module",
			version.Must(version.NewVersion("0.12.0")),
		},
		{
			"no exact match, multiple records",
			[]records{
				{
					"path/to/module",
					nil,
				},
				{
					"path/another/module",
					nil,
				},
				{
					"root",
					version.Must(version.NewVersion("0.12.0")),
				},
			},
			"path/random/module",
			version.Must(version.NewVersion("0.12.0")),
		},
		{
			"exact match, multiple records",
			[]records{
				{
					"path/to/module",
					nil,
				},
				{
					"path/another/module",
					nil,
				},
				{
					"root",
					version.Must(version.NewVersion("0.12.0")),
				},
			},
			"path/another/module",
			nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			feature, err := NewRootModulesFeature(eventBus, ss, fs, exec.NewMockExecutor(nil))
			if err != nil {
				t.Fatal(err)
			}

			for _, record := range tc.records {
				feature.Store.Add(record.path)
				feature.Store.UpdateTofuAndProviderVersions(record.path, record.version, nil, nil)
			}

			version := feature.TofuVersion(tc.path)

			if diff := cmp.Diff(version, tc.version); diff != "" {
				t.Fatalf("version mismatch for %q: %s", tc.path, diff)
			}
		})
	}
}
