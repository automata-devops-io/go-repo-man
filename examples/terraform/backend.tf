terraform {
  required_version = "~> 1.0"

  backend "http" {
    address        = "https://your.backen.co/api/terraform/repo/go-repo-man/resource/rM-builder"
    lock_address   = "https://your.backend.co/api/terraform/repo/go-repo-man/resource/rM-builder/lock"
    unlock_address = "https://expeditor.chef.io/api/terraform/repo/go-repo-man/resource/rM-builder/lock"

    lock_method   = "POST"
    unlock_method = "DELETE"
  }
}
