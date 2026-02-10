# Example: Create Docker networks
resource "dockhand_network" "frontend" {
  environment_id = dockhand_environment.local.id
  name           = "frontend"
  type           = "bridge"
  driver         = "bridge"

  labels = {
    tier = "frontend"
    app  = "web"
  }
}

resource "dockhand_network" "backend" {
  environment_id = dockhand_environment.local.id
  name           = "backend"
  type           = "bridge"
  driver         = "bridge"

  labels = {
    tier = "backend"
    app  = "api"
  }
}

resource "dockhand_network" "database" {
  environment_id = dockhand_environment.local.id
  name           = "database"
  type           = "bridge"
  driver         = "bridge"

  labels = {
    tier = "database"
    app  = "db"
  }
}

# Example: Create overlay network for Swarm (if using Swarm mode)
resource "dockhand_network" "overlay" {
  environment_id = dockhand_environment.local.id
  name           = "swarm-overlay"
  type           = "overlay"
  driver         = "overlay"
  scope          = "swarm"

  labels = {
    swarm = "true"
  }
}
