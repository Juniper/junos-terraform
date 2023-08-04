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


# --------- configure IPSEC tunnel -------------- #
module "ipsec" {
    source = "./ipsec"
    providers = { junos-vsrx = junos-vsrx }
    depends_on = [ junos-vsrx_destroycommit.commit-main ]
}

# -------- commit ---------- #
resource "junos-vsrx_commit" "commit-main" {
	resource_name = "commit"
	depends_on = [module.interfaces, module.ipsec]
}

resource "junos-vsrx_destroycommit" "commit-main" {
	resource_name = "destroycommit"
}