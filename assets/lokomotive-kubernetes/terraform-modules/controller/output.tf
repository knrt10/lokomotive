output "etcd_servers" {
  value = data.template_file.etcds.*.rendered
}
