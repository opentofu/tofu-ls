// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package decoder

import (
	"github.com/hashicorp/hcl-lang/validator"
	"github.com/opentofu/tofu-ls/internal/features/modules/decoder/validations"
)

var moduleValidators = []validator.Validator{
	validator.BlockLabelsLength{},
	validator.DeprecatedAttribute{},
	validator.DeprecatedBlock{},
	validator.MaxBlocks{},
	validator.MinBlocks{},
	validations.MissingRequiredAttribute{},
	validator.UnexpectedAttribute{},
	validator.UnexpectedBlock{},
}
