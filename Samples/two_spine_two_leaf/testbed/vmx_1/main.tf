
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

resource "junos-vmx_GroupsSystemRoot-authenticationEncrypted-password" "vmx_3"{
	resource_name = "vmx_3"
	encrypted-password = "1$ZUlES4dp$OUwWo1g7cLoV/aMWpHUnC/"
}

resource "junos-vmx_GroupsSystemLoginClassName" "vmx_4"{
	resource_name = "vmx_4"
	name = "wheel"
}

resource "junos-vmx_GroupsSystemLoginClassPermissions" "vmx_5"{
	resource_name = "vmx_5"
	permissions = "snmp"
}


resource "junos-vmx_GroupsSystemLoginUserName" "vmx_6"{
	resource_name = "vmx_6"
	name = "regress"
}

resource "junos-vmx_GroupsSystemLoginUserUid" "vmx_7"{
	resource_name = "vmx_7"
	uid = "928"
}

resource "junos-vmx_GroupsSystemLoginUserClass" "vmx_8"{
	resource_name = "vmx_8"
	class = "superuser"
}

resource "junos-vmx_GroupsSystemLoginUserUndocumentedShell" "vmx_9"{
	resource_name = "vmx_9"
	shell = "csh"
}

resource "junos-vmx_GroupsSystemLoginUserAuthenticationEncrypted-password" "vmx_10"{
	resource_name = "vmx_10"
	encrypted-password = "$1$kPU..$w.4FGRAGanJ8U4Yq6sbj7."
}

resource "junos-vmx_GroupsSystemLoginUserName" "vmx_11"{
	resource_name = "vmx_11"
	name = "remote"
}

resource "junos-vmx_GroupsSystemLoginUserUid" "vmx_12"{
	resource_name = "vmx_12"
	uid = "2000"
	class = "super-user"
}

resource "junos-vmx_GroupsSystemLoginUserClass" "vmx_13"{
	resource_name = "vmx_13"
	class = "super-user"
}


resource "junos-vmx_GroupsSystemServicesSshRoot-login" "vmx_14"{
	resource_name = "vmx_14"
	root-login = "allow"
}

resource "junos-vmx_GroupsSystemDomain-name" "vmx_15"{
	resource_name = "vmx_15"
	domain-name = "englab.juniper.net"
}

resource "junos-vmx_GroupsSystemDomain-search" "vmx_16"{
	resource_name = "vmx_16"
	domain-search = "englab.juniper.net"
}

resource "junos-vmx_GroupsSystemDomain-search" "vmx_17"{
	resource_name = "vmx_17"
	domain-search = "juniper.net"
}

resource "junos-vmx_GroupsSystemDomain-search" "vmx_18"{
	resource_name = "vmx_18"
	domain-search = "jnpr.net"
}

resource "junos-vmx_GroupsSystemTime-zone" "vmx_19"{
	resource_name = "vmx_19"
	time-zone = "America/Los_Angeles"
}

resource "junos-vmx_GroupsSystemAuthentication-order" "vmx_20"{
	resource_name = "vmx_20"
	authentication-order = "password" 
}

resource "junos-vmx_GroupsSystemAuthentication-order" "vmx_21"{
	resource_name = "vmx_21"
	authentication-order = "radius" 
}

resource "junos-vmx_GroupsSystemName-serverName" "vmx_22"{
	resource_name = "vmx_22"
	name = "10.49.32.95"
}

resource "junos-vmx_GroupsSystemName-serverName" "vmx_23"{
	resource_name = "vmx_23"
	name = "10.49.32.97"
}

resource "junos-vmx_GroupsSystemRadius-serverName" "vmx_24"{
	resource_name = "vmx_24"
	name = "10.48.144.16"
}

resource "junos-vmx_GroupsSystemRadius-serverSecret" "vmx_25"{
	resource_name = "vmx_25"
	secret = "$9$TQ/t1IcSrKAt0IRheK8X7VYgaZDm5zNdiqmTn6"
}

resource "junos-vmx_GroupsSystemRadius-serverName" "vmx_26"{
	resource_name = "vmx_26"
	name = "10.48.144.17"
}

