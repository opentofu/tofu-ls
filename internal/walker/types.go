// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package walker

import "github.com/opentofu/tofu-ls/internal/document"

type DocumentStore interface {
	HasOpenDocuments(dirHandle document.DirHandle) (bool, error)
}
