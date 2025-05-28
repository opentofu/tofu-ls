// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
)

func TestListPopularProvidersFiltered(t *testing.T) {
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/top/providers?limit=5" {
			w.Write([]byte(`[
				{"addr":"hashicorp/aws","version":"v1.0.0","popularity":10291},
				{"addr":"hashicorp/azurerm","version":"v1.0.0","popularity":4722},
			    {"addr":"hashicorp/google","version":"v1.0.0","popularity":2466},
				{"addr":"telmate/proxmox","version":"v1.0.0","popularity":2399},
				{"addr":"hashicorp/unsupported","version":"v1.0.0","popularity":100}
			]`))
			return
		} else if !strings.Contains(r.RequestURI, "unsupported") {
			w.Write([]byte(fmt.Sprintf(`{
				"versions": [
					{
						"version": "1.0.0",
						"protocols": ["5.0"],
						"platforms": [
							{"os": "%s", "arch": "%s"}
						]
					}
				]
			}`, runtime.GOOS, runtime.GOARCH)))
			return
		} else if strings.Contains(r.RequestURI, "unsupported") {
			w.Write([]byte(`{
				"versions": [
					{
						"version": "v1.0.0",
						"protocols": ["5.0"],
						"platforms": [
						]
					}
				]
			}`))
			return
		}

		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	client.BaseRegistryURL = srv.URL
	t.Cleanup(srv.Close)

	providers, err := client.ListPopularProviders(5)
	if err != nil {
		t.Fatal(err)
	}

	expectedProviders := []Provider{
		{Addr: "hashicorp/aws", Version: "v1.0.0"},
		{Addr: "hashicorp/azurerm", Version: "v1.0.0"},
		{Addr: "hashicorp/google", Version: "v1.0.0"},
		{Addr: "telmate/proxmox", Version: "v1.0.0"},
	}
	sortFn := func(a, b Provider) int {
		if a.Addr < b.Addr {
			return -1
		}
		if a.Addr > b.Addr {
			return 1
		}
		return 0
	}

	slices.SortFunc(providers, sortFn)
	slices.SortFunc(expectedProviders, sortFn)
	if diff := cmp.Diff(expectedProviders, providers); diff != "" {
		t.Fatalf("unexpected providers: %s", diff)
	}
}

func TestListPopularProviders(t *testing.T) {
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/top/providers?limit=4" {
			w.Write([]byte(`[
				{"addr":"hashicorp/aws","version":"v6.0.0-beta1","popularity":10291},
				{"addr":"hashicorp/azurerm","version":"v4.27.0","popularity":4722},
				{"addr":"hashicorp/google","version":"v6.34.0","popularity":2466},
				{"addr":"telmate/proxmox","version":"v3.0.1-rc8","popularity":2399}
			]`))
			return
		}

		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	t.Cleanup(srv.Close)

	providers, err := client.listPopularProvidersRaw(4)
	if err != nil {
		t.Fatal(err)
	}

	expectedProviders := []Provider{
		{Addr: "hashicorp/aws", Version: "v6.0.0-beta1"},
		{Addr: "hashicorp/azurerm", Version: "v4.27.0"},
		{Addr: "hashicorp/google", Version: "v6.34.0"},
		{Addr: "telmate/proxmox", Version: "v3.0.1-rc8"},
	}
	if diff := cmp.Diff(expectedProviders, providers); diff != "" {
		t.Fatalf("unexpected providers: %s", diff)
	}
}

func TestGetLatestProviderVersion(t *testing.T) {
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/v1/providers/hashicorp/aws/versions" {
			w.Write([]byte(`{
			"versions": [
			{
				"version": "6.0.0-beta1",
				"protocols": [
				"5.0"
				],
				"platforms": [
				{
					"os": "darwin",
					"arch": "amd64"
				},
				{
					"os": "linux",
					"arch": "386"
				},
				{
					"os": "windows",
					"arch": "386"
				}
				]
			}
		]}`))
			return
		}

		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	client.BaseRegistryURL = srv.URL
	t.Cleanup(srv.Close)

	pAddr := tfaddr.NewProvider(tfaddr.DefaultProviderRegistryHost, "hashicorp", "aws")

	resp, err := client.checkProviderVersionSupported(pAddr)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := &providerVersionResponse{
		Versions: []ProviderVersion{
			{
				Version: "6.0.0-beta1",
				Platforms: []ProviderVersionPlatform{
					{
						OS:   "darwin",
						Arch: "amd64",
					},
					{
						OS:   "linux",
						Arch: "386",
					},
					{
						OS:   "windows",
						Arch: "386",
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expectedResponse, resp); diff != "" {
		t.Fatalf("unexpected response: %s", diff)
	}
}

func TestProviderSupport(t *testing.T) {
	ver, err := version.NewVersion("6.0.0-beta1")

	if err != nil {
		t.Fatal(err)
	}

	pVersions := []ProviderVersion{
		{
			Version: "6.0.0-beta1",
			Platforms: []ProviderVersionPlatform{
				{
					OS:   "darwin",
					Arch: "amd64",
				},
				{
					OS:   "linux",
					Arch: "386",
				},
				{
					OS:   "windows",
					Arch: "386",
				},
			},
		},
	}

	result := ProviderVersionSupportsOsAndArch(*ver, pVersions, "linux", "386")
	if !result {
		t.Fatalf("expecting linux 386 to be a supported version")
	}
}
