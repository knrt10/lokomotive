locals {
  worker_bootstrap_tokens = [
    for index in range(var.worker_count) : {
      token_id     = random_string.bootstrap-token-id[index].result
      token_secret = random_string.bootstrap-token-secret[index].result
    }
  ]
}

data "template_file" "bootstrap-kubeconfig" {
  count = var.worker_count

  template = file("${path.module}/cl/bootstrap-kubeconfig.yaml.tmpl")

  vars = {
    token_id     = random_string.bootstrap-token-id[count.index].result
    token_secret = random_string.bootstrap-token-secret[count.index].result
    ca_cert      = var.ca_cert
    server       = "https://${var.apiserver}:6443"
  }
}

# Generate a cryptographically random token id (public).
resource random_string "bootstrap-token-id" {
  count = var.worker_count

  length  = 6
  upper   = false
  special = false
}

# Generate a cryptographically random token secret.
resource random_string "bootstrap-token-secret" {
  count = var.worker_count

  length  = 16
  upper   = false
  special = false
}
