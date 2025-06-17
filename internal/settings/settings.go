// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/opentofu/tofu-ls/internal/tofu/datadir"
)

type ExperimentalFeatures struct {
	ValidateOnSave        bool `mapstructure:"validateOnSave"`
	PrefillRequiredFields bool `mapstructure:"prefillRequiredFields"`
}

type ValidationOptions struct {
	EnableEnhancedValidation bool `mapstructure:"enableEnhancedValidation" default:"true"`
}

type Indexing struct {
	IgnoreDirectoryNames []string `mapstructure:"ignoreDirectoryNames"`
	IgnorePaths          []string `mapstructure:"ignorePaths"`
}

type OpenTofu struct {
	Path        string `mapstructure:"path"`
	Timeout     string `mapstructure:"timeout"`
	LogFilePath string `mapstructure:"logFilePath"`
}

type Options struct {
	CommandPrefix string   `mapstructure:"commandPrefix"`
	Indexing      Indexing `mapstructure:"indexing"`

	// ExperimentalFeatures encapsulates experimental features users can opt into.
	ExperimentalFeatures ExperimentalFeatures `mapstructure:"experimentalFeatures"`

	Validation ValidationOptions `mapstructure:"validation"`

	IgnoreSingleFileWarning bool `mapstructure:"ignoreSingleFileWarning"`

	OpenTofu OpenTofu `mapstructure:"opentofu"`

	XLegacyModulePaths          []string `mapstructure:"rootModulePaths"`
	XLegacyExcludeModulePaths   []string `mapstructure:"excludeModulePaths"`
	XLegacyIgnoreDirectoryNames []string `mapstructure:"ignoreDirectoryNames"`
	XLegacyTofuExecPath         string   `mapstructure:"tofuExecPath"`
	XLegacyTofuExecTimeout      string   `mapstructure:"tofuExecTimeout"`
	XLegacyTofuExecLogFilePath  string   `mapstructure:"tofuExecLogFilePath"`
}

func (o *Options) Validate() error {
	if o.OpenTofu.Path != "" {
		path := o.OpenTofu.Path
		if !filepath.IsAbs(path) {
			return fmt.Errorf("expected absolute path for tofu binary, got %q", path)
		}
		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("unable to find tofu binary: %s", err)
		}
		if stat.IsDir() {
			return fmt.Errorf("expected a tofu binary, got a directory: %q", path)
		}
	}

	if len(o.Indexing.IgnoreDirectoryNames) > 0 {
		for _, directory := range o.Indexing.IgnoreDirectoryNames {
			if directory == datadir.DataDirName {
				return fmt.Errorf("cannot ignore directory %q", datadir.DataDirName)
			}

			if strings.Contains(directory, string(filepath.Separator)) {
				return fmt.Errorf("expected directory name, got a path: %q", directory)
			}
		}
	}

	return nil
}

type DecodedOptions struct {
	Options    *Options
	UnusedKeys []string
}

func DecodeOptions(input interface{}) (*DecodedOptions, error) {
	var md mapstructure.Metadata
	options := new(Options)

	// We explicitly set the defaults here before decoding the options.
	// If we were to supply a zero value of a type via our input,
	// setting the default afterwards would override it.
	defaults.SetDefaults(options)

	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   &options,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	if err := decoder.Decode(input); err != nil {
		return nil, err
	}

	return &DecodedOptions{
		Options:    options,
		UnusedKeys: md.Unused,
	}, nil
}
