terraform {
  required_providers {
    runtime = {
      source  = "jfrog/runtime"
      version = "1.0.0"
    }
  }
}

provider "runtime" {
  url = "https://myinstance.jfrog.io"
  // supply JFROG_ACCESS_TOKEN (Identity Token with Admin privileges) as env var
}

# Resource to manage registration tokens
resource "runtime_registration_token" "token" {
  rotate = false # Set to true to delete current token and create a new one
}

