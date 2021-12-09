# ---------------------- #
# Configuring IKE Phase1  
# ---------------------- #

resource "junos-vsrx_SecurityIkeProposalAuthentication__Method" "ikeAuth" {
	resource_name = "JTAF-IKE1"
	name = "ike-ph1-proposal"
	authentication__method = "pre-shared-keys"
}

resource "junos-vsrx_SecurityIkeProposalDh__Group" "ikeDh" {
	resource_name = "JTAF-IKE2"
	name = "ike-ph1-proposal"
	dh__group = "group2"
}

resource "junos-vsrx_SecurityIkeProposalAuthentication__Algorithm" "ikeAuthAlgo" {
	resource_name = "JTAF-IKE3"
	name = "ike-ph1-proposal"
	authentication__algorithm = "sha1"
}

resource "junos-vsrx_SecurityIkeProposalEncryption__Algorithm" "ikeEncryptAlgo" {
	resource_name = "JTAF-IKE4"
	name = "ike-ph1-proposal"
	encryption__algorithm = "aes-256-cbc"
}

resource "junos-vsrx_SecurityIkeProposalLifetime__Seconds" "ikeLifeSec" {
	resource_name = "JTAF-IKE5"
	name = "ike-ph1-proposal"
	lifetime__seconds = "10800"
}

resource "junos-vsrx_SecurityIkePolicyProposals" "ikePolicyProp" {
	resource_name = "JTAF-IKE6"
	name = "ike-ph1-policy"
	proposals = "ike-ph1-proposal"
}

resource "junos-vsrx_SecurityIkePolicyPre__Shared__KeyAscii__Text" "ikePsk" {
	resource_name = "JTAF-IKE7"
	name = "ike-ph1-policy"
	ascii__text = "hello"
	
	// This helps to ensure Terraform doesn't get upset when
	// data at rest is encrypted
	// TODO: Comment this in the docs
	// show | display inheritance no-comments
	lifecycle {
    ignore_changes = [
      ascii__text
    ]
  }
}

resource "junos-vsrx_SecurityIkeGatewayAddress" "ikeGwAddress" {
	resource_name = "JTAF-IKE8"
	name = "ike-gw"
	address = "101.101.101.13"
}

resource "junos-vsrx_SecurityIkeGatewayIke__Policy" "ikeGwPolicy"{
	resource_name = "JTAF-IKE9"
	name = "ike-gw"
	ike__policy = "ike-ph1-policy"
}


resource "junos-vsrx_SecurityIkeGatewayLocal__Address" "ikeGwLocalAddr" {
	resource_name = "JTAF-IKE10"
	name = "ike-gw"
	local__address = "101.101.101.3"

}

resource "junos-vsrx_SecurityIkeGatewayExternal__Interface" "ikeGwExtIntf" {
	resource_name = "JTAF-IKE11"
	name = "ike-gw"
	external__interface = "ge-0/0/0.0"
}

resource "junos-vsrx_SecurityIkeGatewayVersion" "ikeGwVer" {
	resource_name = "JTAF-IKE12"
	name = "ike-gw"
	version = "v2-only"
}
