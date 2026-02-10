terraform {
  required_providers {
    dockhand = {
      source  = "finsys/dockhand"
      version = "~> 0.1"
    }
  }
}

provider "dockhand" {
  endpoint = "http://localhost:3000"
  api_key  = "your-api-key-here"
  timeout  = 30
}

# Example: Create an environment (Docker host)
resource "dockhand_environment" "local" {
  name = "Local Docker"
  type = "local"
}

resource "dockhand_environment" "remote" {
  name = "Remote Docker Host"
  type = "ssh"
  host = "docker.example.com"
  port = 2375
  labels = {
    location = "production"
    team     = "platform"
  }
}
