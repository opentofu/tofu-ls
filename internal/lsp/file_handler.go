// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lsp

import (
	"github.com/opentofu/tofu-ls/internal/document"
	lsp "github.com/opentofu/tofu-ls/internal/protocol"
)

func HandleFromDocumentURI(docUri lsp.DocumentURI) document.Handle {
	return document.HandleFromURI(string(docUri))
}

func DirHandleFromDirURI(dirUri lsp.DocumentURI) document.DirHandle {
	return document.DirHandleFromURI(string(dirUri))
}
