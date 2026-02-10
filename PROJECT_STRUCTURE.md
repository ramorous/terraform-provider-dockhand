# Project Structure

This document describes the directory structure and organization of the Terraform Provider for Dockhand.

```
terraform-provider-dockhand/
├── main.go                          # Provider entry point
├── go.mod                           # Go module definition
├── go.sum                           # Go dependencies (generated)
├── Makefile                         # Build and development commands
├── README.md                        # Main documentation
├── QUICKSTART.md                    # Quick start guide
├── CONTRIBUTING.md                 # Contribution guidelines
├── TESTING.md                       # Testing guide
├── CHANGELOG.md                     # Release history
├── LICENSE                          # License file
├── .gitignore                       # Git ignore rules
│
├── internal/
│   ├── client/                      # Dockhand API client
│   │   ├── client.go                # HTTP client and API methods
│   │   └── models.go                # Data models for API resources
│   │
│   └── provider/                    # Terraform provider implementation
│       ├── provider.go              # Main provider definition
│       ├── container_resource.go    # Container resource
│       ├── compose_stack_resource.go # Compose stack resource
│       ├── environment_resource.go  # Environment resource
│       ├── network_resource.go      # Network resource
│       ├── volume_resource.go       # Volume resource
│       ├── image_resource.go        # Image resource
│       ├── image_pull_resource.go   # Image pull resource
│       ├── containers_data_source.go # Containers data source
│       ├── compose_stacks_data_source.go # Compose stacks data source
│       ├── environments_data_source.go # Environments data source
│       ├── networks_data_source.go  # Networks data source
│       ├── volumes_data_source.go   # Volumes data source
│       └── images_data_source.go    # Images data source
│
└── examples/                        # Example Terraform configurations
    ├── terraform.tf                 # Terraform configuration
    ├── provider.tf                  # Provider configuration examples
    ├── containers.tf                # Container examples
    ├── compose_stacks.tf            # Compose stack examples
    ├── networks.tf                  # Network examples
    ├── volumes.tf                   # Volume examples
    ├── images.tf                    # Image pull examples
    ├── variables.tf                 # Input variables
    ├── outputs.tf                   # Output values
    ├── docker-compose.yaml          # Sample Docker Compose file
    └── README.md                    # Example documentation
```

## Key Components

### `main.go`
The entry point for the provider. Configures and serves the provider via the Terraform plugin protocol.

### `internal/client/`
- **client.go**: Implements the HTTP client for communicating with the Dockhand API. Includes methods for all CRUD operations on resources.
- **models.go**: Defines Go structs for all Dockhand API data types (Container, ComposeStack, Environment, etc.).

### `internal/provider/`
Each resource and data source gets its own file:
- **resource**: Manages the resource lifecycle (Create, Read, Update, Delete)
- **data_source**: Queries and returns resource information

Files are organized as:
- `{resource_name}_resource.go` for resources
- `{resource_name}_data_source.go` for data sources

### `examples/`
Complete, working Terraform configurations demonstrating how to use each resource type. Can be used as starting templates.

## Development Workflow

### Adding a New Resource

1. Create `internal/provider/{resource}_resource.go`
2. Define resource model struct
3. Implement `Resource` interface methods
4. Add to `NewResources()` in `provider.go`
5. Add tests
6. Add examples in `examples/`

### Adding a New Data Source

1. Create `internal/provider/{resource}_data_source.go`
2. Define data source model struct  
3. Implement `DataSource` interface methods
4. Add to `NewDataSources()` in `provider.go`
5. Add tests

### Adding a New API Method

1. Add method to `Client` struct in `internal/client/client.go`
2. Add corresponding model types to `internal/client/models.go`
3. Update resources/data sources to use new method
4. Add tests

## Build Output

```
terraform-provider-dockhand  # Main binary

~/.terraform.d/plugins/
└── registry.terraform.io/
    └── ramorous/
        └── dockhand/
            └── 0.1.0/
                └── linux_amd64/
                    └── terraform-provider-dockhand  # Installed provider
```

## Code Organization Principles

1. **Separation of Concerns**: Client code is separate from provider code
2. **DRY**: Common patterns are extracted to helper functions
3. **Error Handling**: Comprehensive error handling with meaningful messages
4. **Documentation**: Internal functions are documented with godoc comments
5. **Testing**: Resources and data sources should be testable
6. **Examples**: Real-world examples for each resource type

## Dependencies

See `go.mod` for complete dependency list. Key dependencies:

- `github.com/hashicorp/terraform-plugin-framework`: Terraform plugin SDK
- `github.com/go-resty/resty/v2`: HTTP client library

## Adding New Dependencies

To add a new dependency:

```bash
go get github.com/some/library
go mod tidy
```

## Running Tests

```bash
# Unit tests
make test

# Acceptance tests
make testacc

# Format code
make fmt

# Lint code
make lint
```

## Documentation

- **README.md**: Comprehensive provider documentation
- **QUICKSTART.md**: Get started in minutes  
- **TESTING.md**: Testing procedures
- **CONTRIBUTING.md**: Contribution guide
- **Code comments**: Every exported function has a godoc comment

## CI/CD

The project can be integrated with GitHub Actions for:
- Automated testing
- Code linting
- Building artifacts
- Releasing to Terraform Registry

## Versioning

Uses semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

See CHANGELOG.md for version history.
