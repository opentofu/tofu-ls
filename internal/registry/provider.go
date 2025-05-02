// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Provider struct {
	Addr    string `json:"addr"`
	Version string `json:"version"`
}

func (c Client) ListProviders() ([]Provider, error) {
	var providers []Provider
	url := fmt.Sprintf("%s/top/providers?limit=500", c.BaseURL)
	fmt.Printf("using URL %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
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

	return providers, nil
}
