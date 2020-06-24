locals {
  controller_count = length(var.controller_names)
  worker_count     = length(var.worker_names)

  controller_bootstrap_tokens = [
    for index in range(local.controller_count) : {
      token_id     = random_string.bootstrap_token_id_controller[index].result
      token_secret = random_string.bootstrap_token_secret-controller[index].result
    }
  ]

  worker_bootstrap_tokens = [
    for index in range(local.worker_count) : {
      token_id     = random_string.bootstrap_token_id_worker[index].result
      token_secret = random_string.bootstrap_token_secret-worker[index].result
    }
  ]
}

# Generate a cryptographically random token id (public).
resource random_string "bootstrap_token_id_controller" {
  count = local.controller_count

  length  = 6
  upper   = false
  special = false
}

# Generate a cryptographically random token secret.
resource random_string "bootstrap_token_secret-controller" {
  count = local.controller_count

  length  = 16
  upper   = false
  special = false
}

# Generate a cryptographically random token id (public).
resource random_string "bootstrap_token_id_worker" {
  count = local.worker_count

  length  = 6
  upper   = false
  special = false
}

# Generate a cryptographically random token secret.
resource random_string "bootstrap_token_secret-worker" {
  count = local.worker_count

  length  = 16
  upper   = false
  special = false
}
