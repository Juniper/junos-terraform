
terraform {
	required_providers {
		junos-vmx = {
			source = "juniper/providers/junos-vmx"
			version = "22.3"
		}
	}
}

resource "junos-vmx_SystemBackup-RouterAddress" "vmx_1" {
	resource_name = "vmx_1"
	host-name = "Spine1"
	address = "10.56.31.254"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_2" {
	resource_name = "vmx_2"
	name = "fxp0"
	name__1 = "0"
	name__2 = "10.56.16.246/19"
}

resource "junos-vmx_GroupsSystemRoot-authentication" "vmx_3"{
	resource_name = "vmx_3"
	encrypted-password = "1$ZUlES4dp$OUwWo1g7cLoV/aMWpHUnC/"
}

resource "junos-vmx_GroupsSystemLoginClass" "vmx_4"{
	resource_name = "vmx_4"
	name = "wheel"
	permissions = "snmp"
}

resource "junos-vmx_GroupsSystemLoginUserUndocumented" "vmx_5"{
	resource_name = "vmx_5"
	name = "regress"
	uid = "928"
	class = "superuser"
	shell = "csh"
}

resource "junos-vmx_GroupsSystemLoginUserAuthentication" "vmx_6"{
	resource_name = "vmx_6"
	encrypted-password = "$1$kPU..$w.4FGRAGanJ8U4Yq6sbj7."
}

resource "junos-vmx_GroupsSystemLoginUser" "vmx_7"{
	resource_name = "vmx_7"
	name = "remote"
	uid = "2000"
	class = "super-user"
}

resource "junos-vmx_GroupsSystemServicesSsh" "vmx_8"{
	resource_name = "vmx_8"
	root-login = "allow"
}

resource "junos-vmx_GroupsSystem" "vmx_9"{
	resource_name = "vmx_9"
	domain-name = "englab.juniper.net"
	domain-search = "englab.juniper.net"
	domain-search__1 = "juniper.net"
	domain-search__2 = "jnpr.net"
	time-zone = "American/Los_Angeles"
	authentication-order = "password
	authentication_order__1 = "radius"
}

resource "junos-vmx_GroupsSystemName-server" "vmx_10"{
	resource_name = "vmx_10"
	name = "10.49.32.95"
}

resource "junos-vmx_GroupsSystemName-server" "vmx_11"{
	resource_name = "vmx_11"
	name = "10.49.32.97"
}

resource "junos-vmx_GroupsSystemRadius-server" "vmx_12"{
	resource_name = "vmx_12"
	name = "10.48.144.16"
	secret = "$9$TQ/t1IcSrKAt0IRheK8X7VYgaZDm5zNdiqmTn6"
}

resource "junos-vmx_GroupsSystemRadius-server" "vmx_13"{
	resource_name = "vmx_13"
	name = "10.48.144.17"
	secret = "$9$GqUqf3nCApOPfQn690ORhSeWL7-boZDylsgoGHk"
}

resource "junos-vmx_GroupsSystemSyslogHostContents "vmx_14"{
	resource_name = "vmx_14"
	name = "log"
	name__1 = "kernel"
}

resource "junos-vmx_GroupsSystemSyslogHostContents "vmx_15"{
	resource_name = "vmx_15"
	name = "any"
}

resource "junos-vmx_GroupsSystemSyslogHostContents "vmx_16"{
	resource_name = "vmx_16"
	name = "pfe"
}

resource "junos-vmx_GroupsSystemSyslogHostContents "vmx_17"{
	resource_name = "vmx_17"
	name = "interactive-commands"
}

resource "junos-vmx_GroupsSystemSyslogFileContents "vmx_18"{
	resource_name = "vmx_18"
	name = "messages"
	name__1 = "kernel"
}

resource "junos-vmx_GroupsSystemSyslogFileContents "vmx_19"{
	resource_name = "vmx_19"
	name = "any"
}

resource "junos-vmx_GroupsSystemSyslogFileContents "vmx_20"{
	resource_name = "vmx_20"
	name = "authorization"
}

resource "junos-vmx_GroupsSystemSyslogFileContents "vmx_21"{
	resource_name = "vmx_21"
	name = "pfe"
}

resource "junos-vmx_GroupsSystemSyslogFileContents "vmx_22"{
	resource_name = "vmx_22"
	name = "security"
	name__1 = "interactive-commands"
}

resource "junos-vmx_GroupsSystemProcessDaemon-process" "vmx_23"{
	resource_name = "vmx_23"
	name = "snmp"
}

