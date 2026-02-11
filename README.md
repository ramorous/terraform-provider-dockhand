# Terraform Provider for Dockhand

A comprehensive Terraform provider for managing Docker resources through the [Dockhand API](https://dockhand.pro/). This provider enables you to manage containers, Docker Compose stacks, networks, volumes, images, and environments directly from Terraform.

## Provider Still Under Development

Use at your own risk. Currently no maintenance or support for the time beign.

## Features

- **Container Management**: Create, update, and delete containers with full configuration support
- **Compose Stacks**: Deploy and manage Docker Compose stacks with Git integration
- **Environments**: Manage multiple Docker hosts (local and remote)
- **Networks**: Create and manage Docker networks (bridge, overlay, etc.)
- **Volumes**: Create and manage Docker volumes with custom drivers
- **Images**: Pull and manage Docker images from public and private registries
- **Data Sources**: Query existing containers, stacks, environments, networks, volumes, and images

## Requirements

- Terraform >= 1.0
- Go >= 1.22 (for building from source)
- Dockhand API access

## Building the Provider

```bash
# Clone the repository
git clone https://github.com/ramorous/terraform-provider-dockhand.git
cd terraform-provider-dockhand

# Build the provider
make build

# Install the provider locally
make install

# Run tests
make test

# Generate documentation
make doc
```

## Configuration

### Provider Configuration

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
  endpoint = "http://localhost:3000"  # Dockhand API endpoint
  cookie   = "your-session-cookie"    # Session cookie for authentication
  timeout  = 30                       # Request timeout in seconds
}
```

### Environment Variables

You can also configure the provider using environment variables:

```bash
export DOCKHAND_ENDPOINT="http://localhost:3000"
export DOCKHAND_COOKIE="your-session-cookie"
```

## Resources

### `dockhand_environment`

Manages a Docker environment (host/daemon).

```hcl
resource "dockhand_environment" "local" {
  name = "Local Docker"
  type = "local"  # local, ssh, docker_socket, tcp
}

resource "dockhand_environment" "remote" {
  name = "Remote Host"
  type = "ssh"
  host = "docker.example.com"
  port = 2375
  labels = {
    location = "production"
  }
}
```

**Supported types:**
- `local`: Local Docker daemon
- `ssh`: Remote Docker via SSH
- `docker_socket`: Remote Docker socket
- `tcp`: Remote Docker via TCP

---

### `dockhand_container`

Manages a Docker container.

```hcl
resource "dockhand_container" "web" {
  environment_id = dockhand_environment.local.id
  name           = "web-server"
  image          = "nginx:latest"
  restart_policy = "unless-stopped"

  ports = ["80:80", "443:443"]

  env = [
    "NGINX_HOST=example.com",
    "NGINX_PORT=80"
  ]

  labels = {
    app     = "web"
    version = "1.0"
  }

  memory = 536870912  # 512MB in bytes
  cpus   = 1.0
}
```

**Arguments:**
- `environment_id` - (Required) Environment ID
- `name` - (Required) Container name
- `image` - (Required) Docker image
- `restart_policy` - (Optional) Restart policy (no, always, on-failure, unless-stopped)
- `ports` - (Optional) Port mappings
- `env` - (Optional) Environment variables
- `labels` - (Optional) Container labels
- `command` - (Optional) Container command
- `args` - (Optional) Command arguments
- `memory` - (Optional) Memory limit in bytes
- `cpus` - (Optional) CPU limit

---

### `dockhand_compose_stack`

Manages a Docker Compose stack.

```hcl
resource "dockhand_compose_stack" "app" {
  environment_id = dockhand_environment.local.id
  name           = "demo-app"
  
  compose = file("${path.module}/docker-compose.yaml")

  labels = {
    project = "demo"
  }

  auto_sync = true
}

