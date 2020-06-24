locals {
  controller_bootstrap_tokens = [
    for index in range(var.controller_count) : {
      token_id     = random_string.bootstrap_token_id[index].result
      token_secret = random_string.bootstrap_token_secret[index].result
    }
  ]
}

# Generate a cryptographically random token id (public).
resource random_string "bootstrap_token_id" {
  count = var.controller_count

  length  = 6
  upper   = false
  special = false
}

# Generate a cryptographically random token secret.
resource random_string "bootstrap_token_secret" {
  count = var.controller_count

  length  = 16
  upper   = false
  special = false
}
