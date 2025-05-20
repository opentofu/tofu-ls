// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package errors

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

type UnsupportedTofuVersion struct {
	Component   string
	Version     string
	Constraints version.Constraints
}

func (utv *UnsupportedTofuVersion) Error() string {
	msg := "OpenTofu version is not supported"
	if utv.Version != "" {
		msg = fmt.Sprintf("OpenTofu version %s is not supported", utv.Version)
	}

	if utv.Component != "" {
		msg += fmt.Sprintf(" in %s", utv.Component)
	}

	if utv.Constraints != nil {
		msg += fmt.Sprintf(" (supported: %s)", utv.Constraints.String())
	}

	return msg
}

func (utv *UnsupportedTofuVersion) Is(err error) bool {
	te, ok := err.(*UnsupportedTofuVersion)
	if !ok {
		return false
	}

	return te.Version == utv.Version
}