# With Git integration
resource "dockhand_compose_stack" "git_app" {
  environment_id = dockhand_environment.local.id
  name           = "git-app"
  
  compose = "version: '3.8'\nservices:\n  web:\n    image: nginx:latest"

  git_repo = {
    url    = "https://github.com/user/repo.git"
    branch = "main"
    path   = "docker"
  }

  auto_sync = true
}
```

**Arguments:**
- `environment_id` - (Required) Environment ID
- `name` - (Required) Stack name
- `compose` - (Required) Docker Compose YAML content
- `labels` - (Optional) Stack labels
- `auto_sync` - (Optional) Enable automatic sync from Git
- `git_repo` - (Optional) Git repository configuration

---

### `dockhand_network`

Manages a Docker network.

```hcl
resource "dockhand_network" "backend" {
  environment_id = dockhand_environment.local.id
  name           = "backend"
  type           = "bridge"
  driver         = "bridge"

  labels = {
    tier = "backend"
  }
}
```

**Arguments:**
- `environment_id` - (Required) Environment ID
- `name` - (Required) Network name
- `type` - (Optional) Network type (bridge, overlay, host, null)
- `driver` - (Optional) Network driver
- `labels` - (Optional) Network labels

---

### `dockhand_volume`

Manages a Docker volume.

```hcl
resource "dockhand_volume" "db_data" {
  environment_id = dockhand_environment.local.id
  name           = "postgres_data"
  driver         = "local"

  labels = {
    app = "database"
  }

  options = {
    type   = "tmpfs"
    device = "tmpfs"
    o      = "size=100m"
  }
}
```

**Arguments:**
- `environment_id` - (Required) Environment ID
- `name` - (Required) Volume name
- `driver` - (Optional) Volume driver
- `labels` - (Optional) Volume labels
- `options` - (Optional) Driver-specific options

---

### `dockhand_image`

Manages a Docker image (read-only).

```hcl
resource "dockhand_image" "nginx" {
  environment_id = dockhand_environment.local.id
  id             = "nginx:latest"
}
```

---

### `dockhand_image_pull`

Pulls a Docker image to an environment.

```hcl
resource "dockhand_image_pull" "nginx" {
  environment_id = dockhand_environment.local.id
  image          = "nginx:latest"
  registry       = "docker.io"
}

# Private registry with authentication
resource "dockhand_image_pull" "private" {
  environment_id = dockhand_environment.local.id
  image          = "myapp:latest"
  registry       = "registry.company.com"
  
  auth_username = var.registry_username
  auth_password = var.registry_password
}
```

**Arguments:**
- `environment_id` - (Required) Environment ID
- `image` - (Required) Image reference
- `registry` - (Optional) Registry URL
- `auth_username` - (Optional) Registry username
- `auth_password` - (Optional) Registry password

---

## Data Sources

### `dockhand_containers`

Query containers in an environment.

```hcl
data "dockhand_containers" "all" {
  environment_id = dockhand_environment.local.id
}

output "containers" {
  value = data.dockhand_containers.all.containers
}
```

---

### `dockhand_compose_stacks`

Query compose stacks in an environment.

```hcl
data "dockhand_compose_stacks" "all" {
  environment_id = dockhand_environment.local.id
}
```

---

### `dockhand_environments`

Query all environments.

```hcl
data "dockhand_environments" "all" {}

output "environments" {
  value = data.dockhand_environments.all.environments
}
```

---

### `dockhand_networks`

Query networks in an environment.

```hcl
data "dockhand_networks" "all" {
  environment_id = dockhand_environment.local.id
}
```

---

### `dockhand_volumes`

Query volumes in an environment.

```hcl
data "dockhand_volumes" "all" {
  environment_id = dockhand_environment.local.id
}
```

---

### `dockhand_images`

Query images in an environment.

```hcl
data "dockhand_images" "all" {
  environment_id = dockhand_environment.local.id
}
```

---

## Complete Example

See the `examples/` directory for complete examples including:

- `provider.tf` - Provider configuration
- `containers.tf` - Container examples
- `compose_stacks.tf` - Compose stack examples
- `networks.tf` - Network examples
- `volumes.tf` - Volume examples
- `images.tf` - Image pull examples

## Development

### Running Tests

```bash
# Run all tests
make test

# Run acceptance tests (requires running Dockhand instance)
make testacc

# Format code
make fmt

# Lint code
make lint
```

### Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This provider is licensed under the Dockhand BSL 1.1 License. See [LICENSE](LICENSE) for details.

For questions or support, visit [Dockhand Documentation](https://dockhand.pro/manual).

## Roadmap

- [ ] Advanced compose stack features
- [ ] Container logs streaming
- [ ] Container exec support
- [ ] Health checks
- [ ] Service discovery
- [ ] Secrets management
- [ ] More data sources and filters

## Support

For issues, feature requests, or questions:
- Open an issue on GitHub
- Check the [Dockhand Documentation](https://dockhand.pro/manual)
- Visit the [Dockhand GitHub Repository](https://github.com/ramorous/dockhand)
