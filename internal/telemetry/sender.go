// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package telemetry

import "context"

type Sender interface {
	SendEvent(ctx context.Context, name string, properties map[string]interface{})
}
