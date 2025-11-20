---
page_title: "Provider: Runtime"
description: |-
  The JFrog Runtime provider is used to interact with the JFrog Runtime API features.
---

# Runtime Provider

The [JFrog](https://jfrog.com/) Runtime provider is used to interact with the Runtime API features. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```terraform
# Configure the Runtime Provider
provider "runtime" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-identity-token-with-admin-privileges"
}

# Retrieve the current registration token
data "runtime_registration_token" "current" {}

# Manage registration tokens
resource "runtime_registration_token" "token" {
  revoked = false
}
```

## Authentication

The Runtime provider supports several ways to authenticate:

### Identity Token (Recommended)

The provider requires an Identity Token with Admin privileges to interact with the Runtime API.

```terraform
provider "runtime" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-identity-token"
}
```

You can also use environment variables:

```bash
export JFROG_URL="https://myinstance.jfrog.io"
export JFROG_ACCESS_TOKEN="my-identity-token"
```

### OIDC Provider

```terraform
provider "runtime" {
  url                = "https://myinstance.jfrog.io"
  oidc_provider_name = "my-oidc-provider"
}
```

### Terraform Cloud Workload Identity

When running in Terraform Cloud with proper OIDC configuration, the provider can automatically use workload identity tokens.

## Argument Reference

* `url` - (Optional) JFrog Platform URL. This can also be sourced from the `JFROG_URL` environment variable.
* `access_token` - (Optional) Identity Token with Admin privileges. This can also be sourced from the `JFROG_ACCESS_TOKEN` environment variable.
* `oidc_provider_name` - (Optional) OIDC provider name. See [Configure an OIDC Integration](https://jfrog.com/help/r/jfrog-platform-administration-documentation/configure-an-oidc-integration) for more details.
* `tfc_credential_tag_name` - (Optional) Terraform Cloud Workload Identity Token tag name. Use for generating multiple TFC workload identity tokens.

## Runtime API Endpoints

This provider uses the following JFrog Runtime API endpoints:

* `POST /runtime/api/v1/registration_token` - Get current registration token
* `DELETE /runtime/api/v1/registration_token/{token}` - Revoke token and create new one

