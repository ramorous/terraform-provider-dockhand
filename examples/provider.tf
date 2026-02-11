terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }
}

provider "dockhand" {
  endpoint = "http://localhost:3000"
  cookie   = "your-session-cookie-here"
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

# Example: Read all environments
data "dockhand_environments" "all" {}

output "environment_ids" {
  value = [for e in data.dockhand_environments.all.environments : e.id]
}
