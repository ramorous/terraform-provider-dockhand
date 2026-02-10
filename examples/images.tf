# Example: Pull Docker images
resource "dockhand_image_pull" "nginx" {
  environment_id = dockhand_environment.local.id
  image          = "nginx:latest"
  registry       = "docker.io"
}

resource "dockhand_image_pull" "postgres" {
  environment_id = dockhand_environment.local.id
  image          = "postgres:15-alpine"
  registry       = "docker.io"
}

resource "dockhand_image_pull" "redis" {
  environment_id = dockhand_environment.local.id
  image          = "redis:7-alpine"
  registry       = "docker.io"
}

# Example: Pull image with specific tag
resource "dockhand_image_pull" "node" {
  environment_id = dockhand_environment.local.id
  image          = "node:18-alpine"
  registry       = "docker.io"
}

# Example: Pull from private registry
resource "dockhand_image_pull" "private_app" {
  environment_id = dockhand_environment.local.id
  image          = "myapp:latest"
  registry       = "registry.company.com"
  
  auth_username = var.registry_username
  auth_password = var.registry_password
}

# Example: Pull from GitHub Container Registry
resource "dockhand_image_pull" "ghcr_app" {
  environment_id = dockhand_environment.local.id
  image          = "myapp:latest"
  registry       = "ghcr.io"
  
  auth_username = var.github_username
  auth_password = var.github_token
}

# Example: Reference pulled image for later use
resource "dockhand_container" "web_from_pulled" {
  depends_on     = [dockhand_image_pull.nginx]
  environment_id = dockhand_environment.local.id
  name           = "web-from-pulled"
  image          = dockhand_image_pull.nginx.image
  
  labels = {
    deployment = "standard"
  }
}
