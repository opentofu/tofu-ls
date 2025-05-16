// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
)

type Provider struct {
	Addr    string `json:"addr"`
	Version string `json:"version"`
}

type ProviderVersion struct {
	Version   string                    `json:"version"`
	Platforms []ProviderVersionPlatform `json:"platforms"`
}

type ProviderVersionPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type providerVersionResponse struct {
	Versions []ProviderVersion `json:"versions"`
}

func (c Client) ListPopularProviders(limit int) ([]Provider, error) {
	providers, err := c.listPopularProvidersRaw(limit)
	if err != nil {
		return nil, err
	}
	return filterUnsupportedProviders(providers)
}

func (c Client) listPopularProvidersRaw(limit int) ([]Provider, error) {
	var providers []Provider
	url := fmt.Sprintf("%s/top/providers?limit=%d", c.BaseURL, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected response: %s: %s", resp.Status, string(bodyBytes))
	}

	var response []Provider
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}
	providers = append(providers, response...)

	return providers, err
}

// filterUnsupportedProviders filters out providers that are not supported
// by the current OS and architecture. It uses a worker pool to check each provider in parallel.
// The function returns a slice of supported providers.
func filterUnsupportedProviders(providers []Provider) ([]Provider, error) {
	log.Printf("Filtering providers, based on OS/ARCH support initial providers %d", len(providers))
	regClient := NewRegistryClient()
	workerCount := 20
	// Put all the providers into a channel to be processed
	providersToCheckCh := make(chan Provider, len(providers))
	for _, provider := range providers {
		providersToCheckCh <- provider
	}
	close(providersToCheckCh)

	// Providers that pass the check will be sent to this channel
	supportedProvidersCh := make(chan Provider, len(providers))
	var wg sync.WaitGroup
	wg.Add(workerCount)
	supportedFound := atomic.Int32{}
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for provider := range providersToCheckCh {
				tfAddr, err := tfaddr.ParseProviderSource(provider.Addr)
				if err != nil {
					continue
				}
				lpv, err := regClient.CheckProviderVersionSupported(tfAddr)
				if err != nil {
					log.Printf("Error checking provider address: %s", provider.Addr)
					continue
				}
				v, err := version.NewVersion(provider.Version)
				if err != nil {
					log.Printf("Error parsing provider version: %s", provider.Version)
					continue
				}
				if !ProviderVersionSupportsOsAndArch(*v, lpv.Versions, runtime.GOOS, runtime.GOARCH) {
					log.Printf("Provider %s version %s does not support %s/%s", provider.Addr, provider.Version, runtime.GOOS, runtime.GOARCH)
					continue
				}
				supportedFound.Add(1)
				log.Printf("(%d/%d) Provider %s version %s supports %s/%s", supportedFound.Load(), len(providers), provider.Addr, provider.Version, runtime.GOOS, runtime.GOARCH)
				supportedProvidersCh <- provider
			}
		}()
	}
	// Wait for all workers to finish
	wg.Wait()
	// Close the channel to signal that no more providers will be sent
	close(supportedProvidersCh)
	supportedProvidersList := make([]Provider, 0, len(providers))
	// Collect the supported providers from the channel and make a new slice
	for provider := range supportedProvidersCh {
		supportedProvidersList = append(supportedProvidersList, provider)
	}
	log.Printf("Finished filtering providers, found %d supported providers out of %d", len(supportedProvidersList), len(providers))
	return supportedProvidersList, nil

}

func (c Client) CheckProviderVersionSupported(pAddr tfaddr.Provider) (*providerVersionResponse, error) {
	url := fmt.Sprintf("%s/v1/providers/%s/%s/versions", c.BaseURL, pAddr.Namespace, pAddr.Type)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected response: %s: %s", resp.Status, string(bodyBytes))
	}

	var response providerVersionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return &response, nil
}

func ProviderVersionSupportsOsAndArch(pVersion version.Version, providerVersions []ProviderVersion, os, arch string) bool {
	for _, version := range providerVersions {
		if version.Version != pVersion.String() {
			continue
		}
		for _, platform := range version.Platforms {
			if platform.OS == os &&
				platform.Arch == arch {
				return true
			}
		}
	}

	return false
}