resource "junos-vmx_GroupsSystemRadius-serverSecret" "vmx_27"{
	resource_name = "vmx_27"
	secret = "$9$GqUqf3nCApOPfQn690ORhSeWL7-boZDylsgoGHk"
}

resource "junos-vmx_GroupsSystemSyslogHostName" "vmx_28"{
	resource_name = "vmx_28"
	name = "log"
}

resource "junos-vmx_GroupsSystemSyslogHostContentsName" "vmx_29"{
	resource_name = "vmx_29"
	name = "kernel"
}


resource "junos-vmx_GroupsSystemSyslogHostContentsName" "vmx_30"{
	resource_name = "vmx_30"
	name = "any"
}

resource "junos-vmx_GroupsSystemSyslogHostContentsName" "vmx_31"{
	resource_name = "vmx_31"
	name = "pfe"
}

resource "junos-vmx_GroupsSystemSyslogHostContentsName" "vmx_32"{
	resource_name = "vmx_32"
	name = "interactive-commands"
}

resource "junos-vmx_GroupsSystemSyslogFileName" "vmx_33"{
	resource_name = "vmx_33"
	name = "messages"
}

resource "junos-vmx_GroupsSystemSyslogFileContentsName" "vmx_34"{
	resource_name = "vmx_34"
	name = "kernel"
}

resource "junos-vmx_GroupsSystemSyslogFileContentsName" "vmx_35"{
	resource_name = "vmx_35"
	name = "any"
}

resource "junos-vmx_GroupsSystemSyslogFileContentsName" "vmx_36"{
	resource_name = "vmx_36"
	name = "authorization"
}

resource "junos-vmx_GroupsSystemSyslogFileContentsName" "vmx_37"{
	resource_name = "vmx_37"
	name = "pfe"
}

resource "junos-vmx_GroupsSystemSyslogFileName" "vmx_38"{
	resource_name = "vmx_38"
	name = "security"
}

resource "junos-vmx_GroupsSystemSyslogFileContentsName" "vmx_39"{
	resource_name = "vmx_39"
	name = "interactive-commands"
}

resource "junos-vmx_GroupsSystemProcessDaemon-processName" "vmx_40"{
	resource_name = "vmx_40"
	name = "snmp"
}

resource "junos-vmx_GroupsSystemProcessDaemon-processName" "vmx_41"{
	resource_name = "vmx_41"
	name = "inet-process"
}

resource "junos-vmx_GroupsSystemProcessDaemon-processName" "vmx_42"{
	resource_name = "vmx_42"
	name = "mib-process"
}

resource "junos-vmx_GroupsSystemNtpUndocumentedBoot-server" "vmx_43"{
	resource_name = "vmx_43"
	boot-server = "66.129.255.75"
}

resource "junos-vmx_GroupsSystemNtpServerName" "vmx_44"{
	resource_name = "vmx_44"
	name = "66.129.255.75"
}

resource "junos-vmx_GroupsSystemNtpServerNAme" "vmx_45"{
	resource_name = "vmx_45"
	name = "172.29.131.60"
}

resource "junos-vmx_GroupsSystemNtpServerName" "vmx_46"{
	resource_name = "vmx_46"
	name = "172.29.135.60"
}

resource "junos-vmx_GroupsSnmpLocation" "vmx_47"{
	resource_name = "vmx_47"
	location = "Systest lab"
}

resource "junos-vmx_GroupsSnmpContact" "vmx_48"{
	resource_name = "vmx_48"
	contact = "Jay Lloyd"
}

resource "junos-vmx_GroupsSnmpInterface" "vmx_49"{
	resource_name = "vmx_49"
	interface = "fxp0.0"
}

resource "junos-vmx_GroupsSnmpCommunityName" "vmx_50"{
	resource_name = "vmx_50"
	name = "public"
}

resource "junos-vmx_GroupsSnmpCommunityAuthorization" "vmx_51"{
	resource_name = "vmx_51"
	authorization = "read-only"
}

resource "junos-vmx_GroupsSnmpCommunityName" "vmx_52"{
	resource_name = "vmx_52"
	name = "private"
}

