terraform {
  required_version = ">= 1.4.0"
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {
  host = var.docker_host != null ? var.docker_host : "unix://${pathexpand("~/.colima/default/docker.sock")}"
}