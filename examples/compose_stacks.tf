# Example: Create a Docker Compose Stack
resource "dockhand_compose_stack" "app" {
  environment_id = dockhand_environment.local.id
  name           = "demo-app"
  
  compose = <<-EOT
    version: '3.8'
    services:
      web:
        image: nginx:latest
        ports:
          - "80:80"
          - "443:443"
        environment:
          - NGINX_HOST=example.com
        restart_policy:
          condition: unless-stopped
        labels:
          com.example.app: "web"
      
      api:
        image: node:18-alpine
        ports:
          - "3000:3000"
        environment:
          - NODE_ENV=production
          - PORT=3000
        restart_policy:
          condition: unless-stopped
        labels:
          com.example.app: "api"
      
      db:
        image: postgres:15-alpine
        environment:
          - POSTGRES_USER=user
          - POSTGRES_PASSWORD=password
          - POSTGRES_DB=myapp
        volumes:
          - postgres_data:/var/lib/postgresql/data
        restart_policy:
          condition: unless-stopped
        labels:
          com.example.app: "database"
    
    volumes:
      postgres_data:
        driver: local
  EOT

  labels = {
    project = "demo"
    env     = "production"
  }

  auto_sync = true
}

# Example: Compose stack with Git integration
resource "dockhand_compose_stack" "git_app" {
  environment_id = dockhand_environment.local.id
  name           = "git-deployed-app"
  
  compose = file("${path.module}/docker-compose.yaml")

  git_repo = {
    url    = "https://github.com/username/repo.git"
    branch = "main"
    path   = "docker"
    auth_type  = "https"
    auth_token = var.github_token # Use variables for sensitive data
  }

  auto_sync = true

  labels = {
    deployment = "git-managed"
    source     = "github"
  }
}
