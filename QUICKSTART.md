# Quick Start Guide

Get up and running with the Terraform Provider for Dockhand in minutes!

## Prerequisites

- Terraform >= 1.0
- A running Dockhand instance
- Dockhand session cookie

## Installation

### Option 1: Using Terraform Registry (Recommended for released versions)

Add to your `terraform.tf`:

```hcl
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }
}
```

Then run:

```bash
terraform init
```

### Option 2: Building from Source

```bash
git clone https://github.com/ramorous/terraform-provider-dockhand.git
cd terraform-provider-dockhand
make build
make install
```

## Configure the Provider

Create a `provider.tf` file:

```hcl
terraform {
  required_providers {
    dockhand = {
      source = "ramorous/dockhand"
    }
  }
}

provider "dockhand" {
  endpoint = "http://localhost:3000"  # Your Dockhand endpoint
  cookie   = var.dockhand_cookie
  timeout  = 30
}
```

Create `terraform.tfvars`:

```hcl
dockhand_cookie = "your-session-cookie-here"
```

Or use environment variables:

```bash
export DOCKHAND_ENDPOINT="http://localhost:3000"
export DOCKHAND_COOKIE="your-session-cookie"
```

## First Configuration

Create a `main.tf` file with your first resources:

### 1. Create an Environment

```hcl
resource "dockhand_environment" "local" {
  name = "My Local Docker"
  type = "local"
}
```

### 2. Create a Network

```hcl
resource "dockhand_network" "app" {
  environment_id = dockhand_environment.local.id
  name           = "app-network"
  type           = "bridge"
}
```

### 3. Create a Volume

```hcl
resource "dockhand_volume" "database" {
  environment_id = dockhand_environment.local.id
  name           = "database_data"
  driver         = "local"

  labels = {
    app = "database"
  }
}
```

### 4. Pull an Image

```hcl
resource "dockhand_image_pull" "postgres" {
  environment_id = dockhand_environment.local.id
  image          = "postgres:15-alpine"
}
```

### 5. Create a Container

```hcl
resource "dockhand_container" "database" {
  depends_on     = [dockhand_image_pull.postgres]
  environment_id = dockhand_environment.local.id
  
  name   = "my-database"
  image  = "postgres:15-alpine"
  restart_policy = "unless-stopped"

  env = [
    "POSTGRES_USER=dbuser",
    "POSTGRES_PASSWORD=dbpassword",
    "POSTGRES_DB=myapp"
  ]

  labels = {
    app = "database"
    env = "development"
  }

  memory = 1073741824  # 1GB
  cpus   = 1.0
}
```

## Apply Configuration

```bash
# Review planned changes
terraform plan

# Apply the configuration
terraform apply
```

## Verify Resources

Check what Terraform created:

```bash
terraform state list
terraform show
```

## Deploy a Compose Stack

For a more complex application with multiple services:

```hcl
resource "dockhand_compose_stack" "myapp" {
  environment_id = dockhand_environment.local.id
  name           = "myapp-stack"
  
  compose = file("${path.module}/docker-compose.yaml")

  labels = {
    project = "myapp"
  }

  auto_sync = true
}
```

Where `docker-compose.yaml` contains your stack definition.

## Query Existing Resources

Use data sources to query resources:

```hcl
data "dockhand_containers" "all" {
  environment_id = dockhand_environment.local.id
}

output "containers" {
  value = data.dockhand_containers.all.containers
}
```

Run Terraform to see the output:

```bash
terraform apply -target=data.dockhand_containers.all
terraform output containers
```

## Manage State

Terraform stores state about your infrastructure. By default, it's stored locally in `terraform.tfstate`.

For team collaboration, use remote state:

```hcl
terraform {
  backend "remote" {
    organization = "my-org"
    workspaces {
      name = "dockhand"
    }
  }
}
```

Initialize remote state:

```bash
terraform init
```

## Destroy Resources

To remove all managed resources:

```bash
terraform destroy
```

Alternatively, remove specific resources:

```bash
terraform destroy -target=dockhand_container.database
```

## Troubleshooting

### "Provider could not be found"

Make sure Terraform has initialized:

```bash
terraform init
```

### "Could not connect to Dockhand"

Check:
1. Dockhand is running
2. Endpoint URL is correct
3. Network connectivity

### "Session cookie is invalid"

Verify:
1. Session cookie is correct
2. Session cookie has not expired
3. Session cookie has necessary permissions

## Next Steps

- [Review detailed documentation](README.md)
- [Check resource examples](examples/)
- [Learn about data sources](README.md#data-sources)
- [Set up remote state](https://www.terraform.io/language/state/remote)
- [Use Terraform workspaces](https://www.terraform.io/language/state/workspaces)

## Getting Help

- [GitHub Issues](https://github.com/ramorous/terraform-provider-dockhand/issues)
- [Dockhand Documentation](https://dockhand.pro/manual)
- [Terraform Documentation](https://www.terraform.io/docs)

## Tips & Best Practices

1. **Use Variables:** Replace hardcoded values with variables
   ```hcl
   variable "cookie" {
     sensitive = true
   }
   ```

2. **Use Outputs:** Export important values
   ```hcl
   output "container_id" {
     value = dockhand_container.myapp.id
   }
   ```

3. **Use Modules:** Organize complex configurations
   ```hcl
   module "database" {
     source = "./modules/database"
     environment_id = dockhand_environment.local.id
   }
   ```

4. **Use Labels:** Tag your resources for organization
   ```hcl
   labels = {
     project = "myapp"
     env     = "production"
     managed = "terraform"
   }
   ```

5. **Use Dependencies:** Ensure correct creation order
   ```hcl
   depends_on = [dockhand_image_pull.nginx]
   ```

Happy Terraforming! 🚀