resource "junos-vmx_GroupsSystemProcessDaemon-process" "vmx_24"{
	resource_name = "vmx_24"
	name = "inet-process"
}

resource "junos-vmx_GroupsSystemProcessDaemon-process" "vmx_25"{
	resource_name = "vmx_25"
	name = "mib-process"
}

resource "junos-vmx_GroupsSystemNtpUndocumented" "vmx_26"{
	resource_name = "vmx_26"
	boot-server = "66.129.255.75"
}

resource "junos-vmx_GroupsSystemNtpServer" "vmx_27"{
	resource_name = "vmx_27"
	name = "66.129.255.75"
}

resource "junos-vmx_GroupsSystemNtpServer" "vmx_28"{
	resource_name = "vmx_28"
	name = "172.29.131.60"
}

resource "junos-vmx_GroupsSystemNtpServer" "vmx_29"{
	resource_name = "vmx_29"
	name = "172.29.135.60"
}

resource "junos-vmx_GroupsSnmpCommunity" "vmx_30"{
	resource_name = "vmx_30"
	location = "Systest lab"
	contact = "Jay Lloyd"
	interface = "fxp0.0"
	name = "public"
	authorization = "read-only"
}

resource "junos-vmx_GroupsSnmpCommunity" "vmx_31"{
	resource_name = "vmx_31"
	name = "private"
	authorization = "read-write"
}

resource "junos-vmx_GroupsRouting-optionsStaticRoute" "vmx_32"{
	resource_name = "vmx_32"
	name = "172.16.0.0/12"
	next-hop = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRoute" "vmx_33"{
	resource_name = "vmx_33"
	name = "192.168.0.0/16"
	next-hop = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRoute" "vmx_34"{
	resource_name = "vmx_34"
	name = "10.0.0.0/8"
	next-hop = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRoute" "vmx_35"{
	resource_name = "vmx_35"
	name = "66.129.0.0/16"
	next-hop = "10.56.31.254"
}

resource "junos-vmx_Groups" "vmx_36"{
	resource_name = "vmx_36"
	apply-groups = "global"
}

resource "junos-vmx_Groups" "vmx_37"{
	resource_name = "vmx_37"
	apply-groups = "member0"
}

resource "junos-vmx_SystemServicesExtension-serviceRequest-responseGrpcSsl" "vmx_38"{
	resource_name = "vmx_38"
	host-name = "spine1"
	port = "32767"
	local-certificate = "aos_grpc"
}

resource "junos-vmx_SystemServicesExtension-serviceRequest-responseGrpc" "vmx_39"{
	resource_name = "vmx_39"
	routing-instance = "mgmt_junos"
}

resource "junos-vmx_ChasisFpcPic" "vmx_40"{
	resource_name = "vmx_40"
	name = "0"
	name__1 = "0"
	number-of-ports = "96"
}