resource "junos-vmx_GroupsSnmpCommunityAuthorization" "vmx_53"{
	resource_name = "vmx_53"
	authorization = "read-write"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteName" "vmx_54"{
	resource_name = "vmx_54"
	name = "172.16.0.0/12"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteNext-hop" "vmx_55"{
	resource_name = "vmx_55"
	name = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteName" "vmx_56"{
	resource_name = "vmx_56"
	name = "192.168.0.0/16"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteNext-hop" "vmx_57"{
	resource_name = "vmx_57"
	name = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteName" "vmx_58"{
	resource_name = "vmx_58"
	name = "10.0.0.0/8"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteNext-hop" "vmx_59"{
	resource_name = "vmx_59"
	name = "10.56.31.254"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteName" "vmx_60"{
	resource_name = "vmx_60"
	name = "66.129.0.0/16"
}

resource "junos-vmx_GroupsRouting-optionsStaticRouteNext-hop" "vmx_61"{
	resource_name = "vmx_61"
	name = "10.56.31.254"
}

resource "junos-vmx_Apply-groups" "vmx_62"{
	resource_name = "vmx_62"
	apply-groups = "global"
}

resource "junos-vmx_Apply-groups" "vmx_63"{
	resource_name = "vmx_63"
	apply-groups = "member0"
}

resource "junos-vmx_SystemHost-name" "vmx_64"{
	resource_name = "vmx_64"
	host-name = "spine1"
}

resource "junos-vmx_SystemServicesExtension-serviceRequest-responseGrpcSslPort" "vmx_65"{
	resource_name = "vmx_65"
	port = "32767"
}

resource "junos-vmx_SystemServicesExtension-serviceRequest-responseGrpcSslLocal-certificate" "vmx_66"{
	resource_name = "vmx_66"
	local-certificate = "aos_grpc"
}

resource "junos-vmx_SystemServicesExtension-serviceRequest-responseGrpcRouting-instance" "vmx_67"{
	resource_name = "vmx_67"
	routing-instance = "mgmt_junos"
}

resource "junos-vmx_ChasisFpcName" "vmx_68"{
	resource_name = "vmx_68"
	name = "0"
}

resource "junos-vmx_ChasisFpcPicName" "vmx_69"{
	resource_name = "vmx_69"
	name = "0"
}

resource "junos-vmx_ChasisFpcPicNumber-of-ports" "vmx_70"{
	resource_name = "vmx_70"
	number-of-ports = "96"
}

resource "junos-vmx_SecurityCertificateLocalName" "vmx_71"{
	resource_name = "vmx_71"
	name = "aos_grpc"
}

