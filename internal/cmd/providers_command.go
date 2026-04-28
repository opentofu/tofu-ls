// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/opentofu/tofu-ls/internal/schemas"
)

type ProvidersBundledCommand struct {
	Ui cli.Ui
	FS fs.ReadDirFS

	jsonOutput bool
}

func (c *ProvidersBundledCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("providers bundled")

	fs.BoolVar(&c.jsonOutput, "json", false, "output the provider list as JSON")

	fs.Usage = func() { c.Ui.Error(c.Help()) }

	return fs
}

func (c *ProvidersBundledCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	providers, err := schemas.ListBundledProviders(c.FS)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing bundled providers: %s", err))
		return 1
	}

	if len(providers) == 0 {
		c.Ui.Output("No bundled providers found.")
		return 0
	}

	if c.jsonOutput {
		jsonOutput, err := json.MarshalIndent(providers, "", "  ")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error marshalling JSON: %s", err))
			return 1
		}
		c.Ui.Output(string(jsonOutput))
		return 0
	}

	for _, p := range providers {
		c.Ui.Output(fmt.Sprintf("%s %s", p.Addr.ForDisplay(), p.Version))
	}
	return 0
}

func (c *ProvidersBundledCommand) Help() string {
	helpText := `
Usage: tofu-ls providers bundled [-json]

` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c *ProvidersBundledCommand) Synopsis() string {
	return "Lists the provider schemas bundled into tofu-ls binary"
}
