provider "junos-device" {
    host = "10.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr" "addr1" {
    resource_name = "XYZ"
    name = "ge-0/0/1.0"
    name__1 = "10.0.0.17/24"
}

resource "junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr" "addr2" {
    resource_name = "XYZ"
    name = "ge-0/0/0.0"
    name__1 = "10.10.0.64/27"
}

resource "junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr" "addr3" {
    resource_name = "XYZ"
    name = "ge-0/0/0.0"
    name__1 = "10.1.1.23/24"
}


resource "junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr" "addr4" {
    resource_name = "XYZ"
    name = "ge-0/0/0.0"
    name__1 = "10.0.0.40/32"
    ipaddr = "10.0.0.43/32"
}

#Commit is to be done at the end

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr.addr1,
        junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr.addr2,
        junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr.addr3,
        junos-device_SecurityNatProxy__ArpInterfaceAddressToIpaddr.addr4
    ]
}




