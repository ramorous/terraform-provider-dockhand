# Testing Guide for Terraform Provider Dockhand

This guide explains how to test the Terraform Provider for Dockhand locally.

## Prerequisites

- Go 1.22 or later
- Dockhand instance running (for integration tests)
- Terraform 1.0 or later (for Terraform tests)

## Unit Tests

Run unit tests without requiring a Dockhand instance:

```bash
make test
```

This runs all tests in the `*_test.go` files across the project.

## Acceptance Tests

Acceptance tests interact with a real Dockhand API. Before running:

1. Start a Dockhand instance:
   ```bash
   docker run -p 3000:3000 dockhand:latest
   ```

2. Get an API key from the Dockhand UI (or use the default if configured)

3. Run acceptance tests:
   ```bash
   DOCKHAND_ENDPOINT="http://localhost:3000" \
   DOCKHAND_API_KEY="your-api-key" \
   make testacc
   ```

## Testing Specific Resources

Test specific resources using Go's `-run` flag:

```bash
# Test only container resource
TF_ACC=1 go test -v -run TestContainerResource ./...

# Test only compose stack resource
TF_ACC=1 go test -v -run TestComposeStackResource ./...
```

## Local Development Testing

### 1. Build the Provider

```bash
make build
```

### 2. Install Locally

```bash
make install
```

This installs the provider to `~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/`

### 3. Test with Terraform

Create a test Terraform configuration:

```hcl
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "0.1.0"
    }
  }
}

provider "dockhand" {
  endpoint = "http://localhost:3000"
  api_key  = "test-key"
}

resource "dockhand_environment" "test" {
  name = "test-env"
  type = "local"
}
```

### 4. Run Terraform

```bash
terraform plan
terraform apply
```

## Debugging

### Enable Debug Logging

```bash
TF_LOG=DEBUG terraform plan
```

### Run Provider with Debugger

To debug the provider in VS Code or delve:

```bash
go run main.go -debug
```

Then set `TF_REATTACH_PROVIDERS` environment variable and run Terraform.

## Common Issues

### "Provider not found"

Make sure the provider is installed in the correct location:

```bash
ls -la ~/.terraform.d/plugins/registry.terraform.io/ramorous/dockhand/0.1.0/linux_amd64/
```

### "Connection refused"

Ensure:
1. Dockhand instance is running
2. Endpoint URL is correct
3. API key is valid

### "Invalid API response"

Check:
1. Dockhand version compatibility
2. API endpoint path correctness
3. API key permissions

## Performance Testing

For performance testing with many resources:

```bash
# Use Terraform's parallelism setting
terraform apply -parallelism=10
```

## Continuous Integration

The project can be set up with GitHub Actions for automated testing. See `.github/workflows/` for examples.

## Reporting Test Failures

When reporting test failures:

1. Include the test output
2. Specify Go version, Terraform version, and Dockhand version
3. Attach any debug logs (with sensitive info removed)
4. List the steps to reproduce
