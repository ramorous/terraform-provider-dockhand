# Terraform Dockhand Provider - Common Patterns and Use Cases

This guide demonstrates common patterns and best practices for using the Dockhand provider.

## Table of Contents

- [Basic Setup](#basic-setup)
- [Environment Management](#environment-management)
- [Container Deployment](#container-deployment)
- [Multi-Tier Applications](#multi-tier-applications)
- [Compose Stack Management](#compose-stack-management)
- [Monitoring and Inventory](#monitoring-and-inventory)
- [Best Practices](#best-practices)

## Basic Setup

### Provider Configuration

```terraform
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }
}

provider "dockhand" {
  api_url = var.dockhand_api_url
  cookie  = var.dockhand_cookie
}
```

### Environment Variables

Create `terraform.tfvars`:

```hcl
dockhand_api_url = "http://localhost:3000"
dockhand_cookie  = "your-session-cookie-here"
```

Or set environment variables:

```bash
export DOCKHAND_API_URL="http://localhost:3000"
export DOCKHAND_COOKIE="your-session-cookie-here"
```

## Environment Management

### Create Multiple Environments

```terraform
locals {
  environments = {
    dev = {
      name     = "development"
      endpoint = "unix:///var/run/docker.sock"
    }
    staging = {
      name     = "staging"
      endpoint = "tcp://staging-docker:2375"
    }
    prod = {
      name     = "production"
      endpoint = "tcp://prod-docker:2376"
    }
  }
}

resource "dockhand_environment" "all" {
  for_each = local.environments
  
  name = each.value.name
  
  docker_config = {
    endpoint = each.value.endpoint
    tls      = contains(["prod", "staging"], each.key)
  }
  
  labels = {
    environment = each.key
    managed_by  = "terraform"
  }
}

output "environment_ids" {
  value = {
    for k, v in dockhand_environment.all : k => v.id
  }
}
```

### Query Environment Status

```terraform
data "dockhand_environments" "connected" {
  filter = {
    status = "connected"
  }
}

output "healthy_environments" {
  value = data.dockhand_environments.connected.environments[*].name
}
```

## Container Deployment

### Simple Container

```terraform
resource "dockhand_container" "simple" {
  environment_id = dockhand_environment.dev.id
  name           = "nginx"
  image          = "nginx:latest"
  
  ports = [
    {
      container_port = 80
      host_port      = 8080
    }
  ]
  
  restart_policy = "unless-stopped"
}
```

### Container with Full Configuration

```terraform
resource "dockhand_container" "web" {
  environment_id = dockhand_environment.prod.id
  name           = "web-app"
  image          = var.web_image
  
  ports = [
    {
      container_port = 3000
      host_port      = 3000
    },
    {
      container_port = 443
      host_port      = 443
    }
  ]
  
  environment = [
    {
      key   = "NODE_ENV"
      value = "production"
    },
    {
      key   = "DATABASE_URL"
      value = "postgres://user:pass@db:5432/app"
    },
    {
      key   = "LOG_LEVEL"
      value = "info"
    }
  ]
  
  volumes = [
    {
      container_path = "/app/data"
      volume_name    = dockhand_volume.app_data.name
    },
    {
      container_path = "/app/logs"
      volume_name    = dockhand_volume.app_logs.name
    }
  ]
  
  resources = {
    cpu_limit    = "1000m"
    memory_limit = "512m"
  }
  
  networks = [dockhand_network.backend.name]
  
  labels = {
    app         = "web"
    environment = "production"
    team        = "platform"
  }
  
  restart_policy = "unless-stopped"
}
```

### Deploy Multiple Containers from Variables

```terraform
variable "services" {
  type = map(object({
    image  = string
    ports  = list(number)
    memory = string
  }))
}

resource "dockhand_container" "services" {
  for_each = var.services
  
  environment_id = dockhand_environment.prod.id
  name           = each.key
  image          = each.value.image
  
  ports = [
    for port in each.value.ports : {
      container_port = port
      host_port      = port
    }
  ]
  
  resources = {
    memory_limit = each.value.memory
  }
  
  restart_policy = "unless-stopped"
}
```

## Multi-Tier Applications

### Complete Application Stack

```terraform
# Network tier
resource "dockhand_network" "app_network" {
  environment_id = dockhand_environment.prod.id
  name           = "app-stack"
  driver         = "bridge"
  
  ipam_config = {
    subnet = "10.0.0.0/24"
  }
}

# Storage tier
resource "dockhand_volume" "db_data" {
  environment_id = dockhand_environment.prod.id
  name           = "postgres-data"
  
  labels = {
    tier    = "database"
    backup  = "daily"
  }
}

# Database tier
resource "dockhand_container" "db" {
  environment_id = dockhand_environment.prod.id
  name           = "postgres"
  image          = "postgres:15"
  
  environment = [
    {
      key   = "POSTGRES_PASSWORD"
      value = var.db_password
    },
    {
      key   = "POSTGRES_DB"
      value = "appdb"
    }
  ]
  
  volumes = [
    {
      container_path = "/var/lib/postgresql/data"
      volume_name    = dockhand_volume.db_data.name
    }
  ]
  
  networks       = [dockhand_network.app_network.name]
  restart_policy = "unless-stopped"
}

# Cache tier
resource "dockhand_container" "cache" {
  environment_id = dockhand_environment.prod.id
  name           = "redis"
  image          = "redis:7-alpine"
  
  ports = [
    {
      container_port = 6379
      host_port      = 6379
    }
  ]
  
  networks       = [dockhand_network.app_network.name]
  restart_policy = "unless-stopped"
}

# Application tier
resource "dockhand_container" "app" {
  environment_id = dockhand_environment.prod.id
  name           = "app"
  image          = var.app_image
  
  ports = [
    {
      container_port = 3000
      host_port      = 3000
    }
  ]
  
  environment = [
    {
      key   = "DATABASE_URL"
      value = "postgres://postgres:${var.db_password}@postgres:5432/appdb"
    },
    {
      key   = "REDIS_URL"
      value = "redis://redis:6379"
    }
  ]
  
  networks       = [dockhand_network.app_network.name]
  restart_policy = "unless-stopped"
  
  depends_on = [
    dockhand_container.db,
    dockhand_container.cache
  ]
}
```

## Compose Stack Management

### Deploy from Local File

```terraform
resource "dockhand_compose_stack" "local" {
  environment_id = dockhand_environment.dev.id
  name           = "dev-stack"
  compose_file   = file("${path.module}/docker-compose.yml")
}
```

### Deploy from Git with Auto-Update

```terraform
resource "dockhand_compose_stack" "git_repo" {
  environment_id = dockhand_environment.prod.id
  name           = "production"
  
  git_config = {
    repository = "https://github.com/example/docker-compose.git"
    ref        = "main"
    path       = "docker-compose.yml"
  }
  
  environment = {
    ENVIRONMENT = "production"
    LOG_LEVEL   = "info"
  }
  
  auto_update         = true
  poll_interval_seconds = 600
}
```

## Monitoring and Inventory

### Environment Health Check

```terraform
data "dockhand_environments" "all" {}

local {
  environment_status = {
    for env in data.dockhand_environments.all.environments :
    env.name => {
      status      = env.status
      docker_ver  = env.docker_version
      containers  = env.container_count
      memory_gb   = env.total_memory / 1024 / 1024 / 1024
      cpus        = env.total_cpus
    }
  }
}

output "environment_health" {
  value = local.environment_status
}
```

### Container Inventory

```terraform
data "dockhand_containers" "all" {
  environment_id = dockhand_environment.prod.id
}

output "running_services" {
  value = {
    for c in data.dockhand_containers.all.containers :
    c.name => {
      image      = c.image
      status     = c.state
      ports      = c.ports
      created_at = c.created_at
    }
  }
}
```

### Cleanup Unused Resources

```terraform
# Find unused volumes
data "dockhand_volumes" "all" {
  environment_id = dockhand_environment.prod.id
}

locals {
  unused_volumes = [
    for v in data.dockhand_volumes.all.volumes :
    v.name
    if !v.in_use
  ]
}

# Find unused images
data "dockhand_images" "all" {
  environment_id = dockhand_environment.prod.id
}

locals {
  unused_images = [
    for i in data.dockhand_images.all.images :
    "${i.repository}:${i.tag}"
    if !i.in_use
  ]
}

output "cleanup_candidates" {
  value = {
    volumes = local.unused_volumes
    images  = local.unused_images
  }
}
```

## Best Practices

### 1. Use Variables for Configuration

```terraform
variable "environment" {
  type        = string
  description = "Environment name (dev, staging, prod)"
}

variable "app_version" {
  type        = string
  description = "Application version to deploy"
}

# Use in resources
resource "dockhand_container" "app" {
  image = "myapp:${var.app_version}"
  
  labels = {
    deployed_at = timestamp()
    environment = var.environment
  }
}
```

### 2. Use Locals for Computed Values

```terraform
locals {
  is_production = var.environment == "prod"
  
  container_config = {
    memory_limit = local.is_production ? "1g" : "512m"
    cpu_limit    = local.is_production ? "2000m" : "1000m"
    replicas     = local.is_production ? 3 : 1
  }
}
```

### 3. Organize Resources by Module

```
├── main.tf              # Environments
├── networks.tf         # Networks and storage
├── containers.tf       # Containers
├── stacks.tf          # Compose stacks
├── variables.tf       # Variables
├── outputs.tf         # Outputs
└── terraform.tfvars   # Values (not in git)
```

### 4. Use Data Sources for Existing Resources

```terraform
# Query existing resources instead of recreating
data "dockhand_networks" "existing" {
  environment_id = dockhand_environment.prod.id
  
  filter = {
    name = "existing-network"
  }
}

# Reference without managing
resource "dockhand_container" "app" {
  networks = [data.dockhand_networks.existing.networks[0].name]
}
```

### 5. Use Depends-On for Complex Dependencies

```terraform
resource "dockhand_container" "app" {
  # Depends on storage and network being ready
  depends_on = [
    dockhand_volume.app_data,
    dockhand_network.app_network,
    dockhand_container.db  # Wait for DB to start
  ]
}
```

### 6. Document with Labels

```terraform
resource "dockhand_container" "web" {
  labels = {
    "app"             = "web-server"
    "environment"     = var.environment
    "managed_by"      = "terraform"
    "owner"           = "platform-team"
    "backup_policy"   = "daily"
    "monitoring"      = "true"
    "cost_center"     = "engineering"
  }
}
```

### 7. Use Workspaces for Multiple Environments

```bash
# Create workspaces
terraform workspace new dev
terraform workspace new staging
terraform workspace new prod

# Deploy to each
terraform workspace select dev
terraform apply -var-file=dev.tfvars

terraform workspace select staging
terraform apply -var-file=staging.tfvars

terraform workspace select prod
terraform apply -var-file=prod.tfvars
```

### 8. Implement State Locking

```terraform
terraform {
  backend "consul" {
    path = "terraform/dockhand"
  }
}
```

### 9. Use Import for Existing Resources

```bash
# Import existing containers
terraform import dockhand_container.existing env-id:container-name

# Import existing networks
terraform import dockhand_network.existing env-id:network-name
```

### 10. Validate Before Applying

```bash
# Format code
terraform fmt

# Validate syntax
terraform validate

# Check plan
terraform plan -out=tfplan

# Review the plan
cat tfplan

# Apply safe
terraform apply tfplan
```

## Additional Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Dockhand API](https://github.com/ramorous/dockhand)
- [Docker Documentation](https://docs.docker.com)
- Provider Resources: `/docs/resources/`
- Data Sources: `/docs/data-sources/`
