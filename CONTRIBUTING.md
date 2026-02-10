# Contributing to Terraform Provider Dockhand

Thank you for your interest in contributing to the Terraform Provider for Dockhand! This document outlines the process for contributing changes.

## Code of Conduct

Please be respectful and constructive in all interactions with other contributors and maintainers.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/terraform-provider-dockhand.git`
3. Create a feature branch: `git checkout -b feature/my-feature`
4. Set up your development environment: `make build`

## Development Process

### Building

```bash
# Build the provider
make build

# Install locally for testing
make install

# Run tests
make test

# Run acceptance tests (requires Dockhand instance running)
make testacc
```

### Code Style

- Follow standard Go conventions
- Run `make fmt` to format code
- Run `make lint` to check for issues

### Writing Tests

- Add unit tests for new functions
- Add acceptance tests (if applicable) in `*_test.go` files
- Use table-driven tests when appropriate

### Documentation

- Update README.md with new resources/data sources
- Add examples to the `examples/` directory
- Include resource/data source documentation in code comments

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
