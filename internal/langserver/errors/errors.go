// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	e "errors"

	"github.com/opentofu/tofu-ls/internal/tofu/module"
)

func EnrichTfExecError(err error) error {
	if module.IsTerraformNotFound(err) {
		return e.New("Terraform (CLI) is required. " +
			"Please install Terraform or make it available in $PATH")
	}
	return err
}
