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

var terraformConfigTmpl = `module "controllers" {
  source = "../lokomotive-kubernetes/terraform-modules/tinkerbell"

  # Generic configuration.
  asset_dir    = "../cluster-assets"
  cluster_name = "{{.Name}}"
  ssh_keys     = [
  {{- range .SSHPublicKeys}}
    "{{.}}",
  {{- end}}
  ]

  ip_addresses = [
  {{- range .ControllerIPAddresses}}
    "{{.}}",
  {{- end}}
  ]

  dns_zone = "{{.DNSZone}}"

  {{- if .ControllerFlatcarInstallBaseURL}}
  flatcar_install_base_url = "{{.ControllerFlatcarInstallBaseURL}}"
  {{- end}}

  {{if .ControllerCLCSnippets}}
  clc_snippets = [
  {{- range .ControllerCLCSnippets }}
    <<EOF
{{.}}
EOF
    ,
  {{- end}}
  ]
  {{end}}

  {{- if .OSChannel}}
  os_channel = "{{.OSChannel}}"
  {{- end }}

  {{- if .OSVersion}}
  os_version = "{{.OSVersion}}"
  {{- end }}

  enable_aggregation = {{.EnableAggregation}}

  {{- if .NetworkMTU }}
  network_mtu = {{.NetworkMTU}}
  {{- end }}

  {{- if .PodCIDR }}
  pod_cidr = "{{.PodCIDR}}"
  {{- end }}

  {{- if .ServiceCIDR }}
  service_cidr = "{{.ServiceCIDR}}"
  {{- end }}

  {{- if .ClusterDomainSuffix }}
  cluster_domain_suffix = "{{.ClusterDomainSuffix}}"
  {{- end }}

  enable_reporting = {{.EnableReporting}}

  {{- if .CertsValidityPeriodHours }}
  certs_validity_period_hours = {{.CertsValidityPeriodHours}}
  {{- end }}
}

{{ range $index, $pool := .WorkerPools }}
module "worker-{{ $pool.Name }}" {
  source = "../lokomotive-kubernetes/terraform-modules/tinkerbell/workerpool"

  kubeconfig             = module.controllers.kubeconfig
  cluster_dns_service_ip = module.controllers.cluster_dns_service_ip

  cluster_name = "{{$.Name}}"

  {{if .SSHPublicKeys}}
  ssh_keys = [
  {{range $pool.SSHPublicKeys}}
    "{{.}}",
  {{end}}
  ]
  {{end}}

  {{if $pool.FlatcarInstallBaseURL}}
  flatcar_install_base_url = "{{$pool.FlatcarInstallBaseURL}}"
  {{end}}

  ip_addresses = [
  {{range $pool.IPAddresses}}
    "{{.}}",
  {{end}}
  ]

  {{- if $.ClusterDomainSuffix }}
  cluster_domain_suffix = "{{$.ClusterDomainSuffix}}"
  {{- end }}

  {{if $pool.CLCSnippets}}
  clc_snippets = [
  {{range $pool.CLCSnippets}}
    <<EOF
{{.}}
EOF
    ,
  {{end}}
  ]
  {{end}}

  {{if $pool.Labels}}
  kubelet_labels = [
  {{range $pool.Labels}}
    "{{.}}",
  {{end}}
  ]
  {{end}}

  {{if $pool.Taints}}
  ssh_keys = [
  {{range $pool.Taints}}
    "{{.}}",
  {{end}}
  ]
  {{end}}
}
{{end}}

provider "ct" {
  version = "~> 0.3"
}

provider "local" {
  version = "1.4.0"
}

provider "null" {
  version = "~> 2.1"
}

provider "template" {
  version = "~> 2.1"
}

provider "tls" {
  version = "~> 2.0"
}

provider "tinkerbell" {
  version = "~> 0.0.0"
}

# Stub output, which indicates, that Terraform run at least once.
# Used when checking, if we should ask user for confirmation, when
# applying changes to the cluster.
output "initialized" {
  value = true
}

# values.yaml content for all deployed charts.
output "pod-checkpointer_values" {
  value = module.controllers.pod-checkpointer_values
}

output "kube-apiserver_values" {
  value     = module.controllers.kube-apiserver_values
  sensitive = true
}

output "kubernetes_values" {
  value     = module.controllers.kubernetes_values
  sensitive = true
}

output "kubelet_values" {
  value     = module.controllers.kubelet_values
  sensitive = true
}

output "calico_values" {
  value = module.controllers.calico_values
}
`
