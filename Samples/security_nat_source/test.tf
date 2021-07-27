provider "junos-device" {
    host = "172.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_SecurityNatSourcePoolAddressToIpaddr" "ipaddr1" {
    resource_name = "XYZ"
    name = "ut-pool"
    name__1 = "107.127.96.40/32"
    ipaddr = "107.127.96.43/32"
}

resource "junos-device_SecurityNatSourcePoolAddressToIpaddr" "ipaddr2" {
    resource_name = "XYZ"
    name = "t-pool"
    name__1 = "172.16.0.8/32"
    ipaddr = "172.16.0.10/32"
}

resource "junos-device_SecurityNatSourceRule__SetFromZone" "fromzone1" {
    resource_name = "XYZ"
    name = "snat-untrust-to-trust"
    zone = "untrust"
}

resource "junos-device_SecurityNatSourceRule__SetFromZone" "fromzone2" {
    resource_name = "XYZ"
    name = "snat-trust-to-untrust"
    zone = "trust"
}

resource "junos-device_SecurityNatSourceRule__SetToZone" "tozone1" {
    resource_name = "XYZ"
    name = "snat-untrust-to-trust"
    zone = "trust"
}

resource "junos-device_SecurityNatSourceRule__SetToZone" "tozone2" {
    resource_name = "XYZ"
    name = "snat-trust-to-untrust"
    zone = "trust"
}

resource "junos-device_SecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address" "match1" {
    resource_name = "XYZ"
    name = "snat-untrust-to-trust"
    name__1 = "1"
    source__address = "0.0.0.0/0"
}

resource "junos-device_SecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address" "match2" {
    resource_name = "XYZ"
    name = "snat-trust-to-untrust"
    name__1 = "2"
    source__address = "10.0.0.0/17"
}

resource "junos-device_SecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name" "poolname1" {
    resource_name = "XYZ"
    name = "snat-untrust-to-trust"
    name__1 = "1"
    pool__name = "t-pool"
}

resource "junos-device_SecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name" "poolname2" {
    resource_name = "XYZ"
    name = "snat-trust-to-untrust"
    name__1 = "2"
    pool__name = "ut-pool"
}

#Commit is to be done at the end

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_SecurityNatSourcePoolAddressToIpaddr.ipaddr1,
    ]
}




