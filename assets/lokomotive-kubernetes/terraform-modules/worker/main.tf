locals {
  manage_iscsid_service = <<EOF
systemd:
  units:
    - name: iscsid.service
      enable: true
      enabled: true
EOF

  kubeconfig = <<EOF
storage:
  files:
    - path: /etc/kubernetes/kubeconfig
      filesystem: root
      mode: 0644
      contents:
        inline: |
          ${indent(10, var.kubeconfig)}
EOF

  snippets = [
    local.manage_iscsid_service,
    local.kubeconfig,
  ]
}

data "ct_config" "config" {
  count = var.node_count

  pretty_print = false

  content = templatefile("${path.module}/templates/node.yaml.tmpl", {
    ssh_keys               = jsonencode(var.ssh_keys)
    cluster_dns_service_ip = var.cluster_dns_service_ip
    cluster_domain_suffix  = var.cluster_domain_suffix
    kubelet_rkt_extra_args = [
      # Workers should have iscsiadm mounted for storage solutions support.
      "--volume iscsiadm,kind=host,source=/usr/sbin/iscsiadm",
      "--mount volume=iscsiadm,target=/usr/sbin/iscsiadm",
    ]
    # Here we set default labels for worker nodes.
    kubelet_labels = length(var.kubelet_labels) > 0 ? var.kubelet_labels : [
      "node.kubernetes.io/node"
    ]
    kubelet_taints = var.kubelet_taints

    kubelet_image_name = var.kubelet_image_name != "" ? var.kubelet_image_name : null
    kubelet_image_tag  = var.kubelet_image_tag != "" ? var.kubelet_image_tag : null
  })

  snippets = concat(local.snippets, var.clc_snippets, [
    # Allow to pass unique snippets per controller node. For example, to set the hostname.
    format(var.clc_snippet_index, count.index),
  ])
}
