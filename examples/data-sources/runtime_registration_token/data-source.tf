data "runtime_registration_token" "current" {}

output "registration_token" {
  value     = data.runtime_registration_token.current.access_token
  sensitive = true
  description = "The current registration token for cluster node registration"
}

