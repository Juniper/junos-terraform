variable "HOST" {
	description = "Host name"
	type = string
	}

variable "PORT" {
	description = "Port number to make netconf connection to"
	type = number
	default = 830
}

variable "USERNAME" {}
variable "PASSWORD" {}
variable "SSHKEY" {}

terraform {
	required_providers {
		junos-vsrx = {
			source = "juniper/providers/junos-vsrx"
			version = "1.1"
		}
	}
}

provider "junos-vsrx" {
	host = var.HOST
	port = var.PORT
	username = var.USERNAME
	password = var.PASSWORD
	sshkey = var.SSHKEY
}

# --------- configure security policies -------- #
# policies include the below
# 1. application
# 2. Address book
# 3. Address set
# 4. security policy
# ---------------------------------------------- #
module "security-policies" {
	source = "./security-policies"
	providers = { junos-vsrx = junos-vsrx }
	depends_on = [ junos-vsrx_destroycommit.commit-main ]
}


# -------- commit ---------- #
resource "junos-vsrx_commit" "commit-main" {
	resource_name = "commit"
	depends_on = [module.security-policies]
}

resource "junos-vsrx_destroycommit" "commit-main" {
	resource_name = "destroycommit"
}