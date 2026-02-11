# Changelog

All notable changes to the Terraform Provider for Dockhand will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [0.1.17] - 2026-02-11

### Breaking
- Replaced provider configuration attribute `api_key` with `cookie`. Users must update provider blocks and environment variables accordingly.

### Changed
- Authentication now uses a session cookie sent in the `Cookie` header rather than a bearer API key.
- Updated client to set the `Cookie` header instead of `Authorization: Bearer`.
- Updated examples, documentation, and tests to use `cookie` and `DOCKHAND_COOKIE`.

### Notes
- This is a breaking change for provider configuration: configurations using `api_key` or `DOCKHAND_API_KEY` must be migrated to `cookie` or `DOCKHAND_COOKIE`.

## [0.1.6] - 2024-02-11

### Fixed

### Data Sources
- `dockhand_environments` - Fixed some output and testing.

## [0.1.0] - 2024-02-10

### Added

#### Resources
- `dockhand_container` - Manage Docker containers
- `dockhand_compose_stack` - Manage Docker Compose stacks
- `dockhand_environment` - Manage environments for containers
- `dockhand_network` - Manage Docker networks
- `dockhand_volume` - Manage Docker volumes
- `dockhand_image` - Manage Docker images
- `dockhand_image_pull` - Pull Docker images

#### Data Sources
- `dockhand_containers` - Query Docker containers
- `dockhand_compose_stacks` - Query Docker Compose stacks
- `dockhand_environments` - Query environments
- `dockhand_networks` - Query Docker networks
- `dockhand_volumes` - Query Docker volumes
- `dockhand_images` - Query Docker images

#### Features
- Provider configuration with API authentication
- Support for environment variable configuration
- Comprehensive error handling and logging
- Full CRUD operations for all resources
- Filtering and querying capabilities for data sources

### Documentation
- Complete README with provider overview
- Quickstart guide for getting started
- Testing documentation with troubleshooting
- Project structure documentation
- Implementation summary with architecture details

---

## Version Format

Versions follow the format: `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`

- **MAJOR**: Breaking changes to the provider API
- **MINOR**: New backward-compatible features
- **PATCH**: Backward-compatible bug fixes
- **PRERELEASE**: Alpha, beta, or release candidate versions
- **BUILD**: Build metadata (rarely used)

## How to Upgrade

To upgrade to a new version of the provider:

1. Update your Terraform configuration:

```hcl
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.2"  # Update version constraint
    }
  }
}
```

2. Run `terraform init` to fetch the new provider version

3. Review the changelog for breaking changes if upgrading a MAJOR version

## Reporting Issues

If you encounter issues, please:

1. Check if a similar issue exists: https://github.com/ramorous/terraform-provider-dockhand/issues
2. Provide detailed reproduction steps
3. Include terraform version, provider version, and error logs
4. Create a new issue if needed

## Contributing Changes

To contribute:

1. Fork the repository
2. Create a feature branch
3. Make your changes and test thoroughly
4. Submit a pull request with description of changes
5. Wait for review and approval

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.
  - Provider configuration via environment variables or Terraform configuration
  - Support for multiple Docker environments (local, ssh, docker_socket, tcp)
  - Git integration for Compose stack deployments with auto-sync
  - Private registry authentication for image pulls
  - Container resource management with environment variables, labels, and resource limits
  - Docker network creation and management
  - Docker volume creation with custom drivers and options
  - Complete examples and documentation

### Planned for Future Releases

- [ ] Container logs streaming
- [ ] Container exec support
- [ ] Health check management
- [ ] Service discovery integration
- [ ] Secrets management integration
- [ ] RBAC support for Enterprise Dockhand
- [ ] Advanced debugging and monitoring
- [ ] Container build support
- [ ] Registry management

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute.

## License

This provider is licensed under the Business Source License 1.1 (BSL 1.1). See [LICENSE](LICENSE) for details.
