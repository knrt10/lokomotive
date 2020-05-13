module "controller" {
  source = "../controller"

  cluster_name           = var.cluster_name
  dns_zone               = var.dns_zone
  node_count             = var.node_count
  cluster_dns_service_ip = module.bootkube.cluster_dns_service_ip
  ssh_keys               = var.ssh_keys
  clc_snippets           = var.clc_snippets
  cluster_domain_suffix  = var.cluster_domain_suffix
  bootkube_rkt_extra_args = [
    # So /etc/hosts changes passed via CLC snippets have effect in bootkube rkt container.
    # This allows to workaround a requirement of DNS server resolving etcd DNS names etc.
    "--hosts-entry=host",
  ]
  clc_snippet_index = <<EOF
storage:
  files:
  - path: /etc/hostname
    filesystem: root
    mode: 0644
    contents:
      inline: |
        controller%d
EOF
}

resource "tinkerbell_template" "main" {
  count = var.node_count

  name = "${var.cluster_name}-controller-${count.index}"

  content = templatefile("${path.module}/templates/flatcar-install.tmpl", {
    ignition_config          = module.controller.clc_configs[count.index]
    flatcar_install_base_url = var.flatcar_install_base_url
    machine                  = "${var.cluster_name}_controller_${count.index}"
    os_version               = var.os_version
    os_channel               = var.os_channel
  })
}

resource "tinkerbell_workflow" "main" {
  count = var.node_count

  hardwares = <<EOF
{"${var.cluster_name}_controller_${count.index}": "${var.ip_addresses[count.index]}"}
EOF
  template  = tinkerbell_template.main[count.index].id
}
