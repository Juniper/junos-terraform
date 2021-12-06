# ----------------- #
# addressbook IP
# ----------------- #
resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "SecurityAddress__BookAddressIp__Prefix_1" {
	resource_name = "JTAF-ADDRESSBOOK"
    	name = "global"
    	name__1 = "source-1"
    	ip__prefix = "10.0.7.9/32"
}

resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "SecurityAddress__BookAddressIp__Prefix_2" {
	resource_name = "JTAF-ADDRESSBOOK"
    	name = "global"
    	name__1 = "dest-1"
    	ip__prefix = "10.0.48.24/32"
}

resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "SecurityAddress__BookAddressIp__Prefix_3" {
	resource_name = "JTAF-ADDRESSBOOK"
	name = "global"
	name__1 = "source-2"	
	ip__prefix = "10.0.9.9/32"
}

resource "junos-vsrx_SecurityAddress__BookAddressIp__Prefix" "SecurityAddress__BookAddressIp__Prefix_4" {
    	resource_name = "JTAF-ADDRESSBOOK"
    	name = "global"
    	name__1 = "dest-2"
    	ip__prefix = "10.0.10.9/32"
}
