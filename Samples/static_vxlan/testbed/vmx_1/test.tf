terraform {
  required_providers {
    junos-vmx = {
      source = "juniper/providers/junos-vmx"
      version = "22.3"
    }
  }
}

resource "junos-vmx_SystemSyslogUser" "vmx_1" {
resource_name = "vmx_1"
}

resource "junos-vmx_SystemSyslogFileContentsInfo" "vmx_2" {
resource_name = "vmx_2"
name = "/system/syslog/file/name"
name__1 = "/system/syslog/file/contents/name"
info = "/system/syslog/file/contents/info"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsLevel" "vmx_3" {
resource_name = "vmx_3"
level = "/system/processes/dhcp-service/traceoptions/level"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddressName" "vmx_4" {
resource_name = "vmx_4"
name = "/interfaces/interface/name"
name__1 = "/interfaces/interface/unit/name"
name__2 = "/interfaces/interface/unit/family/inet/address/name"
}

resource "junos-vmx_ChassisFpcLite__Mode" "vmx_5" {
resource_name = "vmx_5"
name = "/chassis/fpc/name"
lite__mode = "/chassis/fpc/lite-mode"
}

resource "junos-vmx_ProtocolsRouter__AdvertisementInterface" "vmx_6" {
resource_name = "vmx_6"
}
resource "junos-vmx_SystemRoot__Authentication" "vmx_7" {
resource_name = "vmx_7"
}

resource "junos-vmx_SystemLoginUserUid" "vmx_8" {
resource_name = "vmx_8"
name = "/system/login/user/name"
uid = "/system/login/user/uid"
}

resource "junos-vmx_SystemServices" "vmx_9" {
resource_name = "vmx_9"
}

resource "junos-vmx_SystemSyslogFileName" "vmx_10" {
resource_name = "vmx_10"
name = "/system/syslog/file/name"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsFileSize" "vmx_11" {
resource_name = "vmx_11"
size = "/system/processes/dhcp-service/traceoptions/file/size"
}

resource "junos-vmx_SystemRoot__AuthenticationEncrypted__Password" "vmx_12" {
resource_name = "vmx_12"
encrypted__password = "/system/root-authentication/encrypted-password"
}

resource "junos-vmx_SystemLoginUserAuthentication" "vmx_13" {
resource_name = "vmx_13"
name = "/system/login/user/name"
}

resource "junos-vmx_SystemSyslogUserContentsName" "vmx_14" {
resource_name = "vmx_14"
name = "/system/syslog/user/name"
name__1 = "/system/syslog/user/contents/name"
}

resource "junos-vmx_SystemSyslogFile" "vmx_15" {
resource_name = "vmx_15"
}

resource "junos-vmx_SystemSyslogFileContents" "vmx_16" {
resource_name = "vmx_16"
name = "/system/syslog/file/name"
}

resource "junos-vmx_SystemServicesRestHttpPort" "vmx_17" {
resource_name = "vmx_17"
port = "/system/services/rest/http/port"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsFlagName" "vmx_18" {
resource_name = "vmx_18"
name = "/system/processes/dhcp-service/traceoptions/flag/name"
}

resource "junos-vmx_SystemLogin" "vmx_19" {
resource_name = "vmx_19"
}

resource "junos-vmx_SystemLoginUser" "vmx_20" {
resource_name = "vmx_20"
}

