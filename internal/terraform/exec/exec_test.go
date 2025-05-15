// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exec_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opentofu/tofu-ls/internal/terraform/exec"
	"github.com/opentofu/tofu-ls/internal/testutils"
)

func TestExec_timeout(t *testing.T) {
	// This test is known to fail under '-race'
	// and similar race conditions are reproducible upstream
	// See https://github.com/hashicorp/terraform-exec/issues/129
	t.Skip("upstream implementation prone to race conditions")

	e := testutils.NewTestingExecutor(t, t.TempDir())
	timeout := 1 * time.Millisecond
	e.SetTimeout(timeout)

	expectedErr := exec.ExecTimeoutError("Version", timeout)

	_, _, err := e.Version(t.Context())
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
	e := testutils.NewTestingExecutor(t, t.TempDir())

	ctx, cancelFunc := context.WithCancel(t.Context())
	cancelFunc()

	expectedErr := exec.ExecCanceledError("Version")

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
