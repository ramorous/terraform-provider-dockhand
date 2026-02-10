output "local_environment_id" {
  description = "ID of the local Docker environment"
  value       = dockhand_environment.local.id
}

output "remote_environment_id" {
  description = "ID of the remote Docker environment"
  value       = try(dockhand_environment.remote.id, "")
}

output "web_container_id" {
  description = "ID of the web container"
  value       = try(dockhand_container.web.id, "")
}

output "web_container_name" {
  description = "Name of the web container"
  value       = try(dockhand_container.web.name, "")
}

output "web_container_status" {
  description = "Status of the web container"
  value       = try(dockhand_container.web.status, "")
}

output "database_container_id" {
  description = "ID of the database container"
  value       = try(dockhand_container.database.id, "")
}

output "compose_stack_id" {
  description = "ID of the compose stack"
  value       = try(dockhand_compose_stack.app.id, "")
}

output "compose_stack_webhook_token" {
  description = "Webhook token for the compose stack (for git integration)"
  value       = try(dockhand_compose_stack.app.webhook_token, "")
  sensitive   = true
}

output "frontend_network_id" {
  description = "ID of the frontend network"
  value       = try(dockhand_network.frontend.id, "")
}

output "backend_network_id" {
  description = "ID of the backend network"
  value       = try(dockhand_network.backend.id, "")
}

output "postgres_volume_id" {
  description = "ID of the postgres data volume"
  value       = try(dockhand_volume.postgres_data.id, "")
}

output "postgres_volume_mountpoint" {
  description = "Mountpoint of the postgres data volume"
  value       = try(dockhand_volume.postgres_data.mountpoint, "")
}

output "redis_volume_id" {
  description = "ID of the redis data volume"
  value       = try(dockhand_volume.redis_data.id, "")
}

output "nginx_image_status" {
  description = "Status of the nginx image pull"
  value       = try(dockhand_image_pull.nginx.status, "")
}

output "postgres_image_status" {
  description = "Status of the postgres image pull"
  value       = try(dockhand_image_pull.postgres.status, "")
}
