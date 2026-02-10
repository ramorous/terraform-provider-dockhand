# Terraform Provider Dockhand - Development Setup

Complete setup instructions for local development and testing of the Terraform Provider for Dockhand.

## System Requirements

- **Go**: 1.22 or higher
- **Terraform**: 1.0 or higher
- **Git**: Latest version
- **Linux/macOS/Windows**: All major platforms supported

## Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/ramorous/terraform-provider-dockhand.git
cd terraform-provider-dockhand
```

### 2. Install Dependencies

```bash
# Download Go module dependencies
go mod download

# Verify dependencies
go mod verify
```

### 3. Build Provider

```bash
# Quick build
make build

# Or native go command
go build -o terraform-provider-dockhand
```

You should see a binary file created: `terraform-provider-dockhand`

## Development Environment

### Configure Provider Locally

```bash
# Install to local Terraform directory
make install

# Verify installation
ls -la ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/
```

### Create Test Configuration

Create `terraform.tf` in your project:

```hcl
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }
}

provider "dockhand" {
  api_url = "http://localhost:3000"
  api_key = "your-api-key"
}

# Your resources here
```

## Testing

### Run All Tests

```bash
# Unit tests
make test

# With coverage
go test -v -cover ./...

# With race detector (recommended)
go test -v -race ./...
```

### Run Specific Tests

```bash
# Test a specific package
go test -v ./internal/provider

# Test a specific test
go test -v ./internal/provider -run TestProviderConfigure
```

### Acceptance Tests

Requires a running Dockhand API instance:

```bash
# Set environment variables
export DOCKHAND_API_URL="http://localhost:3000"
export DOCKHAND_API_KEY="your-api-key"

# Run acceptance tests
TF_ACC=1 go test -v ./...

# Run specific acceptance test
TF_ACC=1 go test -v ./internal/provider -run TestAccContainer
```

## Code Quality

### Format Code

```bash
make fmt

# Or use go command directly
go fmt ./...
```

### Lint

```bash
make lint

# Requires golangci-lint
# Install: https://golangci-lint.run/usage/install/
```

### Generate Documentation

```bash
make doc

# Generates documentation from code comments
# Output goes to /docs directory
```

## Common Development Tasks

### Add a New Resource

1. Create file: `internal/provider/new_resource_resource.go`
2. Implement Resource interface
3. Register in `provider.go`
4. Add tests in `new_resource_resource_test.go`
5. Update README.md

### Add a New Data Source

1. Create file: `internal/provider/new_data_source_data_source.go`
2. Implement DataSource interface
3. Register in `provider.go`
4. Add tests in `new_data_source_data_source_test.go`
5. Update README.md

### Debug Provider

Set debug mode:

```bash
# Run provider in debug mode
go run . -debug

# In another terminal, set environment variable before init
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/ramorous/dockhand":{"protocol":"grpc","pid":12345,"test":true,"addr":{"unix":"/var/folders/...}}}'

# Then use terraform
terraform init
terraform apply
```

## Troubleshooting

### Build Fails

```bash
# Clean build
rm terraform-provider-dockhand
go clean -testcache
make build
```

### Tests Fail

```bash
# Check Go version
go version

# Update dependencies
go mod tidy
go mod download

# Verify module
go mod verify
```

### Provider Not Loading

```bash
# Check installation
ls -la ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/

# Verify binary is executable
file terraform-provider-dockhand

# Check permissions
chmod +x terraform-provider-dockhand

# Reinstall
make uninstall
make install
```

### API Connection Issues

1. Verify Dockhand API is running
2. Check API URL and credentials
3. Test connectivity: `curl http://localhost:3000/api/health`
4. Check provider logs for detailed errors

## Development Workflow

### Before Committing

```bash
# Format code
make fmt

# Run linter
make lint

# Run all tests
make test

# Build
make build

# Check for uncommitted changes
git status
```

### Create Feature Branch

```bash
# Update main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/your-feature-name

# Make changes
# Test thoroughly
# Commit with conventional commits

git add .
git commit -m "feat: description of changes"

# Push branch
git push origin feature/your-feature-name

# Create pull request on GitHub
```

### After Merge

```bash
# Sync local main
git checkout main
git pull origin main

# Delete local feature branch
git branch -d feature/your-feature-name

# Delete remote feature branch (or use GitHub UI)
git push origin --delete feature/your-feature-name
```

## Editor Setup

### VS Code

Add to `.vscode/settings.json`:

```json
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.lintArgs": [
    "--fast"
  ],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "go.useLanguageServer": true
}
```

### Recommended Extensions

- Go (golang.go)
- HashiCorp Terraform (hashicorp.terraform)
- gopls language server

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Terraform Registry Publishing](https://registry.terraform.io/docs/providers/publish-provider)
- [Contributing Guide](./CONTRIBUTING.md)
- [Release Process](./RELEASING.md)

## Getting Help

- Review [CONTRIBUTING.md](./CONTRIBUTING.md) for development guidelines
- Check [TESTING.md](./TESTING.md) for testing procedures
- See [RELEASING.md](./RELEASING.md) for release information
- Open an issue: https://github.com/ramorous/terraform-provider-dockhand/issues
