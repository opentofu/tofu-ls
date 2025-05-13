// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build generate
// +build generate

package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/tofu-exec/tfexec"
	lsctx "github.com/opentofu/tofu-ls/internal/context"
	"github.com/opentofu/tofu-ls/internal/registry"
	"github.com/opentofu/tofudl"
)

var terraformVersion = version.MustConstraints(version.NewConstraint("~> 1.0"))

type Provider struct {
	ID      string
	Addr    tfaddr.Provider
	Version *version.Version
}

func main() {
	os.Exit(func() int {
		if err := gen(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return 0
	}())
}

func gen() error {
	ctx := context.Background()

	ctx, cancelFunc := lsctx.WithSignalCancel(context.Background(), log.Default(),
		os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	providers := make([]Provider, 0)

	client := registry.NewClient()
	log.Println("fetching from registry")
	listOfProviders, err := client.ListProviders()
	if err != nil {
		return err
	}
	log.Printf("fetched providers: %d", len(listOfProviders))
	for _, p := range listOfProviders {
		pAddr, err := tfaddr.ParseProviderSource(p.Addr)
		if err != nil {
			// TODO: Better error handling
			fmt.Printf("error processing %s\n", pAddr)
			continue
		}

		ver, err := version.NewVersion(p.Version)
		if err != nil {
			// TODO: Better error handling
			fmt.Printf("error processing version of %s\n", pAddr)
			continue
		}

		providers = append(providers, Provider{
			ID: p.Addr,
			Addr: tfaddr.NewProvider(
				tfaddr.DefaultProviderRegistryHost,
				pAddr.Namespace,
				pAddr.Type,
			),
			Version: ver,
		})
	}

	// find or install Terraform
	log.Println("ensuring tofu is installed")
	tempDir, err := os.MkdirTemp("", "tofuinstall")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dl, err := tofudl.New()
	if err != nil {
		log.Fatalf("error when instantiating tofudl %s", err)
	}

	binary, err := dl.Download(ctx)
	if err != nil {
		log.Fatalf("error when downloading %s", err)
	}

	execPath := filepath.Join(tempDir, "tofu")
	// Windows executable case
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, binary, 0755); err != nil {
		log.Fatalf("error when writing the file %s: %s", execPath, err)
	}

	// log version
	tf, err := tfexec.NewTofu(tempDir, execPath)
	if err != nil {
		return err
	}
	coreVersion, _, err := tf.Version(ctx, true)
	if err != nil {
		return err
	}
	log.Printf("using Terraform %s (%s)", coreVersion, execPath)

	workspacePath, err := filepath.Abs("gen-workspace")
	if err != nil {
		return err
	}
	dataDirPath, err := filepath.Abs("data")
	if err != nil {
		return err
	}

	// remove data from previous run
	err = os.RemoveAll(workspacePath)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dataDirPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			// ensure that git.keep is kept in place
			continue
		}
		err = os.RemoveAll(filepath.Join(dataDirPath, entry.Name()))
		if err != nil {
			return err
		}
	}

	cacheDirPath, err := filepath.Abs("tf-plugin-cache")
	if err != nil {
		return err
	}
	err = os.MkdirAll(cacheDirPath, 0755)
	if err != nil {
		return err
	}
	log.Printf("Terraform plugin cache will be stored at %s", cacheDirPath)

	// install each provider and obtain schema for it
	providerChan := make(chan Inputs)
	go func() {
		for _, p := range providers {
			providerChan <- Inputs{
				TerraformExecPath: execPath,
				WorkspacePath:     workspacePath,
				DataDirPath:       dataDirPath,
				CacheDirPath:      cacheDirPath,
				CoreVersion:       coreVersion,
				Provider:          p,
				ProviderVersion:   p.Version,
			}
		}
		close(providerChan)
	}()

	registryClient := registry.NewRegistryClient()

	var workerWg sync.WaitGroup
	workerCount := runtime.NumCPU()
	log.Printf("worker count: %d", workerCount)
	workerWg.Add(workerCount)
	for i := 1; i <= workerCount; i++ {
		go func(i int) {
			defer workerWg.Done()
			for input := range providerChan {
				log.Printf("%s: obtaining schema ...", input.Provider.Addr.ForDisplay())
				details, err := schemaForProvider(ctx, registryClient, input)

				if err != nil {
					log.Printf("%s: %s", input.Provider.Addr.ForDisplay(), err)
					continue
				}

				log.Printf("%s: obtained schema for %s (%db raw / %db compressed); tofu init: %s",
					input.Provider.Addr.ForDisplay(), input.ProviderVersion,
					details.RawSize, details.CompressedSize, details.InitElapsedTime)
			}
		}(i)
	}
	workerWg.Wait()

	return nil
}

