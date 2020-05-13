// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tinkerbell

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/go-homedir"

	"github.com/kinvolk/lokomotive/pkg/platform"
	"github.com/kinvolk/lokomotive/pkg/terraform"
)

type config struct {
	AssetDir                        string   `hcl:"asset_dir"`
	Name                            string   `hcl:"name"`
	DNSZone                         string   `hcl:"dns_zone"`
	SSHPublicKeys                   []string `hcl:"ssh_public_keys"`
	ControllerIPAddresses           []string `hcl:"controller_ip_addresses"`
	ControllerCLCSnippets           []string `hcl:"controller_clc_snippets,optional"`
	ControllerFlatcarInstallBaseURL string   `hcl:"controller_flatcar_install_base_url,optional"`

	OSChannel string `hcl:"os_channel,optional"`
	OSVersion string `hcl:"os_version,optional"`

	// Generic options.
	EnableAggregation        bool   `hcl:"enable_aggregation,optional"`
	EnableReporting          bool   `hcl:"enable_reporting,optional"`
	PodCIDR                  string `hcl:"pod_cidr,optional"`
	ServiceCIDR              string `hcl:"service_cidr,optional"`
	ClusterDomainSuffix      string `hcl:"cluster_domain_suffix,optional"`
	CertsValidityPeriodHours int    `hcl:"certs_validity_period_hours,optional"`
	NetworkMTU               int    `hcl:"network_mtu,optional"`

	WorkerPools []workerPool `hcl:"worker_pool,block"`
}

type workerPool struct {
	PoolName string `hcl:"name,label"`

	IPAddresses   []string `hcl:"ip_addresses"`
	SSHPublicKeys []string `hcl:"ssh_public_keys"`

	OSChannel             string   `hcl:"os_channel,optional"`
	OSVersion             string   `hcl:"os_version,optional"`
	FlatcarInstallBaseURL string   `hcl:"flatcar_install_base_url,optional"`
	CLCSnippets           []string `hcl:"clc_snippets,optional"`

	// Generic options.
	//
	// TODO: Should we have more specialized type here to use across the platforms?
	// Maybe structured block? Should we validate them somehow?
	Labels []string `hcl:"labels,optional"`
	Taints []string `hcl:"taints,optional"`
}

func (w *workerPool) Name() string {
	return w.PoolName
}

// init registers tinkerbell as a platform.
func init() { //nolint:gochecknoinits
	platform.Register("tinkerbell", newConfig())
}

func (c *config) LoadConfig(configBody *hcl.Body, evalContext *hcl.EvalContext) hcl.Diagnostics {
	if configBody == nil {
		return hcl.Diagnostics{}
	}

	if diags := gohcl.DecodeBody(*configBody, evalContext, c); len(diags) != 0 {
		return diags
	}

	return c.checkValidConfig()
}

// newConfig returns Tinkerbell default configuration.
func newConfig() *config {
	return &config{
		EnableAggregation: true,
	}
}

// Meta is part of Platform interface and returns common information about the platform configuration.
func (c *config) Meta() platform.Meta {
	nodes := len(c.ControllerIPAddresses)
	for _, workerpool := range c.WorkerPools {
		nodes += len(workerpool.IPAddresses)
	}

	return platform.Meta{
		AssetDir:      c.AssetDir,
		ExpectedNodes: nodes,
	}
}

func (c *config) Initialize(ex *terraform.Executor) error {
	assetDir, err := homedir.Expand(c.AssetDir)
	if err != nil {
		return err
	}

	terraformRootDir := terraform.GetTerraformRootDir(assetDir)

	return createTerraformConfigFile(c, terraformRootDir)
}

func (c *config) Apply(ex *terraform.Executor) error {
	if err := c.Initialize(ex); err != nil {
		return err
	}

	return ex.Apply()
}

func (c *config) Destroy(ex *terraform.Executor) error {
	if err := c.Initialize(ex); err != nil {
		return err
	}

	return ex.Destroy()
}

func createTerraformConfigFile(cfg *config, terraformPath string) error {
	tmplName := "cluster.tf"

	t := template.Must(template.New("cluster.tf").Parse(terraformConfigTmpl))

	path := filepath.Join(terraformPath, tmplName)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}

	if err := t.Execute(f, cfg); err != nil {
		return fmt.Errorf("failed to write template to file %q: %w", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed closing file %q: %w", path, err)
	}

	return nil
}

// checkValidConfig validates cluster configuration.
func (c *config) checkValidConfig() hcl.Diagnostics {
	var d hcl.Diagnostics

	x := []platform.WorkerPool{}
	for i := range c.WorkerPools {
		x = append(x, &c.WorkerPools[i])
	}

	d = append(d, platform.WorkerPoolNamesUnique(x)...)

	return d
}
