provider "junos-device" {
    host = "10.x.x.x"
    port = 22
    username = "user"
    password = "user123"
    sshkey = ""
}

resource "junos-device_commit" "commit2" {
    resource_name = "commit"
    depends_on = [
        junos-device_ProtocolsBgpGroupTraceoptionsFile.demo
    ]
}

//resource "junos-device_ProtocolsBgpGroupTraceoptionsFile" "demo" {
//    resource_name = "XYZ"
//    name = "demo"
//    filename = "temp.log"
//}

resource "junos-device_ProtocolsBgpGroupTraceoptionsFile" "demo" {
    resource_name = "XYZ"
    name = "demo"
    filename = "temp.log"
    size = "3m"
    files = 10
}