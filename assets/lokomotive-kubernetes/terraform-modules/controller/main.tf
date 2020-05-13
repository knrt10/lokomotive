locals {
  kubelet_require_kubeconfig = <<EOF
systemd:
  units:
  - name: kubelet.service
    dropins:
    - name: 10-controller.conf
      contents: |
        [Service]
        ConditionPathExists=/etc/kubernetes/kubeconfig
        ExecStartPre=/bin/mkdir -p /etc/kubernetes/checkpoint-secrets
        ExecStartPre=/bin/mkdir -p /etc/kubernetes/inactive-manifests
EOF

  bootkube = templatefile("${path.module}/templates/bootkube.yaml.tmpl", {
    bootkube_rkt_extra_args = var.bootkube_rkt_extra_args
    bootkube_image_name     = var.bootkube_image_name
    bootkube_image_tag      = var.bootkube_image_tag
    kubelet_image_name      = var.kubelet_image_name
    kubelet_image_tag       = var.kubelet_image_tag
  })

  snippets = [
    local.kubelet_require_kubeconfig,
    local.bootkube,
  ]
}

data "ct_config" "config" {
  count = var.node_count

  pretty_print = false

  content = templatefile("${path.module}/templates/node.yaml.tmpl", {
    ssh_keys               = jsonencode(var.ssh_keys)
    cluster_dns_service_ip = var.cluster_dns_service_ip
    cluster_domain_suffix  = var.cluster_domain_suffix
    kubelet_image_name     = var.kubelet_image_name
    kubelet_image_tag      = var.kubelet_image_tag
    kubelet_rkt_extra_args = []
    kubelet_labels = [
      "node.kubernetes.io/master",
      "node.kubernetes.io/controller=true",
    ]
    kubelet_taints = [
      "node-role.kubernetes.io/master=:NoSchedule"
    ]
  })

  snippets = concat(local.snippets, var.clc_snippets, [
    templatefile("${path.module}/templates/etcd.yaml.tmpl", {
      etcd_name   = "etcd${count.index}"
      etcd_domain = data.template_file.etcds[count.index].rendered
      # etcd0=https://cluster-etcd0.example.com,etcd1=https://cluster-etcd1.example.com,...
      etcd_initial_cluster = join(",", [for i, name in data.template_file.etcds.*.rendered : format("etcd%d=https://%s:2380", i, name)])
    }),
    # Allow to pass unique snippets per controller node. For example, to set the hostname.
    format(var.clc_snippet_index, count.index),
  ])
}

data "template_file" "etcds" {
  count = var.node_count

  template = "$${cluster_name}-etcd$${index}.$${dns_zone}"

  vars = {
    index        = count.index
    cluster_name = var.cluster_name
    dns_zone     = var.dns_zone
  }
}
