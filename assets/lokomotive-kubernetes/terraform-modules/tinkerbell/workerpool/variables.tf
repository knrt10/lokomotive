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

variable "kubeconfig" {
  type        = string
  description = "Content of kubelet's kubeconfig file"
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

variable "cluster_dns_service_ip" {
  type        = string
  description = "IP address of cluster DNS Service. Passed to kubelet as --cluster_dns parameter."
  default     = "10.3.0.10"
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
