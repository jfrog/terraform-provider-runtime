---
page_title: "runtime_registration_token Data Source - terraform-provider-runtime"
subcategory: ""
description: |-
  Retrieves the current, active registration token.
---

# runtime_registration_token (Data Source)

Retrieves the current, active registration token. This token is used to register new cluster nodes.

## Security Requirements

Requires a valid Identity Token with Admin privileges.

## Example Usage

```terraform
data "runtime_registration_token" "current" {}

output "registration_token" {
  value     = data.runtime_registration_token.current.access_token
  sensitive = true
}
```

## Schema

### Read-Only

- `access_token` (String, Sensitive) The current, active registration token.

## API Endpoint

This data source uses the following API endpoint:

- `POST /runtime/api/v1/registration_token`

## Response Codes

- `201` - OK
- `401` - Unauthorized
- `403` - Permission denied
- `500` - Internal server error

