module "worker" {
  source = "../../worker"

  kubeconfig = var.kubeconfig

  node_count             = var.node_count
  cluster_dns_service_ip = var.cluster_dns_service_ip
  ssh_keys               = var.ssh_keys
  clc_snippets           = var.clc_snippets
  cluster_domain_suffix  = var.cluster_domain_suffix
  clc_snippet_index      = <<EOF
storage:
  files:
  - path: /etc/hostname
    filesystem: root
    mode: 0644
    contents:
      inline: |
        worker%d
EOF
}

resource "tinkerbell_template" "main" {
  count = var.node_count

  name = "${var.cluster_name}-worker-${count.index}"

  content = templatefile("${path.module}/templates/flatcar-install.tmpl", {
    ignition_config          = module.worker.clc_configs[count.index]
    flatcar_install_base_url = var.flatcar_install_base_url
    machine                  = "${var.cluster_name}_worker_${count.index}"
    os_version               = var.os_version
    os_channel               = var.os_channel
  })
}

resource "tinkerbell_workflow" "main" {
  count = var.node_count

  hardwares = <<EOF
{"${var.cluster_name}_worker_${count.index}": "${var.ip_addresses[count.index]}"}
EOF
  template  = tinkerbell_template.main[count.index].id
}
