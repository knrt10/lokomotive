variable "dns_zone" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "ip_addresses" {
  type = list(string)
}

variable "flatcar_install_base_url" {
  type    = string
  default = ""
}

variable "os_version" {
  type    = string
  default = ""
}

variable "os_channel" {
  type    = string
  default = ""
}

variable "asset_dir" {
  description = "Path to a directory where generated assets should be placed (contains secrets)"
  type        = string
}

# Required variables.
variable "ssh_keys" {
  type        = list(string)
  description = "List of SSH public keys for user `core`. Each element must be specified in a valid OpenSSH public key format, as defined in RFC 4253 Section 6.6, e.g. 'ssh-rsa AAAAB3N...'."
  default     = []
}

# Optional variables.
variable "node_count" {
  type        = number
  description = "Number of nodes to create."
  default     = 1
}

variable "clc_snippets" {
  type        = list(string)
  description = "Extra CLC snippets to include in the configuration."
  default     = []
}

variable "cluster_domain_suffix" {
  type        = string
  description = "Cluster domain suffix. Passed to kubelet as --cluster_domain flag."
  default     = "cluster.local"
}

variable "network_mtu" {
  description = "CNI interface MTU"
  type        = number
  default     = 1500
}

variable "pod_cidr" {
  description = "CIDR IP range to assign Kubernetes pods"
  type        = string
  default     = "10.2.0.0/16"
}

variable "service_cidr" {
  description = <<EOF
CIDR IP range to assign Kubernetes services.
The 1st IP will be reserved for kube_apiserver, the 10th IP will be reserved for kube-dns.
EOF
  type        = string
  default     = "10.3.0.0/24"
}

variable "enable_reporting" {
  type        = bool
  description = "Enable usage or analytics reporting to upstream component owners (Tigera: Calico)"
  default     = false
}

variable "certs_validity_period_hours" {
  description = "Validity of all the certificates in hours"
  type        = number
  default     = 8760
}

variable "enable_aggregation" {
  description = "Enable the Kubernetes Aggregation Layer (defaults to false, recommended)"
  type        = bool
  default     = true
}
