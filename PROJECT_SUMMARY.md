# Terraform Provider Runtime - Project Summary

## Overview

Successfully created a new Terraform provider for JFrog Runtime API v1, implementing support for registration token management.

## Features Implemented

### Data Sources

1. **runtime_registration_token**
   - Retrieves the current, active registration token
   - Uses `POST /runtime/api/v1/registration_token`
   - Returns sensitive access_token field

### Resources

1. **runtime_registration_token**
   - Manages registration tokens for cluster node registration
   - Supports token rotation via `rotate` flag
   - Uses two API endpoints:
     - `POST /runtime/api/v1/registration_token` - Get current token
     - `DELETE /runtime/api/v1/registration_token/{token}` - Delete current token and generate a new one
   - Supports import functionality

## Project Structure

```
terraform-provider-runtime/
├── main.go                           # Provider entry point
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksums
├── GNUmakefile                       # Build and test automation
├── LICENSE                           # Apache 2.0 license
├── README.md                         # User documentation
├── CHANGELOG.md                      # Version history
├── CODEOWNERS                        # Code ownership
├── CONTRIBUTIONS.md                  # Contribution guidelines
├── sample.tf                         # Sample Terraform configuration
├── terraform-registry-manifest.json  # Terraform registry metadata
├── .goreleaser.yml                   # Release configuration
├── .gitignore                        # Git ignore patterns
├── pkg/runtime/
│   ├── provider.go                   # Provider implementation
│   ├── data_source_registration_token.go  # Data source implementation
│   └── resource_registration_token.go     # Resource implementation
├── docs/
│   ├── index.md                      # Provider documentation
│   ├── data-sources/
│   │   └── registration_token.md     # Data source documentation
│   └── resources/
│       └── registration_token.md     # Resource documentation
├── examples/
│   ├── provider/
│   │   └── provider.tf               # Provider configuration examples
│   ├── data-sources/
│   │   └── runtime_registration_token/
│   │       └── data-source.tf        # Data source example
│   └── resources/
│       └── runtime_registration_token/
│           └── resource.tf           # Resource example
├── templates/
│   └── index.md.tmpl                 # Documentation template
└── tools/
    └── tools.go                      # Build tools

```

## Provider Configuration

The provider supports multiple authentication methods:

1. **Identity Token (Access Token)** - Via configuration or environment variable
2. **OIDC** - Via OIDC provider name
3. **Terraform Cloud Workload Identity** - Automatic when running in TFC

Example configuration:

```terraform
provider "runtime" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-identity-token-with-admin-privileges"
}
```

## API Endpoints Implemented

### Get Registration Token
- **Method:** POST
- **Endpoint:** `/runtime/api/v1/registration_token`
- **Description:** Retrieves the current, active registration token
- **Response Codes:** 201 (OK), 401 (Unauthorized), 403 (Permission denied), 500 (Internal server error)

### Delete and Generate Registration Token
- **Method:** DELETE
- **Endpoint:** `/runtime/api/v1/registration_token/{token}`
- **Description:** Deletes the specified token and generates a new one
- **Response Codes:** 201 (OK), 401 (Unauthorized), 403 (Permission denied), 404 (Not found), 500 (Internal server error)

## Building the Provider

```bash
# Initialize dependencies
go mod tidy

# Build the provider
make build

# Install locally for testing
make install

# Run tests
make test

# Run acceptance tests
make acceptance
```

## Usage Examples

### Data Source Example

```terraform
data "runtime_registration_token" "current" {}

output "current_token" {
  value     = data.runtime_registration_token.current.access_token
  sensitive = true
}
```

### Resource Example

```terraform
# Basic usage
resource "runtime_registration_token" "token" {
  rotate = false
}

# Token rotation
resource "runtime_registration_token" "rotated" {
  rotate = true  # Deletes current token and generates a new one
}

output "token" {
  value     = runtime_registration_token.token.access_token
  sensitive = true
}
```

## Security Considerations

1. **Admin Privileges Required:** All operations require an Identity Token with Admin privileges
2. **Sensitive Data:** Access tokens are marked as sensitive in Terraform
3. **Token Rotation:** Use the `rotate` flag to delete the current token and generate a new one
4. **No Deletion Impact:** The resource doesn't delete the token on `terraform destroy` as it's a system-level resource

## Development Notes

- Built with Terraform Plugin Framework (v1.16.1)
- Uses JFrog shared library (v1.30.6) for common functionality
- Follows patterns established by other JFrog Terraform providers
- Compatible with Go 1.24.0+
- Supports Terraform 1.0+

## Version

Initial version: 1.0.0

## Dependencies

Key dependencies:
- github.com/hashicorp/terraform-plugin-framework v1.16.1
- github.com/jfrog/terraform-provider-shared v1.30.6
- github.com/go-resty/resty/v2 v2.16.5
- github.com/hashicorp/terraform-plugin-docs v0.20.1

## Next Steps

1. Add unit tests for data source and resource
2. Add acceptance tests
3. Set up CI/CD pipeline
4. Prepare for Terraform Registry publication
5. Add additional Runtime API endpoints as needed

## License

Apache 2.0 - Copyright (c) 2025 JFrog Ltd

