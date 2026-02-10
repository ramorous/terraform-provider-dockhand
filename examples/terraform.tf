terraform {
  required_version = ">= 1.0"

  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }

  # Uncomment the following to use remote state storage
  # backend "remote" {
  #   organization = "your-org"
  #   workspaces {
  #     name = "dockhand-terraform"
  #   }
  # }

  # Or use local state with the following:
  # backend "local" {
  #   path = "terraform.tfstate"
  # }
}