type Inputs struct {
	TerraformExecPath string
	WorkspacePath     string
	DataDirPath       string
	CacheDirPath      string
	CoreVersion       *version.Version
	Provider          Provider
	ProviderVersion   *version.Version
}

type Outputs struct {
	Version         string
	RawSize         int
	CompressedSize  int64
	InitElapsedTime time.Duration
}

func schemaForProvider(ctx context.Context, client registry.Client, input Inputs) (*Outputs, error) {
	var pVersion *version.Version
	pVersion = input.CoreVersion

	wd := filepath.Join(input.WorkspacePath,
		input.Provider.Addr.Hostname.String(),
		input.Provider.Addr.Namespace,
		input.Provider.Addr.Type,
		pVersion.String())
	err := os.MkdirAll(wd, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create workspace dir: %w", err)
	}

	dataDir := filepath.Join(input.DataDirPath,
		input.Provider.Addr.Hostname.String(),
		input.Provider.Addr.Namespace,
		input.Provider.Addr.Type,
		pVersion.String())
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create data dir: %w", err)
	}

	type templateData struct {
		TerraformVersion string
		LocalName        string
		Source           string
		Version          string
	}
	tmpl, err := template.New("providers").Parse(`terraform {
  required_version = "{{ .TerraformVersion }}"
  required_providers {
    {{ .LocalName }} = {
      source  = "{{ .Source }}"
      {{ with .Version }}version = "{{ . }}"{{ end }}
    }
  }
}
`)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template: %w", err)
	}

	versionFilePath := filepath.Join(wd, "versions.tf")
	configFile, err := os.Create(versionFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create config file: %w", err)
	}

	err = tmpl.Execute(configFile, templateData{
		TerraformVersion: terraformVersion.String(),
		LocalName:        input.Provider.Addr.Type,
		Source:           input.Provider.Addr.ForDisplay(),
		Version:          input.ProviderVersion.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	configFile.Close()

	tf, err := tfexec.NewTofu(wd, input.TerraformExecPath)
	if err != nil {
		return nil, err
	}

	// See https://github.com/hashicorp/terraform-exec/issues/337
	// Terraform would refuse to init any provider otherwise
	// and some providers refuse to give schemas or break
	// the gRPC protocol for some mysterious reason
	env := make(map[string]string, 0)
	for _, rawKeyPair := range os.Environ() {
		parts := strings.Split(rawKeyPair, "=")
		if parts[0] == "" {
			// For unknown reasons on Windows there can be some odd variables
			// such as "=::=::\\", "=C:=C:\\path" or "=ExitCode=00000000"
			// which we ignore here
			continue
		}
		env[parts[0]] = os.Getenv(parts[0])
	}
	// This is to help keep paths short, esp. on Windows
	// (260 characters by default)
	// See https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#maximum-path-length-limitation
	// and also to avoid embedding the provider binaries
	env["TF_PLUGIN_CACHE_DIR"] = input.CacheDirPath

	tf.SetEnv(env)

	var initElapsed time.Duration
	if !input.Provider.Addr.IsBuiltIn() {
		initElapsed, err = retryInit(ctx, tf, input.Provider.Addr.ForDisplay(), 0)
		if err != nil {
			return nil, err
		}

		_, pVersions, err := tf.Version(ctx, true)
		if err != nil {
			return nil, err
		}

		pv, ok := pVersions[input.Provider.Addr.String()]
		if !ok {
			return nil, fmt.Errorf("provider version not found for %q", input.Provider.Addr.ForDisplay())
		}
		if !pv.Equal(input.ProviderVersion) {
			return nil, fmt.Errorf("expected provider version %s to match %s", pv, pVersion)
		}

		lpv, err := client.CheckProviderVersionSupported(input.Provider.Addr)

		if !registry.ProviderVersionSupportsOsAndArch(*input.ProviderVersion, lpv.Versions, runtime.GOOS, runtime.GOARCH) {
			return nil, fmt.Errorf("version %s does not support %s/%s", input.ProviderVersion, runtime.GOOS, runtime.GOARCH)
		}
	}

	// TODO upstream change to have tfexec write to file directly instead of unmarshal/remarshal
	ps, err := retryProviderSchema(ctx, tf, input.Provider.Addr.ForDisplay(), 0)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(dataDir, "schema.json.gz"))
	if err != nil {
		return nil, fmt.Errorf("failed to create schema file: %w", err)
	}
	var rawJson bytes.Buffer
	err = json.NewEncoder(&rawJson).Encode(ps)
	if err != nil {
		return nil, fmt.Errorf("failed to encode schema file: %w", err)
	}
	gzw := gzip.NewWriter(f)
	_, err = gzw.Write(rawJson.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write compressed file: %w", err)
	}
	err = gzw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close compressed file: %w", err)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to check schema file: %w", err)
	}

	return &Outputs{
		Version:         pVersion.String(),
		RawSize:         rawJson.Len(),
		CompressedSize:  fi.Size(),
		InitElapsedTime: initElapsed,
	}, nil
}

