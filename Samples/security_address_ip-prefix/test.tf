provider "junos-device" {
    host = "172.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_SecurityAddress__BookAddressIp__Prefix" "ipaddr1" {
    resource_name = "XYZ"
    name = "global"
    name__1 = "F5"
    ip__prefix = "107.10.10.250/32"
}

#Commit is to be done at the end

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_SecurityAddress__BookAddressIp__Prefix.ipaddr1,
    ]
}