resource "junos-vmx_SecurityCertificateLocalCertificate" "vmx_72"{
	resource_name = "vmx_72"
	certificate = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDS8UeT5MDK8uMl\nVP8JjNmsMqJ6u3loR/5ct5hhfqATxVVoSXhQf/AtIBpd6GQySo74tzqjFzhje3sP\nGZ99R1/xqiPX/VYFAbq1LjSSFKVNQgc87zFEO0LoHuP/3jHWmIL5Z8LR4971vKyM\ntCEQtfEoDDzJFP+Tx/UrHGRms/jkI7f395F47C+0Ync4+6NNoE/4HUZRSimhad8U\nSdPr2phLJZTtW2JnAh7WJFcRDgBzMKPqIWikSs1ruF4RerWLVmV4t8hswD6IR5XK\noEhgZoHVFjUyBik86BClqr3c56ZeZ0b6PKo8o9W14LM5sHabd8FKxO4pC4cPVw58\nvW0R0kz1AgMBAAECggEAG+zdMPMDotpsv7B04urGlklNwTm4ZNnDDrsvbRi6UGO7\nXsb7Ko0FWrP5SZ1ORmoga0S2eojEakcaj3Ew8ADW7sM7Y4iiLp9//CRVyVD6YTPw\niVyRBRtjTB9qx6C+sE4YaLQX0nl8NsP3g9kE15V+i3KzSVuCSioahs6HbCE/PdRP\nJl0IqLReehnjtoXmffaHbyyh47B3HVKQcUQB5T3WVnijerhMAQ23FCi/bIy+f/sV\nOp9eWDeWZ34t6IxX8mrB5ANc1gYd/Pa2MqY6YicDaFGpxBDQCmPnpjZV7GIsP1F2\nby3vnimcwM4pRyWaIJPv/oWuBrZiGRfRP+BgClOLMQKBgQDuUbTpjbt7CJTK2ENQ\nRvQaQBEGWuaDPLoAjAOUbwtrBeZZp8oc1huSMOLFCxjlJ5yeeCcN1WFf0r+5RcV1\n0BvYAmKGoE0LigkRtVJNNge7ZzO5xtZo4aqnXi+LXFZSY0HIDSKtY7FQypfLVMug\ni8DqLe7b5YXQnSVRa47BYUEMDwKBgQDil56gqvQQbEtRfCwyVjFNgpRr2+G84jfI\ns9Agw4PjYGxg0rvqw3u7NfZIGam3L+nKf+MxZKF49Ir0no4q4JA7Hahw2bMrWTjg\nGppulZ2E8bBgk1uJQu2Qoz9ajRNf+cw3Q9caHXwNZ5lWrsCusbQ6EIMVFbdCjywf\nRj8ULoOiuwKBgQCzXaGIypiZX/sEtEHtcImWHDP0EYQ+r/zaHocvN2hjk1wsjEFs\n9KRpogt6/AAvsGGLT6ktGzUfdrG/0shsBEaAjEL1Sj7SeYCb6FlaLVkibekbYJDM\n/VOAhj3IaKW7emLvGxvHXin9QY1/hoF/gs+eMRX+QMA3I58LqaUW3IildwKBgHTn\nbLvP92ORu7oqqlK+DWnD+Pc81lYxED7IJpUAykbBNKkvkCwq9uc7i/je5KTSX3hO\nStTs6jIRWK+kLg0JFuCpCIJQjxPKUIGuuAZdwosYdrscy5khopeA4erB5kEaC7Zn\nlc6vb5Uq4f3K4zY+EJH0euyh6463dzgbcefjp56JAoGBAI7KjnGOyvlNEy6eUYC0\nPyCT4mfrqrND5vdHDA2dj8N/iYfL7Y+CKU7UuP/YTh9iI6FeQDLV0LwV0TVKOV+1\nPpyNf7UoBAMFs8i/vB/+S8eybXcThBr24MVJTAJJ111d+nveCv3qSc4MxFyJNvuG\nFPdib3rZ5vnSAjZjf6HNf3C1\n-----END PRIVATE KEY-----\n-----BEGIN CERTIFICATE-----\nMIIDUjCCAjqgAwIBAgICA+gwDQYJKoZIhvcNAQELBQAwVjELMAkGA1UEBhMCVVMx\nCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApNZW5sbyBQYXJrMRMwEQYDVQQKDApBcHN0\ncmEgSW5jMRAwDgYDVQQDDAdBT1NfVExTMB4XDTI0MDkyNDA4MzYzNloXDTM0MDky\nMjA4MzYzNlowVjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApN\nZW5sbyBQYXJrMRMwEQYDVQQKDApBcHN0cmEgSW5jMRAwDgYDVQQDDAdBT1NfVExT\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0vFHk+TAyvLjJVT/CYzZ\nrDKiert5aEf+XLeYYX6gE8VVaEl4UH/wLSAaXehkMkqO+Lc6oxc4Y3t7DxmffUdf\n8aoj1/1WBQG6tS40khSlTUIHPO8xRDtC6B7j/94x1piC+WfC0ePe9bysjLQhELXx\nKAw8yRT/k8f1KxxkZrP45CO39/eReOwvtGJ3OPujTaBP+B1GUUopoWnfFEnT69qY\nSyWU7VtiZwIe1iRXEQ4AczCj6iFopErNa7heEXq1i1ZleLfIbMA+iEeVyqBIYGaB\n1RY1MgYpPOgQpaq93OemXmdG+jyqPKPVteCzObB2m3fBSsTuKQuHD1cOfL1tEdJM\n9QIDAQABoyowKDAPBgNVHRMECDAGAQH/AgEAMBUGA1UdEQQOMAyCCmFwc3RyYS5j\nb20wDQYJKoZIhvcNAQELBQADggEBAKmUGAQgU97CHFDlAGLVVF9nxKmfGJwoCjAz\n+J4q9cvZE8F0S5yEHQ/iLOu6lIJy68TTEuvSQZDxGo6XGJnD86UEO6tAGKxkrvWx\nyyLoZEPdL6Ob2vPaq9jDNbXjcTCwzXvG6y4pX9384mPHN99RV9igiO1LQQzVov1v\nZ5TREif2m80JDvOxAb4aktEYr9nELPHJaPxeUHB9Nv4MGW3buwXe/LPF2hplw3FX\nDQMI123sD5gLx8U0xSwyatKJ9YbK7Ia0Iwh/EFAnIv/pu5DbH1j83yUcy4/vSsE2\nY61lFTJtIVaBUPH8p1zEvAr7qxmfwMTboe1rQ0Bkv3aPfMtDzJI=\n-----END CERTIFICATE-----\n"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_73" {
	resource_name = "vmx_73"
	name = "ge-0/0/0"
}

