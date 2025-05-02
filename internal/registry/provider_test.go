// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListProviders(t *testing.T) {
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[
    {"addr":"hashicorp/aws","version":"v6.0.0-beta1","popularity":10291},
    {"addr":"hashicorp/azurerm","version":"v4.27.0","popularity":4722},
    {"addr":"hashicorp/google","version":"v6.34.0","popularity":2466},
    {"addr":"telmate/proxmox","version":"v3.0.1-rc8","popularity":2399}
]`))
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseURL = srv.URL
	client.ProviderPageSize = 2
	t.Cleanup(srv.Close)

	providers, err := client.ListProviders()
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
