// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package discovery

type MockDiscovery struct {
	Path string
}

func (d *MockDiscovery) LookPath() (string, error) {
	return d.Path, nil
}