resource "junos-vmx_InterfacesInterfaceDescription" "vmx_74" {
	resource_name = "vmx_74"
	description = "facing_mayank-rack-001-leaf3:ge-0/0/0"
}

resource "junos-vmx_InterfacesInterfaceMtu" "vmx_75" {
	resource_name = "vmx_75"
	mtu = "9192"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_76" {
	resource_name = "vmx_76"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetMtu" "vmx_77" {
	resource_name = "vmx_77"
	mtu = "9170"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_78" {
	resource_name = "vmx_78"
	name = "10.10.10.4/31"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_79" {
	resource_name = "vmx_79"
	name = "ge-0/0/1"
}

resource "junos-vmx_InterfacesInterfaceDescription" "vmx_80" {
	resource_name = "vmx_80"
	description = "facing_mayank-rack-001-leaf4:ge-0/0/0"
}

resource "junos-vmx_InterfacesInterfaceMtu" "vmx_81" {
	resource_name = "vmx_81"
	mtu = "9192"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_82" {
	resource_name = "vmx_82"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetMtu" "vmx_83" {
	resource_name = "vmx_83"
	mtu = "9170"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_84" {
	resource_name = "vmx_84"
	name = "10.10.10.4/31"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_85" {
	resource_name = "vmx_85"
	name = "ge-0/0/2"
}

resource "junos-vmx_InterfacesInterfaceDescription" "vmx_86" {
	resource_name = "vmx_86"
	description = "facing_mayank-rack-001-leaf1:ge-0/0/0"
}

resource "junos-vmx_InterfacesInterfaceMtu" "vmx_87" {
	resource_name = "vmx_87"
	mtu = "9192"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_88" {
	resource_name = "vmx_88"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetMtu" "vmx_89" {
	resource_name = "vmx_89"
	mtu = "9170"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_90" {
	resource_name = "vmx_90"
	name = "10.10.10.4/31"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_91" {
	resource_name = "vmx_91"
	name = "ge-0/0/3"
}

resource "junos-vmx_InterfacesInterfaceDescription" "vmx_92" {
	resource_name = "vmx_92"
	description = "facing_mayank-rack-001-leaf2:ge-0/0/0"
}

resource "junos-vmx_InterfacesInterfaceMtu" "vmx_93" {
	resource_name = "vmx_93"
	mtu = "9192"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_94" {
	resource_name = "vmx_94"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetMtu" "vmx_95" {
	resource_name = "vmx_95"
	mtu = "9170"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_96" {
	resource_name = "vmx_96"
	name = "10.10.10.4/31"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_97" {
	resource_name = "vmx_97"
	name = "ge-0/0/4"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_98" {
	resource_name = "vmx_98"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_99" {
	resource_name = "vmx_99"
	name = "ge-0/0/5"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_100" {
	resource_name = "vmx_100"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_101" {
	resource_name = "vmx_101"
	name = "ge-0/0/6"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_102" {
	resource_name = "vmx_102"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_103" {
	resource_name = "vmx_103"
	name = "ge-0/0/7"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_104" {
	resource_name = "vmx_104"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_105" {
	resource_name = "vmx_105"
	name = "ge-0/0/8"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_106" {
	resource_name = "vmx_106"
	name = "0"
}


