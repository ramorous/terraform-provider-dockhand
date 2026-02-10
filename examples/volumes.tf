# Example: Create Docker volumes
resource "dockhand_volume" "postgres_data" {
  environment_id = dockhand_environment.local.id
  name           = "postgres_data"
  driver         = "local"

  labels = {
    app      = "database"
    data     = "persistent"
    backup   = "daily"
  }
}

resource "dockhand_volume" "redis_data" {
  environment_id = dockhand_environment.local.id
  name           = "redis_data"
  driver         = "local"

  labels = {
    app  = "cache"
    type = "kvstore"
  }
}

resource "dockhand_volume" "app_logs" {
  environment_id = dockhand_environment.local.id
  name           = "app_logs"
  driver         = "local"

  options = {
    type   = "tmpfs"
    device = "tmpfs"
    o      = "size=100m"
  }

  labels = {
    app      = "application"
    type     = "logs"
    retention = "7days"
  }
}

# Example: Named volume for shared data
resource "dockhand_volume" "shared_data" {
  environment_id = dockhand_environment.local.id
  name           = "shared_data"
  driver         = "local"

  labels = {
    purpose = "shared"
    access  = "multi-container"
  }
}

# Example: NFS volume for distributed storage
resource "dockhand_volume" "nfs_storage" {
  environment_id = dockhand_environment.local.id
  name           = "nfs_mount"
  driver         = "local"

  options = {
    type   = "nfs"
    o      = "addr=192.168.1.100,vers=4,soft,timeo=180,bg,tcp,rw"
    device = ":/export/data"
  }

  labels = {
    storage = "nfs"
    network = "mounted"
  }
}
