# -------------------------- #
# configure traffic selector 
# -------------------------- #

resource "junos-vsrx_SecurityIpsecVpnTraffic__SelectorRemote__Ip" "TrafficSelector_VPN1" {
    resource_name = "JTAF_IPSEC1"
    name = "ipsec-1"
    remote__ip = "30.30.30.0/24"
    name__1 = "ipsec-1-TS-1"
    depends_on = [
        junos-vsrx_SecurityIpsecVpnBind__Interface.ipsecBindIntf,
        junos-vsrx_SecurityIpsecVpnIkeGateway.ipsecIkeGw
     ]
}

resource "junos-vsrx_SecurityIpsecVpnTraffic__SelectorLocal__Ip" "TrafficSelector_VPN2" {
    resource_name = "JTAF_IPSEC2"
    name = "ipsec-1"
    local__ip = "20.20.20.0/24"
    name__1 = "ipsec-1-TS-1"
    depends_on = [
        junos-vsrx_SecurityIpsecVpnBind__Interface.ipsecBindIntf,
        junos-vsrx_SecurityIpsecVpnIkeGateway.ipsecIkeGw
     ]
}