resource "junos-vmx_InterfacesInterfaceName" "vmx_107" {
	resource_name = "vmx_107"
	name = "ge-0/0/9"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_108" {
	resource_name = "vmx_108"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_109" {
	resource_name = "vmx_109"
	name = "ge-0/0/10"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_110" {
	resource_name = "vmx_110"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_111" {
	resource_name = "vmx_111"
	name = "ge-0/0/11"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_112" {
	resource_name = "vmx_112"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_113" {
	resource_name = "vmx_113"
	name = "ge-0/0/12"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_114" {
	resource_name = "vmx_114"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_115" {
	resource_name = "vmx_115"
	name = "ge-0/0/13"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_116" {
	resource_name = "vmx_116"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_117" {
	resource_name = "vmx_117"
	name = "ge-0/0/14"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_118" {
	resource_name = "vmx_118"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_119" {
	resource_name = "vmx_119"
	name = "lo0"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_120" {
	resource_name = "vmx_120"
	name = "0"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_121" {
	resource_name = "vmx_121"
	name = "10.10.11.0/32"
}

resource "junos-vmx_Policy-optionsPolicy-statementName" "vmx_122" {
	resource_name = "vmx_122"
	name = "AllPodNetworks"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_123" {
	resource_name = "vmx_123"
	name = "AllPodNetworks-10"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFromFamily" "vmx_124" {
	resource_name = "vmx_124"
	family = "inet"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFromProtocol" "vmx_125" {
	resource_name = "vmx_125"
	protocol = "direct"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunityCommunity-name" "vmx_126" {
	resource_name = "vmx_126"
	community-name = "DEFAULT_DIRECT_V4"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_127" {
	resource_name = "vmx_127"
	name = "AllPodNetworks-100"
}

resource "junos-vmx_Policy-optionsPolicy-statementName" "vmx_128" {
	resource_name = "vmx_128"
	name = "BGP-AOS-Policy"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_129" {
	resource_name = "vmx_129"
	name = "BGP-AOS-Policy-10"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFromPolicy" "vmx_130" {
	resource_name = "vmx_130"
	policy = "AllPodNetworks"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_131" {
	resource_name = "vmx_131"
	name = "BGP-AOS-Policy-20"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermFromProtocol" "vmx_132" {
	resource_name = "vmx_132"
	protocol = "bgp"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_133" {
	resource_name = "vmx_133"
	name = "BGP-AOS-Policy-100"
}

resource "junos-vmx_Policy-optionsPolicy-statementName" "vmx_134" {
	resource_name = "vmx_134"
	name = "PFE-LB"
}

resource "junos-vmx_Policy-optionsPolicy-statementName" "vmx_135" {
	resource_name = "vmx_135"
	name = "SPINE_TO_LEAF_EVPN_OUT"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_136" {
	resource_name = "vmx_136"
	name = "SPINE_TO_LEAF_EVPN_OUT-10"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunity-name" "vmx_137" {
	resource_name = "vmx_137"
	name = "FROM_SPINE_EVPN_TIER"
}

resource "junos-vmx_Policy-optionsPolicy-statementName" "vmx_138" {
	resource_name = "vmx_138"
	name = "SPINE_TO_LEAF_FABRIC_OUT"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermName" "vmx_139" {
	resource_name = "vmx_139"
	name = "SPINE_TO_LEAF_FABRIC_OUT-10"
}

resource "junos-vmx_Policy-optionsPolicy-statementTermThenCommunity-name" "vmx_140" {
	resource_name = "vmx_140"
	name = "FROM_SPINE_FABRIC_TIER"
}