resource "junos-vmx_SecurityCertificateLocal" "vmx_41"{
	resource_name = "vmx_41"
	name = "aos_grpc"
	certificate = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDS8UeT5MDK8uMl\nVP8JjNmsMqJ6u3loR/5ct5hhfqATxVVoSXhQf/AtIBpd6GQySo74tzqjFzhje3sP\nGZ99R1/xqiPX/VYFAbq1LjSSFKVNQgc87zFEO0LoHuP/3jHWmIL5Z8LR4971vKyM\ntCEQtfEoDDzJFP+Tx/UrHGRms/jkI7f395F47C+0Ync4+6NNoE/4HUZRSimhad8U\nSdPr2phLJZTtW2JnAh7WJFcRDgBzMKPqIWikSs1ruF4RerWLVmV4t8hswD6IR5XK\noEhgZoHVFjUyBik86BClqr3c56ZeZ0b6PKo8o9W14LM5sHabd8FKxO4pC4cPVw58\nvW0R0kz1AgMBAAECggEAG+zdMPMDotpsv7B04urGlklNwTm4ZNnDDrsvbRi6UGO7\nXsb7Ko0FWrP5SZ1ORmoga0S2eojEakcaj3Ew8ADW7sM7Y4iiLp9//CRVyVD6YTPw\niVyRBRtjTB9qx6C+sE4YaLQX0nl8NsP3g9kE15V+i3KzSVuCSioahs6HbCE/PdRP\nJl0IqLReehnjtoXmffaHbyyh47B3HVKQcUQB5T3WVnijerhMAQ23FCi/bIy+f/sV\nOp9eWDeWZ34t6IxX8mrB5ANc1gYd/Pa2MqY6YicDaFGpxBDQCmPnpjZV7GIsP1F2\nby3vnimcwM4pRyWaIJPv/oWuBrZiGRfRP+BgClOLMQKBgQDuUbTpjbt7CJTK2ENQ\nRvQaQBEGWuaDPLoAjAOUbwtrBeZZp8oc1huSMOLFCxjlJ5yeeCcN1WFf0r+5RcV1\n0BvYAmKGoE0LigkRtVJNNge7ZzO5xtZo4aqnXi+LXFZSY0HIDSKtY7FQypfLVMug\ni8DqLe7b5YXQnSVRa47BYUEMDwKBgQDil56gqvQQbEtRfCwyVjFNgpRr2+G84jfI\ns9Agw4PjYGxg0rvqw3u7NfZIGam3L+nKf+MxZKF49Ir0no4q4JA7Hahw2bMrWTjg\nGppulZ2E8bBgk1uJQu2Qoz9ajRNf+cw3Q9caHXwNZ5lWrsCusbQ6EIMVFbdCjywf\nRj8ULoOiuwKBgQCzXaGIypiZX/sEtEHtcImWHDP0EYQ+r/zaHocvN2hjk1wsjEFs\n9KRpogt6/AAvsGGLT6ktGzUfdrG/0shsBEaAjEL1Sj7SeYCb6FlaLVkibekbYJDM\n/VOAhj3IaKW7emLvGxvHXin9QY1/hoF/gs+eMRX+QMA3I58LqaUW3IildwKBgHTn\nbLvP92ORu7oqqlK+DWnD+Pc81lYxED7IJpUAykbBNKkvkCwq9uc7i/je5KTSX3hO\nStTs6jIRWK+kLg0JFuCpCIJQjxPKUIGuuAZdwosYdrscy5khopeA4erB5kEaC7Zn\nlc6vb5Uq4f3K4zY+EJH0euyh6463dzgbcefjp56JAoGBAI7KjnGOyvlNEy6eUYC0\nPyCT4mfrqrND5vdHDA2dj8N/iYfL7Y+CKU7UuP/YTh9iI6FeQDLV0LwV0TVKOV+1\nPpyNf7UoBAMFs8i/vB/+S8eybXcThBr24MVJTAJJ111d+nveCv3qSc4MxFyJNvuG\nFPdib3rZ5vnSAjZjf6HNf3C1\n-----END PRIVATE KEY-----\n-----BEGIN CERTIFICATE-----\nMIIDUjCCAjqgAwIBAgICA+gwDQYJKoZIhvcNAQELBQAwVjELMAkGA1UEBhMCVVMx\nCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApNZW5sbyBQYXJrMRMwEQYDVQQKDApBcHN0\ncmEgSW5jMRAwDgYDVQQDDAdBT1NfVExTMB4XDTI0MDkyNDA4MzYzNloXDTM0MDky\nMjA4MzYzNlowVjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApN\nZW5sbyBQYXJrMRMwEQYDVQQKDApBcHN0cmEgSW5jMRAwDgYDVQQDDAdBT1NfVExT\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0vFHk+TAyvLjJVT/CYzZ\nrDKiert5aEf+XLeYYX6gE8VVaEl4UH/wLSAaXehkMkqO+Lc6oxc4Y3t7DxmffUdf\n8aoj1/1WBQG6tS40khSlTUIHPO8xRDtC6B7j/94x1piC+WfC0ePe9bysjLQhELXx\nKAw8yRT/k8f1KxxkZrP45CO39/eReOwvtGJ3OPujTaBP+B1GUUopoWnfFEnT69qY\nSyWU7VtiZwIe1iRXEQ4AczCj6iFopErNa7heEXq1i1ZleLfIbMA+iEeVyqBIYGaB\n1RY1MgYpPOgQpaq93OemXmdG+jyqPKPVteCzObB2m3fBSsTuKQuHD1cOfL1tEdJM\n9QIDAQABoyowKDAPBgNVHRMECDAGAQH/AgEAMBUGA1UdEQQOMAyCCmFwc3RyYS5j\nb20wDQYJKoZIhvcNAQELBQADggEBAKmUGAQgU97CHFDlAGLVVF9nxKmfGJwoCjAz\n+J4q9cvZE8F0S5yEHQ/iLOu6lIJy68TTEuvSQZDxGo6XGJnD86UEO6tAGKxkrvWx\nyyLoZEPdL6Ob2vPaq9jDNbXjcTCwzXvG6y4pX9384mPHN99RV9igiO1LQQzVov1v\nZ5TREif2m80JDvOxAb4aktEYr9nELPHJaPxeUHB9Nv4MGW3buwXe/LPF2hplw3FX\nDQMI123sD5gLx8U0xSwyatKJ9YbK7Ia0Iwh/EFAnIv/pu5DbH1j83yUcy4/vSsE2\nY61lFTJtIVaBUPH8p1zEvAr7qxmfwMTboe1rQ0Bkv3aPfMtDzJI=\n-----END CERTIFICATE-----\n"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_42" {
	resource_name = "vmx_42"
	name = "ge-0/0/0"
	description = "facing_mayank-rack-001-leaf3:ge-0/0/0"
	mtu = "9192"
	name__1 = "0"
	mtu__1 = "9170"
	name__2 = "10.10.10.4/31"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_43" {
	resource_name = "vmx_43"
	name = "ge-0/0/1"
	description = "facing_mayank-rack-001-leaf4:ge-0/0/0"
	mtu = "9192"
	name__1 = "0"
	mtu__1 = "9170"
	name__2 = "10.10.10.6/31"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_44" {
	resource_name = "vmx_44"
	name = "ge-0/0/2"
	description = "facing_mayank-rack-001-leaf4:ge-0/0/0"
	mtu = "9192"
	name__1 = "0"
	mtu__1 = "9170"
	name__2 = "10.10.10.6/31"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_45" {
	resource_name = "vmx_45"
	name = "ge-0/0/3"
	description = "facing_mayank-rack-001-leaf2:ge-0/0/0"
	mtu = "9192"
	name__1 = "0"
	mtu__1 = "9170"
	name__2 = "10.10.10.2/31"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_46" {
	resource_name = "vmx_46"
	name = "ge-0/0/4"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_47" {
	resource_name = "vmx_47"
	name = "ge-0/0/5"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_48" {
	resource_name = "vmx_48"
	name = "ge-0/0/6"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_49" {
	resource_name = "vmx_49"
	name = "ge-0/0/7"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_50" {
	resource_name = "vmx_50"
	name = "ge-0/0/8"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_51" {
	resource_name = "vmx_51"
	name = "ge-0/0/9"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_52" {
	resource_name = "vmx_52"
	name = "ge-0/0/10"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_53" {
	resource_name = "vmx_53"
	name = "ge-0/0/11"
	name__1 = "0"
}

