# Terraform version and plugin versions

terraform {
  required_version = ">= 0.13"

  required_providers {
    ct = {
      source  = "poseidon/ct"
      version = "0.6.1"
    }
    local = {
      source  = "hashicorp/local"
      version = "1.4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "2.1.2"
    }
    template = {
      source  = "hashicorp/template"
      version = "2.1.2"
    }
    packet = {
      source  = "packethost/packet"
      version = "3.0.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "2.3.0"
    }
  }
}
