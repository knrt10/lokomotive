output "clc_configs" {
  value = data.ct_config.config.*.rendered
}