// retryInit runs "terraform init" and attempts to retry
// on known (typically network-related) transient errors
func retryInit(ctx context.Context, tf *tfexec.Tofu, fullName string, retried int) (time.Duration, error) {
	maxRetries := 5
	backoffPeriod := 2 * time.Second

	startTime := time.Now()
	err := tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		retried++
		if retried >= maxRetries {
			timeElapsed := time.Now().Sub(startTime)
			return timeElapsed, fmt.Errorf("%s: final error after 5 retries: %w", fullName, err)
		}

		if shortErr, ok := initErrorIsRetryable(err); ok {
			log.Printf("%s: %s", fullName, err)
			log.Printf("%s: will retry init (attempt %d) in %s due to %s", fullName, retried, backoffPeriod, shortErr)
			time.Sleep(backoffPeriod)
			return retryInit(ctx, tf, fullName, retried)
		}
		return 0, err
	}

	timeElapsed := time.Now().Sub(startTime)
	return timeElapsed, nil
}

func retryProviderSchema(ctx context.Context, tf *tfexec.Tofu, fullName string, retried int) (*tfjson.ProviderSchemas, error) {
	maxRetries := 5
	backoffPeriod := 2 * time.Second

	ps, err := tf.ProvidersSchema(ctx)
	if err != nil {
		retried++
		if retried >= maxRetries {
			return nil, fmt.Errorf("%s: final error after 5 retries: %w", fullName, err)
		}

		// It's unclear why, but some providers just panic
		// (especially on Windows) on the first attempt.
		// While this shouldn't be happening at all,
		// retrying is the easiest workaround here.
		if strings.Contains(err.Error(), "Failed to load plugin schemas") {
			log.Printf("%s: %s", fullName, err)
			log.Printf("%s: will retry provider schema (attempt %d) in %s", fullName, retried, backoffPeriod)
			time.Sleep(backoffPeriod)
			return retryProviderSchema(ctx, tf, fullName, retried)
		}

		return nil, fmt.Errorf("%s: final error after 5 retries: %w", fullName, err)
	}
	return ps, nil
}

func initErrorIsRetryable(err error) (string, bool) {
	if strings.Contains(err.Error(), "i/o timeout") {
		return "i/o timeout", true
	}
	if strings.Contains(err.Error(), "request canceled while waiting for connection") {
		return "connection timeout", true
	}
	if strings.Contains(err.Error(), "handshake timeout") {
		return "handshake timeout", true
	}
	if strings.Contains(err.Error(), "no route to host") {
		return "no route to host", true
	}
	if strings.Contains(err.Error(), "context deadline exceeded") {
		return "context deadline exceeded", true
	}
	if strings.Contains(err.Error(), "503 Service Unavailable") {
		return "503 Service Unavailable", true
	}
	return "", false
}
