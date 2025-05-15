// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/opentofu/tofudl"
)

func TestExec_timeout(t *testing.T) {
	// This test is known to fail under '-race'
	// and similar race conditions are reproducible upstream
	// See https://github.com/hashicorp/terraform-exec/issues/129
	t.Skip("upstream implementation prone to race conditions")

	e := newExecutor(t)
	timeout := 1 * time.Millisecond
	e.SetTimeout(timeout)

	expectedErr := ExecTimeoutError("Version", timeout)

	_, _, err := e.Version(context.Background())
	if err != nil {
		if errors.Is(err, expectedErr) {
			return
		}

		t.Fatalf("errors don't match.\nexpected: %#v\ngiven:    %#v\n",
			expectedErr, err)
	}

	t.Fatalf("expected timeout error: %#v, given: %#v", expectedErr, err)
}

func TestExec_cancel(t *testing.T) {
	e := newExecutor(t)

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	expectedErr := ExecCanceledError("Version")

	_, _, err := e.Version(ctx)
	if err != nil {
		if errors.Is(err, expectedErr) {
			return
		}

		t.Fatalf("errors don't match.\nexpected: %#v\ngiven:    %#v\n",
			expectedErr, err)
	}

	t.Fatalf("expected cancel error: %#v, given: %#v", expectedErr, err)
}

func newExecutor(t *testing.T) TerraformExecutor {
	ctx := context.Background()
	workDir := t.TempDir()

	dl, err := tofudl.New()
	if err != nil {
		log.Fatalf("error when instantiating tofudl %s", err)
	}

	binary, err := dl.Download(ctx)
	if err != nil {
		log.Fatalf("error when downloading %s", err)
	}

	execPath := filepath.Join(workDir, "tofu")
	// Windows executable case
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, binary, 0755); err != nil {
		log.Fatalf("error when writing the file %s: %s", execPath, err)
	}

	t.Cleanup(func() {
		if err := os.Remove(execPath); err != nil {
			t.Fatal(err)
		}
	})

	e, err := NewExecutor(workDir, execPath)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TempDir(t *testing.T) string {
	tmpDir := filepath.Join(os.TempDir(), "tofu-ls", t.Name())

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
	})
	return tmpDir
}