resource "junos-vmx_SystemSyslogUserContents" "vmx_21" {
resource_name = "vmx_21"
name = "/system/syslog/user/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInetAddress" "vmx_22" {
resource_name = "vmx_22"
name = "/interfaces/interface/name"
name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_Routing__OptionsStaticRouteNext__Hop" "vmx_23" {
resource_name = "vmx_23"
name = "/routing-options/static/route/name"
next__hop = "/routing-options/static/route/next-hop"
}

resource "junos-vmx_SystemScriptsLanguage" "vmx_24" {
resource_name = "vmx_24"
language = "/system/scripts/language"
}

resource "junos-vmx_SystemLoginUserAuthenticationEncrypted__Password" "vmx_25" {
resource_name = "vmx_25"
name = "/system/login/user/name"
encrypted__password = "/system/login/user/authentication/encrypted-password"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptions" "vmx_26" {
resource_name = "vmx_26"
}

resource "junos-vmx_InterfacesInterfaceUnitName" "vmx_27" {
resource_name = "vmx_27"
name = "/interfaces/interface/name"
name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_InterfacesInterfaceUnitFamily" "vmx_28" {
resource_name = "vmx_28"
name = "/interfaces/interface/name"
name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_SystemScripts" "vmx_29" {
resource_name = "vmx_29"
}

resource "junos-vmx_SystemLoginUserClass" "vmx_30" {
resource_name = "vmx_30"
name = "/system/login/user/name"
class = "/system/login/user/class"
}

resource "junos-vmx_SystemServicesRestEnable__Explorer" "vmx_31" {
resource_name = "vmx_31"
enable__explorer = "/system/services/rest/enable-explorer"
}

resource "junos-vmx_SystemSyslog" "vmx_32" {
resource_name = "vmx_32"
}

resource "junos-vmx_SystemProcesses" "vmx_33" {
resource_name = "vmx_33"
}

resource "junos-vmx_ProtocolsRouter__Advertisement" "vmx_34" {
resource_name = "vmx_34"
}

resource "junos-vmx_ProtocolsRouter__AdvertisementInterfaceName" "vmx_35" {
resource_name = "vmx_35"
name = "/protocols/router-advertisement/interface/name"
}

resource "junos-vmx_SystemLoginUserName" "vmx_36" {
resource_name = "vmx_36"
name = "/system/login/user/name"
}

resource "junos-vmx_ChassisFpc" "vmx_37" {
resource_name = "vmx_37"
}

resource "junos-vmx_ChassisFpcPicNumber__Of__Ports" "vmx_38" {
resource_name = "vmx_38"
name = "/chassis/fpc/name"
name__1 = "/chassis/fpc/pic/name"
number__of__ports = "/chassis/fpc/pic/number-of-ports"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsFileFilename" "vmx_39" {
resource_name = "vmx_39"
filename = "/system/processes/dhcp-service/traceoptions/file/filename"
}

resource "junos-vmx_Routing__OptionsStaticRoute" "vmx_40" {
resource_name = "vmx_40"
}

resource "junos-vmx_SystemServicesSshRoot__Login" "vmx_41" {
resource_name = "vmx_41"
root__login = "/system/services/ssh/root-login"
}

resource "junos-vmx_ChassisFpcPic" "vmx_42" {
resource_name = "vmx_42"
name = "/chassis/fpc/name"
}

resource "junos-vmx_ChassisFpcPicName" "vmx_43" {
resource_name = "vmx_43"
name = "/chassis/fpc/name"
name__1 = "/chassis/fpc/pic/name"
}

resource "junos-vmx_InterfacesInterface" "vmx_44" {
resource_name = "vmx_44"
}

resource "junos-vmx_InterfacesInterfaceUnitFamilyInet" "vmx_45" {
resource_name = "vmx_45"
name = "/interfaces/interface/name"
name__1 = "/interfaces/interface/unit/name"
}

resource "junos-vmx_Routing__OptionsStaticRouteName" "vmx_46" {
resource_name = "vmx_46"
name = "/routing-options/static/route/name"
}

resource "junos-vmx_SystemServicesNetconfSsh" "vmx_47" {
resource_name = "vmx_47"
}

resource "junos-vmx_SystemServicesRestHttp" "vmx_48" {
resource_name = "vmx_48"
}

resource "junos-vmx_SystemSyslogFileContentsName" "vmx_49" {
resource_name = "vmx_49"
name = "/system/syslog/file/name"
name__1 = "/system/syslog/file/contents/name"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsFlag" "vmx_50" {
resource_name = "vmx_50"
}

resource "junos-vmx_SystemHost__Name" "vmx_51" {
resource_name = "vmx_51"
host__name = "/host-name"
}

resource "junos-vmx_SystemSyslogFileContentsNotice" "vmx_52" {
resource_name = "vmx_52"
name = "/system/syslog/file/name"
name__1 = "/system/syslog/file/contents/name"
notice = "/system/syslog/file/contents/notice"
}

resource "junos-vmx_Routing__OptionsStatic" "vmx_53" {
resource_name = "vmx_53"
}

resource "junos-vmx_SystemServicesSsh" "vmx_54" {
resource_name = "vmx_54"
}

resource "junos-vmx_SystemServicesNetconf" "vmx_55" {
resource_name = "vmx_55"
}

resource "junos-vmx_SystemServicesRest" "vmx_56" {
resource_name = "vmx_56"
}

resource "junos-vmx_SystemSyslogUserContentsEmergency" "vmx_57" {
resource_name = "vmx_57"
name = "/system/syslog/user/name"
name__1 = "/system/syslog/user/contents/name"
emergency = "/system/syslog/user/contents/emergency"
}

resource "junos-vmx_ChassisFpcName" "vmx_58" {
resource_name = "vmx_58"
name = "/chassis/fpc/name"
}

resource "junos-vmx_SystemSyslogUserName" "vmx_59" {
resource_name = "vmx_59"
name = "/system/syslog/user/name"
}

resource "junos-vmx_SystemProcessesDhcp__Service" "vmx_60" {
resource_name = "vmx_60"
}

resource "junos-vmx_InterfacesInterfaceName" "vmx_61" {
resource_name = "vmx_61"
name = "/interfaces/interface/name"
}

resource "junos-vmx_SystemSyslogFileContentsAny" "vmx_62" {
resource_name = "vmx_62"
name = "/system/syslog/file/name"
name__1 = "/system/syslog/file/contents/name"
any = "/system/syslog/file/contents/any"
}

resource "junos-vmx_SystemProcessesDhcp__ServiceTraceoptionsFile" "vmx_63" {
resource_name = "vmx_63"
}

resource "junos-vmx_InterfacesInterfaceUnit" "vmx_64" {
resource_name = "vmx_64"
name = "/interfaces/interface/name"
}