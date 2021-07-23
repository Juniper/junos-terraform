provider "junos-device" {
    host = "10.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_SecurityAddress__BookAddress__SetAddressName" "addr1" {
    resource_name = "XYZ"
    name = "global"
    name__1 = "mvm"
    name__2 = "mvm-aks1-az1"
}

resource "junos-device_SecurityAddress__BookAddress__SetAddressName" "addr2" {
    resource_name = "XYZ"
    name = "global"
    name__1 = "mvm"
    name__2 = "mvm-aks1-az2"
}

resource "junos-device_SecurityAddress__BookAddress__SetAddressName" "addr3" {
    resource_name = "XYZ"
    name = "global"
    name__1 = "mvm"
    name__2 = "mvm-aks2-az1"
}

resource "junos-device_SecurityAddress__BookAddress__SetAddressName" "addr4" {
    resource_name = "XYZ"
    name = "global"
    name__1 = "mvm"
    name__2 = "mvm-aks2-az2"
}

resource "junos-device_SecurityNatSourcePoolAddressToIpaddr" "ipaddr1" {
    resource_name = "XYZ"
    name = "pool1"
    name__1 = "107.0.9.40/32"
    ipaddr = "107.0.9.43/32"
}

resource "junos-device_SecurityNatSourcePoolAddressToIpaddr" "ipaddr2" {
    resource_name = "XYZ"
    name = "pool2"
    name__1 = "172.16.0.8/32"
    ipaddr = "172.16.0.10/32"
}

#Commit is to be done at the end

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_SecurityNatSourcePoolAddressToIpaddr.ipaddr1,
        junos-device_SecurityNatSourcePoolAddressToIpaddr.ipaddr2,
        junos-device_SecurityAddress__BookAddress__SetAddressName.addr1,
        junos-device_SecurityAddress__BookAddress__SetAddressName.addr2,
        junos-device_SecurityAddress__BookAddress__SetAddressName.addr3,
        junos-device_SecurityAddress__BookAddress__SetAddressName.addr4
    ]
}




