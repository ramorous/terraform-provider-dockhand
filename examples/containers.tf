# Example: Create a Docker container
resource "dockhand_container" "web" {
  environment_id = dockhand_environment.local.id
  name           = "web-server"
  image          = "nginx:latest"
  restart_policy = "unless-stopped"

  ports = [
    "80:80",
    "443:443"
  ]

  env = [
    "NGINX_HOST=example.com",
    "NGINX_PORT=80"
  ]

  labels = {
    app     = "web"
    version = "1.0"
    env     = "production"
  }

  memory = 536870912 # 512MB in bytes
  cpus   = 1.0
}

# Example: Create a container with mounts
resource "dockhand_container" "database" {
  environment_id = dockhand_environment.local.id
  name           = "postgres"
  image          = "postgres:15-alpine"
  restart_policy = "always"

  env = [
    "POSTGRES_USER=admin",
    "POSTGRES_PASSWORD=secretpassword",
    "POSTGRES_DB=myapp"
  ]

  labels = {
    app  = "database"
    tier = "backend"
  }

  memory = 1073741824 # 1GB in bytes
  cpus   = 2.0
}

# Example: Create a Redis container
resource "dockhand_container" "cache" {
  environment_id = dockhand_environment.local.id
  name           = "redis"
  image          = "redis:7-alpine"
  restart_policy = "on-failure"

  labels = {
    app  = "cache"
    tier = "data"
  }

  memory = 268435456 # 256MB in bytes
  cpus   = 0.5
}
