# Terraform Provider for JFrog Runtime

## Quick Start

Create a new Terraform file with `runtime` resource:

### HCL Example

```terraform
# Required for Terraform 1.0 and later
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

# Data source to retrieve the current registration token
data "runtime_registration_token" "current" {}

output "current_token" {
  value     = data.runtime_registration_token.current.access_token
  sensitive = true
}

# Resource to manage registration tokens
resource "runtime_registration_token" "token" {
  rotate = false  # Set to true to delete current token and create a new one
}

output "managed_token" {
  value     = runtime_registration_token.token.access_token
  sensitive = true
}

output "last_rotated" {
  value = runtime_registration_token.token.last_rotated
}
```

Initialize Terraform:
```sh
$ terraform init
```

Plan (or Apply):
```sh
$ terraform plan
```

Detailed documentation of the resource and attributes will be available on [Terraform Registry](https://registry.terraform.io/providers/jfrog/runtime/latest/docs).

## Resources and Data Sources

### Data Sources

- **runtime_registration_token**: Retrieves the current, active registration token used to register new cluster nodes.

### Resources

- **runtime_registration_token**: Manages registration tokens for cluster node registration. This resource handles the lifecycle of registration tokens including creation and rotation.

## Requirements

- Terraform 1.0+
- JFrog Platform with Runtime API access
- Identity Token with Admin privileges

## Authentication

The provider supports the following authentication methods:

1. **Access Token** (recommended): Set via `access_token` attribute or `JFROG_ACCESS_TOKEN` environment variable
2. **OIDC**: Configure via `oidc_provider_name` attribute
3. **Terraform Cloud Workload Identity**: Automatic when running in Terraform Cloud with proper configuration

## API Endpoints

This provider uses the following JFrog Runtime API endpoints:

- `POST /runtime/api/v1/registration_token` - Get current registration token
- `DELETE /runtime/api/v1/registration_token/{token}` - Delete current token and create new one

## Versioning

In general, this project follows [semver](https://semver.org/) as closely as we can for tagging releases of the package. We've adopted the following versioning policy:

* We increment the **major version** with any incompatible change to functionality, including changes to the exported Go API surface or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to functionality.
* We increment the **patch version** with any backwards-compatible bug fixes.

## License

Copyright (c) 2025 JFrog.

Apache 2.0 licensed, see [LICENSE][LICENSE] file.

[LICENSE]: ./LICENSE

