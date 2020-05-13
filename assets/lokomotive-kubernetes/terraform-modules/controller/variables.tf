variable "dns_zone" {
  type = string
}

variable "bootkube_rkt_extra_args" {
  description = "Extra parameters to pass to bootkube rkt container"
  type        = list(string)
  default     = []
}

variable "bootkube_image_name" {
  description = "Docker image name to use for rkt container running bootkube"
  type        = string
  default     = "quay.io/kinvolk/bootkube"
}

variable "bootkube_image_tag" {
  description = "Docker image tag to use for rkt container running bootkube"
  type        = string
  default     = "v0.14.0-helm-amd64"
}

# Duplicated variable.
variable "cluster_name" {
  description = "Cluster name"
  type        = string
}

variable "cluster_domain_suffix" {
  type        = string
  description = "Cluster domain suffix. Passed to kubelet as --cluster_domain flag."
  default     = "cluster.local"
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

variable "clc_snippet_index" {
  description = "CLC snippet, which will be formatted with index of the controller."
  type        = string
  default     = ""
}

variable "kubelet_image_name" {
  type        = string
  description = "Source of kubelet Docker image"
  default     = "k8s.gcr.io/hyperkube"
}

variable "kubelet_image_tag" {
  type        = string
  description = "Tag for kubelet Docker image"
  default     = "v1.18.2"
}
