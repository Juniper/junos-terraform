provider "junos-device" {
    host = "10.x.x.x"
    port = 22
    username = "user"
    password = "user123"
}

resource "junos-device_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress" "ipaddr1" {
    resource_name = "XYZ"
    name = "DEF-IMPORT-FWTRUST-TABLE"
    name__1 = "t1"
    choice__ident = "exact"
    address = "10.0.0.0/17"
}

#Commit is to be done at the end

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress.ipaddr1,
    ]
}




