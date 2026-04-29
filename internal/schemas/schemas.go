// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"compress/gzip"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
)

//go:embed data
var FS embed.FS

type ProviderSchema struct {
	File    io.Reader
	Version *version.Version
}

type SchemaNotAvailable struct {
	Addr tfaddr.Provider
}

func (e SchemaNotAvailable) Error() string {
	return fmt.Sprintf("embedded schema not available for %s", e.Addr)
}

type BundledProvider struct {
	Addr    tfaddr.Provider
	Version *version.Version
}

// ListBundledProviders returns a list of all bundled providers to the binary on the `data` folder.
func ListBundledProviders(filesystem fs.ReadDirFS) ([]BundledProvider, error) {
	var providers []BundledProvider

	// We only bundle the latest version of each provider, so we expect exactly one
	// version directory per provider. If there are more than one, that's an error
	// If there are zero, then that provider is not bundled and we skip it.
	hostname := tfaddr.DefaultProviderRegistryHost
	namespacePath := path.Join("data", hostname.String())
	namespaces, err := fs.ReadDir(filesystem, namespacePath)
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces {
		if !namespace.IsDir() {
			continue
		}

		typePath := path.Join(namespacePath, namespace.Name())
		types, err := fs.ReadDir(filesystem, typePath)
		if err != nil {
			return nil, err
		}

		for _, providerType := range types {
			if !providerType.IsDir() {
				continue
			}

			versionPath := path.Join(typePath, providerType.Name())
			entries, err := fs.ReadDir(filesystem, versionPath)
			if err != nil {
				return nil, err
			}

			if len(entries) != 1 {
				return nil, fmt.Errorf("%s/%s: expected one version, found %d", namespace.Name(), providerType.Name(), len(entries))
			}

			v, err := version.NewVersion(entries[0].Name())
			if err != nil {
				return nil, err
			}

			providers = append(providers, BundledProvider{
				Addr:    tfaddr.NewProvider(hostname, namespace.Name(), providerType.Name()),
				Version: v,
			})
		}
	}

	return providers, nil
}

func FindProviderSchemaFile(filesystem fs.ReadDirFS, pAddr tfaddr.Provider) (*ProviderSchema, error) {
	providerPath := path.Join("data", pAddr.Hostname.String(), pAddr.Namespace, pAddr.Type)

	entries, err := fs.ReadDir(filesystem, providerPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, SchemaNotAvailable{Addr: pAddr}
		}
		return nil, err
	}

	if len(entries) != 1 {
		return nil, fmt.Errorf("%q: schema not found", pAddr)
	}

	rawVersion := entries[0].Name()

	filePath := path.Join(providerPath, rawVersion, "schema.json.gz")
	file, err := filesystem.Open(filePath)
	if err != nil {
		return nil, err
	}

	version, err := version.NewVersion(rawVersion)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	return &ProviderSchema{
		File:    gzipReader,
		Version: version,
	}, nil
}