resource "InterfacesInterfaceUnitFamilyInet" "vmx_54" {
	resource_name = "vmx_54"
	name = "ge-0/0/12"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_55" {
	resource_name = "vmx_55"
	name = "ge-0/0/13"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_56" {
	resource_name = "vmx_56"
	name = "ge-0/0/14"
	name__1 = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_57" {
	resource_name = "vmx_57"
	name = "lo0"
	name__1 = "0"
	name__2 = "10.10.11.0/32"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFrom" "vmx_58" {
	resource_name = "vmx_58"
	name = "AllPodNetworks"
	name__1 = "AllPodNetworks-10"
	family = "inet"
	protocols = "direct"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunity" "vmx_59" {
	resource_name = "vmx_59"
	community-name = "DEFAULT_DIRECT_V4"
}

resource "junos-vmx_Policy-optionsPolicy-statementTerm" "vmx_60" {
	resource_name = "vmx_60"
	name = "AllPodNetworks-100"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFrom" "vmx_61" {
	resource_name = "vmx_61"
	name = "BGP-AOS-Policy"
	name__1 = "BGP-AOS-Policy-10"
	policy = "AllPodNetworks"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFrom" "vmx_62" {
	resource_name = "vmx_62"
	name = "BGP-AOS-Policy-20"
	protocol = "bgp"
}

resource "junos-vmx_Policy-optionsPolicy-statementTerm" "vmx_63" {
	resource_name = "vmx_63"
	name = "BGP-AOS-Policy-100"
}

resource "junos-vmx_Policy-optionsPolicy-statement" "vmx_64" {
	resource_name = "vmx_64"
	name = "PFE-LB"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunity" "vmx_65" {
	resource_name = "vmx_65"
	name = "SPINE_TO_LEAF_EVPN_OUT"
	name__1 = "SPINE_TO_LEAF_EVPN_OUT-10"
	community-name = "FROM_SPINE_EVPN_TIER"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunity" "vmx_66" {
	resource_name = "vmx_66"
	name = "SPINE_TO_LEAF_FABRIC_OUT"
	name__1 = "SPINE_TO_LEAF_FABRIC_OUT-10"
	community-name = "FROM_SPINE_FABRIC_TIER"
}

