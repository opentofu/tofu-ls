// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
)

func TestGetModuleData(t *testing.T) {
	ctx := context.Background()
	addr, err := tfaddr.ParseModuleSource("puppetlabs/deployment/ec")
	if err != nil {
		t.Fatal(err)
	}

	cons := version.MustConstraints(version.NewConstraint("0.0.8"))

	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/registry/docs/modules/puppetlabs/deployment/ec/index.json" {
			w.Write([]byte(moduleVersionsMockResponse))
			return
		}
		if r.RequestURI == "/registry/docs/modules/puppetlabs/deployment/ec/v0.0.8/index.json" {
			w.Write([]byte(moduleDataMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	t.Cleanup(srv.Close)

	data, err := client.GetModuleData(ctx, addr, cons)
	if err != nil {
		t.Fatal(err)
	}
	expectedData := &ModuleResponse{
		Version:     "v0.0.8",
		PublishedAt: time.Date(2021, time.August, 5, 0, 26, 01, 0, time.UTC),
		Inputs: map[string]Input{
			"autoscale": {
				Type:        "string",
				Description: "Enable autoscaling of elasticsearch",
				Default:     "true",
				Required:    false,
			},
			"ec_stack_version": {
				Type:        "string",
				Description: "Version of Elastic Cloud stack to deploy",
				Default:     "",
				Required:    false,
			},
			"name": {
				Type:        "string",
				Description: "Name of resources",
				Default:     "ecproject",
				Required:    false,
			},
			"traffic_filter_sourceip": {
				Type:        "string",
				Description: "traffic filter source IP",
				Default:     "",
				Required:    false,
			},
			"ec_region": {
				Type:        "string",
				Description: "cloud provider region",
				Default:     "gcp-us-west1",
				Required:    false,
			},
			"deployment_templateid": {
				Type:        "string",
				Description: "ID of Elastic Cloud deployment type",
				Default:     "gcp-io-optimized",
				Required:    false,
			},
		},
		Outputs: map[string]Output{
			"elasticsearch_password": {
				Description: "elasticsearch password",
			},
			"deployment_id": {
				Description: "Elastic Cloud deployment ID",
			},
			"elasticsearch_version": {
				Description: "Stack version deployed",
			},
			"elasticsearch_cloud_id": {
				Description: "Elastic Cloud project deployment ID",
			},
			"elasticsearch_https_endpoint": {
				Description: "elasticsearch https endpoint",
			},
			"elasticsearch_username": {
				Description: "elasticsearch username",
			},
		},
		Submodules: map[string]Submodule{},
	}
	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("mismatched data: %s", diff)
	}
}

func TestGetMatchingModuleVersion(t *testing.T) {
	ctx := context.Background()
	addr, err := tfaddr.ParseModuleSource("puppetlabs/deployment/ec")
	if err != nil {
		t.Fatal(err)
	}
	cons := version.MustConstraints(version.NewConstraint(">=0.0.7"))
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/registry/docs/modules/puppetlabs/deployment/ec/index.json" {
			w.Write([]byte(moduleVersionsMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	t.Cleanup(srv.Close)

	v, err := client.GetMatchingModuleVersion(ctx, addr, cons)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersion := version.Must(version.NewVersion("0.0.8"))
	if !expectedVersion.Equal(v) {
		t.Fatalf("expected version: %s, given: %s", expectedVersion, v)
	}
}

func TestCancellationThroughContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 50*time.Millisecond)
	t.Cleanup(cancelFunc)

	addr, err := tfaddr.ParseModuleSource("puppetlabs/deployment/ec")
	if err != nil {
		t.Fatal(err)
	}
	cons := version.MustConstraints(version.NewConstraint(">=0.0.7"))
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		if r.RequestURI == "/registry/docs/modules/puppetlabs/deployment/ec/index.json" {
			w.Write([]byte(moduleVersionsMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseAPIURL = srv.URL
	t.Cleanup(srv.Close)

	_, err = client.GetMatchingModuleVersion(ctx, addr, cons)
	e, ok := err.(*url.Error)
	if !ok {
		t.Fatalf("expected error, got %#v", err)
	}

	if e.Err != context.DeadlineExceeded {
		t.Fatalf("expected error: %#v, given: %#v", context.DeadlineExceeded, e.Err)
	}
}
