variable "github_token" {
  description = "GitHub personal access token for Git-based deployments"
  type        = string
  sensitive   = true
  default     = ""
}

variable "registry_username" {
  description = "Username for private container registry authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "registry_password" {
  description = "Password for private container registry authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "github_username" {
  description = "GitHub username for GitHub Container Registry authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "environment_type" {
  description = "Type of Docker environment to create"
  type        = string
  default     = "local"
  
  validation {
    condition     = contains(["local", "ssh", "docker_socket", "tcp"], var.environment_type)
    error_message = "Environment type must be one of: local, ssh, docker_socket, tcp."
  }
}

variable "remote_host" {
  description = "Host address for remote environments"
  type        = string
  default     = ""
}

variable "remote_port" {
  description = "Port for remote environments"
  type        = number
  default     = 2375
}

variable "container_image" {
  description = "Default Docker image to use for containers"
  type        = string
  default     = "alpine:latest"
}

variable "network_driver" {
  description = "Default network driver"
  type        = string
  default     = "bridge"
  
  validation {
    condition     = contains(["bridge", "overlay", "host", "null"], var.network_driver)
    error_message = "Network driver must be one of: bridge, overlay, host, null."
  }
}

variable "volume_driver" {
  description = "Default volume driver"
  type        = string
  default     = "local"
}

variable "default_restart_policy" {
  description = "Default restart policy for containers"
  type        = string
  default     = "unless-stopped"
  
  validation {
    condition     = contains(["no", "always", "on-failure", "unless-stopped"], var.default_restart_policy)
    error_message = "Restart policy must be one of: no, always, on-failure, unless-stopped."
  }
}

variable "tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    managed_by = "terraform"
    provider   = "dockhand"
  }
}
