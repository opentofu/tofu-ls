// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/tofu-ls/internal/features/modules/state"
	"github.com/opentofu/tofu-ls/internal/registry"
	globalState "github.com/opentofu/tofu-ls/internal/state"
	"github.com/zclconf/go-cty/cty"
)

// moduleVersionsMockResponse represents the shortened response from https://api.opentofu.org/registry/docs/modules/terraform-aws-modules/vpc/aws/index.json
var moduleVersionsMockResponse = `{
  "addr": {
    "display": "terraform-aws-modules/vpc/aws",
    "namespace": "terraform-aws-modules",
    "name": "vpc",
    "target": "aws"
  },
  "description": "Terraform module to create AWS VPC resources 🇺🇦",
  "versions": [
    {
      "id": "v5.21.0",
      "published": "2025-04-21T23:55:13Z"
    },
    {
      "id": "v2.72.0",
      "published": "2021-02-22T19:00:52Z"
    },
    {
      "id": "v1.0.0",
      "published": "2017-09-12T15:53:29Z"
    }
  ],
  "is_blocked": false,
  "popularity": 3080,
  "fork_count": 4525,
  "fork_of": {
    "display": "//",
    "namespace": "",
    "name": "",
    "target": ""
  },
  "upstream_popularity": 0,
  "upstream_fork_count": 0
}`

func TestHooks_RegistryModuleVersions(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	ctx = decoder.WithPath(ctx, lang.Path{
		Path:       tmpDir,
		LanguageID: "opentofu",
	})
	ctx = decoder.WithPos(ctx, hcl.Pos{
		Line:   2,
		Column: 5,
		Byte:   5,
	})
	ctx = decoder.WithFilename(ctx, "main.tf")
	ctx = decoder.WithMaxCandidates(ctx, 3)
	s, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	store, err := state.NewModuleStore(s.ProviderSchemas, s.RegistryModules, s.ChangeStore)
	if err != nil {
		t.Fatal(err)
	}

	regClient := registry.NewClient()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/registry/docs/modules/terraform-aws-modules/vpc/aws/index.json" {
			w.Write([]byte(moduleVersionsMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	regClient.BaseAPIURL = srv.URL
	t.Cleanup(srv.Close)

	h := &Hooks{
		ModStore:       store,
		RegistryClient: regClient,
	}

	err = store.Add(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	metadata := &tfmod.Meta{
		Path: tmpDir,
		ModuleCalls: map[string]tfmod.DeclaredModuleCall{
			"vpc": {
				LocalName:  "vpc",
				SourceAddr: tfaddr.MustParseModuleSource("registry.opentofu.org/terraform-aws-modules/vpc/aws"),
				RangePtr: &hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 1},
					End:      hcl.Pos{Line: 4, Column: 2, Byte: 20},
				},
			},
		},
	}
	err = store.UpdateMetadata(tmpDir, metadata, nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedCandidates := []decoder.Candidate{
		{
			Label:         `"5.21.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"5.21.0"`,
			SortText:      "  0",
		},
		{
			Label:         `"2.72.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"2.72.0"`,
			SortText:      "  1",
		},
		{
			Label:         `"1.0.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"1.0.0"`,
			SortText:      "  2",
		},
	}

	candidates, _ := h.RegistryModuleVersions(ctx, cty.StringVal(""))
	if diff := cmp.Diff(expectedCandidates, candidates); diff != "" {
		t.Fatalf("mismatched candidates: %s", diff)
	}
}