resource "junos-vmx_Policy-optionsCommunityName" "vmx_141"{
	resource_name = "vmx_141"
	name = "DEFAULT_DIRECT_V4"
}

resource "junos-vmx_Policy-optionsCommunityMembers" "vmx_142"{
	resource_name = "vmx_142"
	members = "1:20007"
	members__1 = "21001:26000"
}

resource "junos-vmx_Policy-optionsCommunityName" "vmx_143"{
	resource_name = "vmx_143"
	name = "FROM_SPINE_EVPN_TIER"
}

resource "junos-vmx_Policy-optionsCommunityMembers" "vmx_144"{
	resource_name = "vmx_144"
	members = "0:14"
}

resource "junos-vmx_Policy-optionsCommunityName" "vmx_145"{
	resource_name = "vmx_145"
	name = "FROM_SPINE_FABRIC_TIER"
}

resource "junos-vmx_Policy-optionsCommunityMembers" "vmx_146"{
	resource_name = "vmx_146"
	members = "0:15"
}

resource "junos-vmx_Routing-instancesInstanceName" "vmx_147"{
	resource_name = "vmx_147"
	name = "mgmt_junos"
}

resource "junos-vmx_Routing-instancesInstanceRouting-optionsStaticRouteName" "vmx_148"{
	resource_name = "vmx_148"
	name = "0.0.0.0/0"
}

resource "junos-vmx_Routing-instancesInstanceRouting-optionsStaticRouteNext-hop" "vmx_149"{
	resource_name = "vmx_149"
	next-hop = "10.56.31.254"
}

resource "junos-vmx_Routing-optionsRouter-id" "vmx_150"{
	resource_name = "vmx_150"
	router-id = "10.10.11.0"
}

resource "junos-vmx_Routing-optionsAutonomous-systemAs-number" "vmx_151"{
	resource_name = "vmx_151"
	as-number = "65001"
}

resource "junos-vmx_Routing-optionsForwarding-tableExport" "vmx_152"{
	resource_name = "vmx_152"
	export = "PFE-LB"
}

resource "junos-vmx_ProtocolsBgpGroupName" "vmx_153"{
	resource_name = "vmx_153"
	name = "l3clos-s"
}

resource "junos-vmx_ProtocolsBgpGroupType" "vmx_154"{
	resource_name = "vmx_154"
	type = "external"
}

resource "junos-vmx_ProtocolsBgpGroupBfd-liveness-detectionMinimal-interval" "vmx_155"{
	resource_name = "vmx_155"
	minimum-interval = "1000"
}

