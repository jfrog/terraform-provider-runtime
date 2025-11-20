# Basic usage - retrieve and manage the registration token
resource "runtime_registration_token" "token" {
  revoked = false
}

output "token" {
  value     = runtime_registration_token.token.access_token
  sensitive = true
  description = "The managed registration token"
}

# Token rotation example - set revoked to true to rotate
resource "runtime_registration_token" "rotated" {
  revoked = true  # This will revoke the current token and generate a new one
}

output "rotated_token" {
  value     = runtime_registration_token.rotated.access_token
  sensitive = true
  description = "The newly generated registration token after rotation"
}

