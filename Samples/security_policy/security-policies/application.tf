# ------------------------ #
# Configure applications
# ------------------------ #
resource "junos-vsrx_ApplicationsApplicationProtocol" "TCP-26" {
	resource_name = "JTAF-APPLICATION"
	name = "TCP-26"
	protocol = "tcp"
}

resource "junos-vsrx_ApplicationsApplicationDestination__Port" "DPORT-26" {
	resource_name = "JTAF-APPLICATION"
	name = "TCP-26"
	destination__port = "26"
}

resource "junos-vsrx_ApplicationsApplicationProtocol" "TCP-27" {
	resource_name = "JTAF-APPLICATION"
	name = "TCP-27"
	protocol = "tcp"
}

resource "junos-vsrx_ApplicationsApplicationDestination__Port" "DPORT-27" {
	resource_name = "JTAF-APPLICATION"
	name = "TCP-27"
	destination__port = "27"
}

