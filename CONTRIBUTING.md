# Contributing to Terraform Provider for Dockhand

Thank you for your interest in contributing to the Terraform Provider for Dockhand! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions. Treat all contributors with courtesy.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Terraform 1.0 or higher (for testing)
- Git

### Setup Development Environment

```bash
git clone https://github.com/ramorous/terraform-provider-dockhand.git
cd terraform-provider-dockhand
go mod download
```

### Building

```bash
make build
# or
go build -o terraform-provider-dockhand
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests (requires Dockhand API running)
make testacc

# Linting
make lint

# Code formatting
make fmt
```

## Development Workflow

### 1. Fork and Clone

Fork the repository on GitHub and clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/terraform-provider-dockhand.git
cd terraform-provider-dockhand
```

### 2. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-new-resource`: For new features
- `fix/bug-description`: For bug fixes
- `docs/update-readme`: For documentation changes

### 3. Make Changes

Edit files as needed. Keep commits focused and atomic:

```bash
git add .
git commit -m "feat: add support for resource X"
```

### 4. Test Your Changes

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Build
make build
```

### 5. Update Documentation

If adding new resources or data sources:
- Update README.md with new resource/data source
- Add inline documentation in the resource/data source code
- Update examples in the `examples/` directory

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Create a pull request on GitHub with:
- Clear title and description
- Reference to any related issues
- Summary of changes
- Testing performed

## Adding New Resources

### Resource Template

1. Create `internal/provider/{resource_name}_resource.go`

```go
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ramorous/terraform-provider-dockhand/internal/client"
)

type NewResourceResource struct {
	client *client.Client
}

func (r *NewResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_resource"
}

func (r *NewResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a new resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *NewResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Implementation
}

func (r *NewResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Implementation
}

func (r *NewResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Implementation
}

func (r *NewResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Implementation
}
```

2. Register the resource in `internal/provider/provider.go`
3. Add tests
4. Update documentation

## Adding New Data Sources

Similar to resources, create data sources in `internal/provider/{data_source_name}_data_source.go` and register them in the provider.

## Adding New API Methods

Update `internal/client/client.go` to add new API methods that your resources need.

## Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Add comments for exported functions

## Commit Guidelines

Use conventional commits:

```
type: subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Tests
- `chore`: Build, dependencies, etc.
- `refactor`: Code refactoring

Example:
```
feat: add image_pull resource

Add support for pulling Docker images through the Dockhand API.
Includes full CRUD operations and resource testing.

Closes #123
```

## Review Process

- All pull requests require review
- CI/CD checks must pass
- Documentation should be updated
- Tests should be included for new functionality

## Reporting Issues

When reporting bugs, include:
- Terraform version
- Provider version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages/logs

## Documentation

Documentation is generated from code comments. For resources and data sources:

```go
// NewResourceResource manages a new resource
type NewResourceResource struct {
	client *client.Client
}
```

Public functions should have comments explaining their purpose.

## Testing

### Unit Tests

Place unit tests in `*_test.go` files:

```go
func TestNewResource_Create(t *testing.T) {
	// Test implementation
}
```

### Acceptance Tests

For integration tests:

```go
func TestAccNewResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test steps
		},
	})
}
```

Run with:
```bash
TF_ACC=1 go test -v ./...
```

## Questions?

- Check existing issues: https://github.com/ramorous/terraform-provider-dockhand/issues
- Review documentation: https://registry.terraform.io/providers/ramorous/dockhand
- Reference Terraform Plugin Framework: https://developer.hashicorp.com/terraform/plugin/framework

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## Submitting Changes

1. Commit your changes with clear messages: `git commit -am "Add feature X"`
2. Push to your fork: `git push origin feature/my-feature`
3. Create a Pull Request with:
   - Clear description of changes
   - Reference to any related issues
   - Updated documentation
   - Tests where applicable

## PR Review Process

- Maintainers will review PRs within a reasonable timeframe
- Please respond to feedback and questions
- Once approved, your PR will be merged

## Reporting Issues

- Use GitHub Issues to report bugs
- Include:
  - Terraform version
  - Provider version
  - Dockhand version
  - Steps to reproduce
  - Expected vs actual behavior
  - Any error messages

## Feature Requests

Feature requests are welcome! Please:

1. Check if the feature is already requested
2. Describe the use case and expected behavior
3. Provide examples of how it would be used
4. Explain why it would be valuable

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (BSL 1.1).