resource "junos-vmx_Policy-optionsCommunity" "vmx_67"{
	resource_name = "vmx_67"
	name = "DEFAULT_DIRECT_V4"
	members = "1:20007"
	members__1 = "21001:26000"
}

resource "junos-vmx_Policy-optionsCommunity" "vmx_68"{
	resource_name = "vmx_68"
	name = "FROM_SPINE_EVPN_TIER"
	members = "0:14"
}

resource "junos-vmx_Policy-optionsCommunity" "vmx_69"{
	resource_name = "vmx_68"
	name = "FROM_SPINE_FABRIC_TIER"
	members = "0:15"
}

resource "junos-vmx_Routing-instancesInstanceRouting-optionsStaticRoute" "vmx_70"{
	resource_name = "vmx_70"
	name = "mgmt_junos"
	name__1 = "0.0.0.0/0"
	nexthop = "10.56.31.254"
}

resource "junos-vmx_Routing-optionsAutonomous-system" "vmx_71"{
	resource_name = "vmx_71"
	router-id = "10.10.11.0"
	as-number = "65001"
}

resource "junos-vmx_Routing-optionsForwarding-table" "vmx_72"{
	resource_name = "vmx_72"
	export = "PFE-LB"
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_73"{
	resource_name = "vmx_73"
	name = "l3clos-s"
	type = "external"
	name__1 = "10.10.10.1"
	description = "facing_mayank-rack-001-leaf1"
	local-address = "10.10.10.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_74"{
	resource_name = "vmx_74"
	name = "10.10.10.3"
	description = "facing_mayank-rack-001-leaf2"
	local-address = "10.10.10.2"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
	peer-as = "65004"
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_75"{
	resource_name = "vmx_75"
	name = "10.10.10.5"
	description = "facing_mayank-rack-001-leaf3"
	local-address = "10.10.10.4"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
	peer-as = "65005"
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_76"{
	resource_name = "vmx_76"
	name = "10.10.10.7"
	description = "facing_mayank-rack-001-leaf4"
	local-address = "10.10.10.6"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
	peer-as = "65006"
}

resource "junos-vmx_ProtocolsBgpGroupMultihop" "vmx_77"{
	resource_name = "vmx_77"
	name = "l3clos-s-evpn"
	type = "external"
	ttl = "1"
}

resource "junos-vmx_ProtocolsBgpGroupFamilyEvpnSignalingLoops" "vmx_78"{
	resource_name = "vmx_78"
	loops = "2"
}

resource "junos-vmx_ProtocolsBgpGroupBfd-liveness-detection" "vmx_79"{
	resource_name = "vmx_79"
	minimum-interval = "3000"
	multiplier = "3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_80"{
	resource_name = "vmx_80"
	name = "10.10.11.2"
	description = "facing_mayank-rack-001-leaf1-evpn-overlay"
	local-address = "10.10.11.0"
	export = "( SPINE_TO_LEAF_EVPN_OUT )" 
	peer-as = 65003
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_81"{
	resource_name = "vmx_81"
	name = "10.10.11.3"
	description = "facing_mayank-rack-001-leaf2-evpn-overlay"
	local-address = "10.10.11.0"
	export = "( SPINE_TO_LEAF_EVPN_OUT )" 
	peer-as = 65004
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_82"{
	resource_name = "vmx_82"
	name = "10.10.11.4"
	description = "facing_mayank-rack-001-leaf3-evpn-overlay"
	local-address = "10.10.11.0"
	export = "( SPINE_TO_LEAF_EVPN_OUT )" 
	peer-as = 65005
}

resource "junos-vmx_ProtocolsBgpGroupNeighbor" "vmx_83"{
	resource_name = "vmx_83"
	name = "10.10.11.5"
	description = "facing_mayank-rack-001-leaf4-evpn-overlay"
	local-address = "10.10.11.0"
	export = "( SPINE_TO_LEAF_EVPN_OUT )" 
	peer-as = 65006
}

resource "junos-vmx_ProtocolsLldpInterface" "vmx_84"{
	resource_name = "vmx_84"
	port-id-subtype = "interface-name"
	port-description-type = "interface-description"
	neighbour-port-info-display = "port-id"
	name = "all"
}