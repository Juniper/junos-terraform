# ---------------------- #
# Configuring IKE Phase1  
# ---------------------- #

resource "junos-vsrx_SecurityIkeProposalAuthentication__Method" "ikeAuth" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	authentication__method = "pre-shared-keys"
}

resource "junos-vsrx_SecurityIkeProposalDh__Group" "ikeDh" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	dh__group = "group2"
}

resource "junos-vsrx_SecurityIkeProposalAuthentication__Algorithm" "ikeAuthAlgo" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	authentication__algorithm = "sha1"
}

resource "junos-vsrx_SecurityIkeProposalEncryption__Algorithm" "ikeEncryptAlgo" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	encryption__algorithm = "aes-256-cbc"
}

resource "junos-vsrx_SecurityIkeProposalLifetime__Seconds" "ikeLifeSec" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	lifetime__seconds = "10800"
}

resource "junos-vsrx_SecurityIkePolicyProposals" "ikePolicyProp" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-policy"
	proposals = "ike-ph1-proposal"
}

resource "junos-vsrx_SecurityIkePolicyPre__Shared__KeyAscii__Text" "ikePsk" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-policy"
	ascii__text = "hello"
}

resource "junos-vsrx_SecurityIkeGatewayAddress" "ikeGwAddress" {
	resource_name = "JTAF-IKE1"
	name = "ike-gw"
	address = "10.1.1.2"
}

resource "junos-vsrx_SecurityIkeGatewayIke__Policy" "ikeGwPolicy"{
	resource_name = "JTAF-IKE1"
	name = "ike-gw"
	ike__policy = "ike-ph1-policy"
}


resource "junos-vsrx_SecurityIkeGatewayLocal__Address" "ikeGwLocalAddr" {
	resource_name = "JTAF-IKE1"
	name = "ike-gw"
	local__address = "10.1.1.1"

}

resource "junos-vsrx_SecurityIkeGatewayExternal__Interface" "ikeGwExtIntf" {
	resource_name = "JTAF-IKE1"
	name = "ike-gw"
	external__interface = "ge-0/0/0.0"
}

resource "junos-vsrx_SecurityIkeGatewayVersion" "ikeGwVer" {
	resource_name = "JTAF-IKE1"
	name = "ike-gw"
	version = "v2-only"
}
