terraform {
  required_providers {
    runtime = {
      source  = "jfrog/runtime"
      version = "~> 1.0"
    }
  }
}

# Configure with explicit credentials
provider "runtime" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-identity-token"
}

# Configure with environment variables
# JFROG_URL and JFROG_ACCESS_TOKEN
provider "runtime" {
  # url and access_token will be read from environment variables
}

# Configure with OIDC
provider "runtime" {
  url                = "https://myinstance.jfrog.io"
  oidc_provider_name = "my-oidc-provider"
}

