# ------------------ #
# Address set config
# ------------------ #

resource "junos-vsrx_SecurityAddress__BookAddress__SetAddressName" "SecurityAddress__BookAddress__SetAddressName_1" {
    resource_name = "JTAF-ADDRESS_SET"
    name = "global"
    name__1 = "address-set-source"
    name__2 = "source-1"
}

resource "junos-vsrx_SecurityAddress__BookAddress__SetAddressName" "SecurityAddress__BookAddress__SetAddressName_2" {
    resource_name = "JTAF-ADDRESS_SET"
    name = "global"
    name__1 = "address-set-dest"
    name__2 = "dest-1"
}

