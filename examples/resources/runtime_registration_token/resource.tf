# Basic usage - retrieve and manage the registration token
resource "runtime_registration_token" "token" {
  rotate = false
}

output "token" {
  value     = runtime_registration_token.token.access_token
  sensitive = true
  description = "The managed registration token"
}

output "last_rotated" {
  value     = runtime_registration_token.token.last_rotated
  description = "Timestamp of last token rotation"
}

# Token rotation example - set rotate to true to rotate on each apply
resource "runtime_registration_token" "rotated" {
  rotate = true  # This will delete the current token and create a new one on each apply
}

output "rotated_token" {
  value     = runtime_registration_token.rotated.access_token
  sensitive = true
  description = "The newly generated registration token - rotates on each apply when rotate = true"
}

