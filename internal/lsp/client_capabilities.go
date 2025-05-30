// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"context"
	"errors"

	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

type clientCapsCtxKey struct{}

func SetClientCapabilities(ctx context.Context, caps *lsp.ClientCapabilities) error {
	cc, ok := ctx.Value(clientCapsCtxKey{}).(*lsp.ClientCapabilities)
	if !ok {
		return errors.New("client capabilities not found")
	}

	*cc = *caps
	return nil
}

func WithClientCapabilities(ctx context.Context, caps *lsp.ClientCapabilities) context.Context {
	return context.WithValue(ctx, clientCapsCtxKey{}, caps)
}

func ClientCapabilities(ctx context.Context) (lsp.ClientCapabilities, error) {
	caps, ok := ctx.Value(clientCapsCtxKey{}).(*lsp.ClientCapabilities)
	if !ok {
		return lsp.ClientCapabilities{}, errors.New("client capabilities not found")
	}

	return *caps, nil
}
