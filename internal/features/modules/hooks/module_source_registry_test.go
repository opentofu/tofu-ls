// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/opentofu/tofu-ls/internal/features/modules/state"
	globalState "github.com/opentofu/tofu-ls/internal/state"
	"github.com/zclconf/go-cty/cty"
)

func TestHooks_RegistryModuleSources(t *testing.T) {
	t.Skip("Skipping because currently we don't support modules auto completion")
	ctx := context.Background()

	s, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	store, err := state.NewModuleStore(s.ProviderSchemas, s.RegistryModules, s.ChangeStore)
	if err != nil {
		t.Fatal(err)
	}

	h := &Hooks{
		ModStore: store,
		Logger:   log.New(io.Discard, "", 0),
	}

	tests := []struct {
		name    string
		value   cty.Value
		want    []decoder.Candidate
		wantErr bool
	}{
		{
			"simple search",
			cty.StringVal("aws"),
			[]decoder.Candidate{
				{
					Label:         `"terraform-aws-modules/vpc/aws"`,
					Detail:        "registry",
					Kind:          lang.StringCandidateKind,
					Description:   lang.PlainText("Terraform module which creates VPC resources on AWS"),
					RawInsertText: `"terraform-aws-modules/vpc/aws"`,
				},
				{
					Label:         `"terraform-aws-modules/eks/aws"`,
					Detail:        "registry",
					Kind:          lang.StringCandidateKind,
					Description:   lang.PlainText("Terraform module to create an Elastic Kubernetes (EKS) cluster and associated resources"),
					RawInsertText: `"terraform-aws-modules/eks/aws"`,
				},
			},
			false,
		},
		{
			"empty result",
			cty.StringVal("foo"),
			[]decoder.Candidate{},
			false,
		},
		{
			"auth error",
			cty.StringVal("err"),
			[]decoder.Candidate{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates, err := h.RegistryModuleSources(ctx, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("Hooks.RegistryModuleSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, candidates); diff != "" {
				t.Fatalf("mismatched candidates: %s", diff)
			}
		})
	}
}

func TestHooks_RegistryModuleSourcesCtxCancel(t *testing.T) {
	t.Skip("Skipping because currently we don't support modules auto completion")
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 50*time.Millisecond)
	t.Cleanup(cancelFunc)

	s, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	store, err := state.NewModuleStore(s.ProviderSchemas, s.RegistryModules, s.ChangeStore)
	if err != nil {
		t.Fatal(err)
	}

	h := &Hooks{
		ModStore: store,
		Logger:   log.New(io.Discard, "", 0),
	}

	_, err = h.RegistryModuleSources(ctx, cty.StringVal("aws"))
	e, ok := err.(net.Error)
	if !ok {
		t.Fatalf("expected error, got %#v", err)
	}

	if !strings.Contains(e.Error(), "context deadline exceeded") {
		t.Fatalf("expected error with: %q, given: %q", "context deadline exceeded", e.Error())
	}
}

func TestHooks_RegistryModuleSourcesIgnore(t *testing.T) {
	ctx := context.Background()

	s, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	store, err := state.NewModuleStore(s.ProviderSchemas, s.RegistryModules, s.ChangeStore)
	if err != nil {
		t.Fatal(err)
	}

	h := &Hooks{
		ModStore: store,
		Logger:   log.New(io.Discard, "", 0),
	}

	tests := []struct {
		name  string
		value cty.Value
		want  []decoder.Candidate
	}{
		{
			"search dot",
			cty.StringVal("."),
			[]decoder.Candidate{},
		},
		{
			"search dot dot",
			cty.StringVal(".."),
			[]decoder.Candidate{},
		},
		{
			"local module",
			cty.StringVal("../aws"),
			[]decoder.Candidate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates, err := h.RegistryModuleSources(ctx, tt.value)

			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, candidates); diff != "" {
				t.Fatalf("mismatched candidates: %s", diff)
			}
		})
	}
}