resource "junos-vmx_ProtocolsBgpGroupBfd-liveness-detectionMultiplier" "vmx_156"{
	resource_name = "vmx_155"
	multiplier = "3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_157"{
	resource_name = "vmx_157"
	name = "10.10.10.1"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_158"{
	resource_name = "vmx_158"
	description = "facing_mayank-rack-001-leaf1"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_159"{
	resource_name = "vmx_159"
	local-address = "10.10.10.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_160"{
	resource_name = "vmx_160"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_161"{
	resource_name = "vmx_161"
	peer-as = "65003"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_162"{
	resource_name = "vmx_162"
	name = "10.10.10.3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_163"{
	resource_name = "vmx_163"
	description = "facing_mayank-rack-001-leaf2"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_164"{
	resource_name = "vmx_164"
	local-address = "10.10.10.2"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_165"{
	resource_name = "vmx_165"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_166"{
	resource_name = "vmx_166"
	peer-as = "65004"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_167"{
	resource_name = "vmx_167"
	name = "10.10.10.5"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_168"{
	resource_name = "vmx_168"
	description = "facing_mayank-rack-001-leaf3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_169"{
	resource_name = "vmx_169"
	local-address = "10.10.10.4"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_170"{
	resource_name = "vmx_170"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_171"{
	resource_name = "vmx_171"
	peer-as = "65005"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_172"{
	resource_name = "vmx_172"
	name = "10.10.10.7"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_173"{
	resource_name = "vmx_173"
	description = "facing_mayank-rack-001-leaf4"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_174"{
	resource_name = "vmx_174"
	local-address = "10.10.10.6"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_175"{
	resource_name = "vmx_175"
	export = "( SPINE_TO_LEAF_FABRIC_OUT &amp;&amp; BGP-AOS-Policy )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_176"{
	resource_name = "vmx_176"
	peer-as = "65006"
}

resource "junos-vmx_ProtocolsBgpGroupName" "vmx_177"{
	resource_name = "vmx_177"
	name = "l3clos-s-evpn"
}

resource "junos-vmx_ProtocolsBgpGroupType" "vmx_178"{
	resource_name = "vmx_178"
	type = "external"
}

resource "junos-vmx_ProtocolsBgpGroupMultihopTtl" "vmx_179"{
	resource_name = "vmx_179"
	ttl = "1"
}

resource "junos-vmx_ProtocolsBgpGroupFamilyEvpnSignalingLoopsLoops" "vmx_180"{
	resource_name = "vmx_180"
	loops = "2"
}

resource "junos-vmx_ProtocolsBgpGroupBfd-liveness-detectionMinimal-interval" "vmx_181"{
	resource_name = "vmx_181"
	minimum-interval = "3000"
}

resource "junos-vmx_ProtocolsBgpGroupBfd-liveness-detectionMultiplier" "vmx_182"{
	resource_name = "vmx_182"
	multiplier = "3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_183"{
	resource_name = "vmx_183"
	name = "10.10.11.2"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_184"{
	resource_name = "vmx_184"
	description = "facing_mayank-rack-001-leaf1-evpn-overlay"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_185"{
	resource_name = "vmx_185"
	local-address = "10.10.11.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_186"{
	resource_name = "vmx_186"
	export = "( SPINE_TO_LEAF_EVPN_OUT )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_187"{
	resource_name = "vmx_187"
	peer-as = "65003"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_188"{
	resource_name = "vmx_188"
	name = "10.10.11.3"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_189"{
	resource_name = "vmx_189"
	description = "facing_mayank-rack-001-leaf2-evpn-overlay"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_190"{
	resource_name = "vmx_190"
	local-address = "10.10.11.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_191"{
	resource_name = "vmx_191"
	export = "( SPINE_TO_LEAF_EVPN_OUT )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_192"{
	resource_name = "vmx_192"
	peer-as = "65004"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_193"{
	resource_name = "vmx_183"
	name = "10.10.11.4"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_194"{
	resource_name = "vmx_194"
	description = "facing_mayank-rack-001-leaf3-evpn-overlay"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_195"{
	resource_name = "vmx_195"
	local-address = "10.10.11.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_196"{
	resource_name = "vmx_196"
	export = "( SPINE_TO_LEAF_EVPN_OUT )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_197"{
	resource_name = "vmx_197"
	peer-as = "65005"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborName" "vmx_198"{
	resource_name = "vmx_198"
	name = "10.10.11.5"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborDescription" "vmx_199"{
	resource_name = "vmx_199"
	description = "facing_mayank-rack-004-leaf1-evpn-overlay"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborLocal-address" "vmx_200"{
	resource_name = "vmx_200"
	local-address = "10.10.11.0"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborExport" "vmx_201"{
	resource_name = "vmx_201"
	export = "( SPINE_TO_LEAF_EVPN_OUT )"
}

resource "junos-vmx_ProtocolsBgpGroupNeighborPeer-as" "vmx_202"{
	resource_name = "vmx_202"
	peer-as = "65006"
}

resource "junos-vmx_ProtocolsLldpPort-id-subtype" "vmx_203"{
	resource_name = "vmx_203"
	port-id-subtype = "interface-name"
}

resource "junos-vmx_ProtocolsLldpPort-description-type" "vmx_204"{
	resource_name = "vmx_204"
	port-description-type = "interface-description"
}

resource "junos-vmx_ProtocolsLldpNeighbor-port-info-display" "vmx_205"{
	resource_name = "vmx_205"
	neighbor-port-info-display = "port-id"
}

resource "junos-vmx_ProtocolsLldpInterfaceName" "vmx_206"{
	resource_name = "vmx_206"
	name = "all"
}