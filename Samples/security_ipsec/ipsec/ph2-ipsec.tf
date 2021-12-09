# --------------------------- #
#  IPSEC Phase2 configuration
#  -------------------------- #

resource "junos-vsrx_SecurityIpsecProposalProtocol" "ipsecPropProto" {
	resource_name ="JTAF_IPSEC3"
	name = "ipsec-ph2"
	protocol = "esp"
}

resource "junos-vsrx_SecurityIpsecProposalAuthentication__Algorithm" "ipsecPropAuthAlgo" {
	resource_name = "JTAF_IPSEC4"
	name = "ipsec-ph2"
	authentication__algorithm = "hmac-sha1-96"
}

resource "junos-vsrx_SecurityIpsecProposalEncryption__Algorithm" "ipsecPropEncryptAlgo" {
	resource_name = "JTAF_IPSEC5"
	name = "ipsec-ph2"
	encryption__algorithm = "aes-256-cbc"
}

resource "junos-vsrx_SecurityIpsecProposalLifetime__Seconds" "ipsecPropLifeSec" {
	resource_name = "JTAF_IPSEC6"
	name = "ipsec-ph2"
	lifetime__seconds = "2000"
}

resource "junos-vsrx_SecurityIpsecPolicyPerfect__Forward__SecrecyKeys" "ipsecPolPfs" {
	resource_name = "JTAF_IPSEC7"
	name = "ipsec-policy"
	keys = "group2" 
}	

resource "junos-vsrx_SecurityIpsecPolicyProposals" "ipsecPolProp" {
	resource_name = "JTAF_IPSEC8"
	name = "ipsec-policy"
	proposals = "ipsec-ph2"
}

resource "junos-vsrx_SecurityIpsecVpnBind__Interface" "ipsecBindIntf" {
	resource_name = "JTAF_IPSEC9"
	name = "ipsec-1"
	bind__interface = "st0.0"
}

resource "junos-vsrx_SecurityIpsecVpnIkeGateway" "ipsecIkeGw" {
	resource_name = "JTAF_IPSEC10"
	name = "ipsec-1"
	gateway = "ike-gw"
	depends_on = [
		junos-vsrx_SecurityIkeGatewayAddress.ikeGwAddress,
        junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficProtocolsName.seczonesprotocol_st0,
        junos-vsrx_SecurityZonesSecurity__ZoneInterfacesHost__Inbound__TrafficSystem__ServicesName.seczonesservices_st0  	
	]
}

resource "junos-vsrx_SecurityIpsecVpnIkeIpsec__Policy" "ipsecPol" {
	resource_name = "JTAF_IPSEC11"
	name = "ipsec-1"
	ipsec__policy = "ipsec-policy"
	
}

resource "junos-vsrx_SecurityIpsecVpnEstablish__Tunnels" "estdTunnel" {
	resource_name = "JTAF_IPSEC12"
	name = "ipsec-1"
	establish__tunnels = "immediately"
	
}
