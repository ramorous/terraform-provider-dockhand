# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-10

### Added

- Initial release of Terraform Provider for Dockhand
- **Resources:**
  - `dockhand_environment` - Manage Docker environments (hosts)
  - `dockhand_container` - Manage Docker containers
  - `dockhand_compose_stack` - Manage Docker Compose stacks with optional Git integration
  - `dockhand_network` - Manage Docker networks
  - `dockhand_volume` - Manage Docker volumes
  - `dockhand_image` - Reference Docker images (read-only)
  - `dockhand_image_pull` - Pull Docker images from registries
  
- **Data Sources:**
  - `dockhand_containers` - Query containers in an environment
  - `dockhand_compose_stacks` - Query compose stacks in an environment
  - `dockhand_environments` - Query all environments
  - `dockhand_networks` - Query networks in an environment
  - `dockhand_volumes` - Query volumes in an environment
  - `dockhand_images` - Query images in an environment

- **Features:**
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
