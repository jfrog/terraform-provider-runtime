---
page_title: "runtime_registration_token Resource - terraform-provider-runtime"
subcategory: ""
description: |-
  Manages registration tokens for cluster node registration.
---

# runtime_registration_token (Resource)

Manages registration tokens for cluster node registration. This resource handles the lifecycle of registration tokens including creation and revocation.

## Security Requirements

Requires a valid Identity Token with Admin privileges.

## Example Usage

### Basic Usage

```terraform
resource "runtime_registration_token" "token" {
  revoked = false
}

output "token" {
  value     = runtime_registration_token.token.access_token
  sensitive = true
}
```

### Token Rotation

To rotate the registration token, set `revoked` to `true`:

```terraform
resource "runtime_registration_token" "token" {
  revoked = true  # This will revoke the current token and generate a new one
}
```

## Schema

### Optional

- `revoked` (Boolean) When set to true, revokes the current token and generates a new one. Use this to rotate tokens.

### Read-Only

- `id` (String) The ID of the registration token resource.
- `access_token` (String, Sensitive) The current, active registration token.

## Import

Registration token can be imported using any identifier:

```terraform
terraform import runtime_registration_token.token registration_token
```

## API Endpoints

This resource uses the following API endpoints:

- `POST /runtime/api/v1/registration_token` - Get current registration token
- `DELETE /runtime/api/v1/registration_token/{token}` - Revoke token and create new one

## Response Codes

### Get Token

- `201` - OK
- `401` - Unauthorized
- `403` - Permission denied
- `500` - Internal server error

### Revoke and Create Token

- `201` - OK
- `401` - Unauthorized
- `403` - Permission denied
- `404` - Not found - The specified token does not exist
- `500` - Internal server error

## Notes

- The resource does not actually revoke the token on `terraform destroy`, as the registration token is a system-level resource that should persist.
- Setting `revoked = true` will revoke the current token and generate a new one in a single operation.
- The `access_token` is marked as sensitive and will not be displayed in plan output.

