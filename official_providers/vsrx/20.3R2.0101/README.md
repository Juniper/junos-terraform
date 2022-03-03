# terraform-provider-junos-vsrx 20.3R2.01 release 01

This is the Terraform Provider documentation for the vSRX 20.3R2.01 release 01 provider on the Hashicorp Terraform registry.

For correct usage with Junos, be sure to us the commit module pattern, found [here](https://github.com/Juniper/junos-terraform/tree/master/Samples/tf_module_template/module). It documents a powerful pattern to use with Terraform modules, which provide the commit sequences required with the Junos NETCONF transaction system. Also be sure to read the guide located [here](https://github.com/Juniper/junos-terraform/wiki) for further insight in using Junos and Terraform together.

It's important to name each `resource_group` something sequential like `my_app_N`, where N is a monotonic number (that increases). You can also call the resource the same name to keep it simple. An example is below for this naming discipline.

```bash
resource "junos-vsrx_InterfacesInterfaceDescription" "vsrx_1" {
    resource_name = "vsrx_1"
    name = "ge-0/0/0"
    description = "LEFT-UNTRUST"
}

resource "junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName" "vsrx_2" {
    resource_name = "vsrx_2"
    name = "ge-0/0/0"
    name__1 = "0"
    name__2 = "10.0.0.1/24"
}
```

The reason behind this is due to the life-cycle interaction between Terraform and Junos. Each Terraform resource becomes a separate configuration group within a Junos configuration which is inherited. When a change is made to a resource, the group is destroyed and re-created. It is therefore imperative to separate each resource so that upon destruction, only the single resource is effected.

## Provider Content

There are 51 resources contained in this provider.

```bash
junos-vsrx_ApplicationsApplicationDestination__Port
junos-vsrx_ApplicationsApplicationProtocol
junos-vsrx_commit
junos-vsrx_destroycommit
junos-vsrx_FirewallFilterTermFromProtocol
junos-vsrx_FirewallFilterTermThenAccept
junos-vsrx_FirewallFilterTermThenSample
junos-vsrx_Forwarding__OptionsSamplingFamilyInetOutputFileFilename
junos-vsrx_Forwarding__OptionsSamplingInputRate
junos-vsrx_InterfacesInterfaceDescription
junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName
junos-vsrx_Policy__OptionsPolicy__StatementTermFromInstance
junos-vsrx_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress
junos-vsrx_Policy__OptionsPolicy__StatementTermThenAccept
junos-vsrx_Policy__OptionsPolicy__StatementTermThenReject
junos-vsrx_Routing__InstancesInstanceInstance__Type
junos-vsrx_Routing__InstancesInstanceInterfaceName
junos-vsrx_Routing__InstancesInstanceRouting__OptionsInstance__Import
junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Hop
junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Table
junos-vsrx_SecurityAddress__BookAddress__SetAddressName
junos-vsrx_SecurityAddress__BookAddressIp__Prefix
junos-vsrx_SecurityAddress__BookAddressRange__AddressToRange__High
junos-vsrx_SecurityNatDestinationPoolAddressIpaddr
junos-vsrx_SecurityNatDestinationPoolAddressPort
junos-vsrx_SecurityNatDestinationPoolRouting__InstanceRi__Name
junos-vsrx_SecurityNatDestinationRule__SetFromInterface
junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__AddressDst__Addr
junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName
junos-vsrx_SecurityNatDestinationRule__SetRuleThenDestination__NatPoolPool__Name
junos-vsrx_SecurityNatProxy__ArpInterfaceAddressToIpaddr
junos-vsrx_SecurityNatSourcePoolAddressToIpaddr
junos-vsrx_SecurityNatSourceRule__SetFromZone
junos-vsrx_SecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address
junos-vsrx_SecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name
junos-vsrx_SecurityNatSourceRule__SetToZone
junos-vsrx_SecurityPoliciesPolicyPolicyMatchApplication
junos-vsrx_SecurityPoliciesPolicyPolicyMatchDestination__Address
junos-vsrx_SecurityPoliciesPolicyPolicyMatchSource__Address
junos-vsrx_SecurityPoliciesPolicyPolicyThenCountApply__Groups
junos-vsrx_SecurityPoliciesPolicyPolicyThenDeny
junos-vsrx_SecurityPoliciesPolicyPolicyThenLogSession__Init
junos-vsrx_SecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy
junos-vsrx_SecurityZonesSecurity__ZoneHost__Inbound__TrafficSystem__ServicesName
junos-vsrx_SecurityZonesSecurity__ZoneInterfacesName
junos-vsrx_SystemRoot__AuthenticationEncrypted__Password
junos-vsrx_SystemServicesSshPort
junos-vsrx_SystemServicesWeb__ManagementHttpInterface
junos-vsrx_SystemServicesWeb__ManagementHttpsInterface
junos-vsrx_SystemServicesWeb__ManagementHttpsSystem__Generated__Certificate
junos-vsrx_SystemSyslogFileContentsAny
```

## Resource Descriptions

__resource_name:__ junos-vsrx_ApplicationsApplicationDestination__Port

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | application name |
| destination_port | string | destination port for application |

---

__resource_name:__ junos-vsrx_ApplicationsApplicationProtocol

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | application name |
| protocol | string | protocol type, like "tcp" |

---

__resource_name:__ junos-vsrx_commit

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |

This is a 'presence' type resource. It doesn't create remote state, but triggers an RPC when `terraform apply` is ran.

---

__resource_name:__ junos-vsrx_destroycommit

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |

This is a 'presence' type resource. It doesn't create remote state, but triggers an RPC when `terraform destroy` is ran.

---

__resource_name:__ junos-vsrx_FirewallFilterTermFromProtocol

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Firewall Filter Name |
| name__1 | string | Firewall Filter Term |
| protocol | string | Firewall Filter From Protocol |

---

__resource_name:__ junos-vsrx_FirewallFilterTermThenAccept

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Firewall Filter Name |
| name__1 | string | Firewall Filter Term |
| accept | string | Firewall Filter Then Accept |

*Accept is a presence key, use it like this:* `accept = ""`

---

__resource_name:__ junos-vsrx_FirewallFilterTermThenSample

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Firewall Filter Name |
| name__1 | string | Firewall Filter Term |
| sample | string | Firewall Filter Then Sample |

*Sample is a presence key, use it like this:* `sample = ""`

---

__resource_name:__ junos-vsrx_Forwarding__OptionsSamplingFamilyInetOutputFileFilename

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| filename | string | Name of file to contain sampled packet dumps |

---

__resource_name:__ junos-vsrx_Forwarding__OptionsSamplingInputRate

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| rate | string | Ratio of packets to be sampled (1 out of N) |

---

__resource_name:__ junos-vsrx_InterfacesInterfaceDescription

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Name of interface |
| description | string | Text description of interface |

---

__resource_name:__ junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Name of interface |
| name__1 | string | Unit of interface |
| name__2 | string | Interface address/destination prefix |

---

__resource_name:__ junos-vsrx_Policy__OptionsPolicy__StatementTermFromInstance

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Statement |
| name__1 | string | Policy Term |
| instance | string | From routing protocol instance |

---

__resource_name:__ junos-vsrx_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress

Example use case
```bash
resource "junos-vsrx_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress" "vsrx_3" {
    resource_name = "vsrx_3"
    name = "DEF-IMPORT-FWTRUST-TABLE"
    name__1 = "t1"
    address = "10.0.0.0/17"
    // orlonger below, is an enumerated type. The single whitespace turns it on.
    // XPath for this type is: <xpath name="/policy-options/policy-statement/term/from/route-filter/address"/>
    // Some enums require input, in which case use the key like any other and insert the data.
    orlonger = " "
}
```

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Statement |
| name__1 | string | Policy Term |
| address | string | IP address or hostname |

*Each of these that follow may or may not have inputs. Use only one*

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| accept_reject | string | May require valid input or just " " if no input |
| add__path | string | May require valid input or just " " if no input |
| address__mask | string | May require valid input or just " " if no input |
| aggregate__bandwidth | string | May require valid input or just " " if no input |
| aigp__adjust | string | May require valid input or just " " if no input |
| aigp__originate | string | May require valid input or just " " if no input |
| analyze | string | May require valid input or just " " if no input |
| apply__advanced | string | May require valid input or just " " if no input |
| as__path | string | May require valid input or just " " if no input |
| as__path__expand | string | May require valid input or just " " if no input |
| bgp__output__queue__priority | string | May require valid input or just " " if no input |
| class | string | May require valid input or just " " if no input |
| color | string | May require valid input or just " " if no input |
| color2 | string | May require valid input or just " " if no input |
| community | string | May require valid input or just " " if no input |
| cos__next__hop__map | string | May require valid input or just " " if no input |
| damping | string | May require valid input or just " " if no input |
| default__action | string | May require valid input or just " " if no input |
| destination__class | string | May require valid input or just " " if no input |
| dynamic__tunnel__attributes | string | May require valid input or just " " if no input |
| exact | string | May require valid input or just " " if no input |
| external | string | May require valid input or just " " if no input |
| forwarding__class | string | May require valid input or just " " if no input |
| get__route__range | string | May require valid input or just " " if no input |
| install__nexthop | string | May require valid input or just " " if no input |
| install__to__fib | string | May require valid input or just " " if no input |
| label | string | May require valid input or just " " if no input |
| label__allocation | string | May require valid input or just " " if no input |
| label__allocation__fallback__reject | string | May require valid input or just " " if no input |
| limit__bandwidth | string | May require valid input or just " " if no input |
| load__balance | string | May require valid input or just " " if no input |
| local__preference | string | May require valid input or just " " if no input |
| longer | string | May require valid input or just " " if no input |
| map__to__interface | string | May require valid input or just " " if no input |
| metric | string | May require valid input or just " " if no input |
| metric2 | string | May require valid input or just " " if no input |
| metric3 | string | May require valid input or just " " if no input |
| metric4 | string | May require valid input or just " " if no input |
| mhop__bfd__port | string | May require valid input or just " " if no input |
| multipath__resolve | string | May require valid input or just " " if no input |
| next | string | May require valid input or just " " if no input |
| next__hop | string | May require valid input or just " " if no input |
| no__backup | string | May require valid input or just " " if no input |
| no__route__localize | string | May require valid input or just " " if no input |
| no_entropy | string | May require valid input or just " " if no input |
| origin | string | May require valid input or just " " if no input |
| orlonger | string | May require valid input or just " " if no input |
| p2mp__lsp__root | string | May require valid input or just " " if no input |
| preference | string | May require valid input or just " " if no input |
| preference2 | string | May require valid input or just " " if no input |
| prefix__length__range | string | May require valid input or just " " if no input |
| prefix__segment | string | May require valid input or just " " if no input |
| priority | string | May require valid input or just " " if no input |
| resolution__map | string | May require valid input or just " " if no input |
| selected__mldp__egress | string | May require valid input or just " " if no input |
| source__class | string | May require valid input or just " " if no input |
| sr__te__template | string | May require valid input or just " " if no input |
| ssm__source | string | May require valid input or just " " if no input |
| tag | string | May require valid input or just " " if no input |
| tag2 | string | May require valid input or just " " if no input |
| through | string | May require valid input or just " " if no input |
| trace | string | May require valid input or just " " if no input |
| tunnel__attribute | string | May require valid input or just " " if no input |
| tunnel__end__point__address | string | May require valid input or just " " if no input |
| upto | string | May require valid input or just " " if no input |
| validation__state | string | May require valid input or just " " if no input |

---

__resource_name:__ junos-vsrx_Policy__OptionsPolicy__StatementTermThenAccept

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Statement |
| name__1 | string | Policy Term |
| accept | string | Accept a route |

*Accept is a presence key, use it like this:* `accept = ""`

---

__resource_name:__ junos-vsrx_Policy__OptionsPolicy__StatementTermThenReject

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Statement |
| name__1 | string | Policy Term |
| reject | string | Reject a route |

*Reject is a presence key, use it like this:* `reject = ""`

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceInstance__Type

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| instance__type | string | Type of routing instance |

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceInstance__Type

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| instance__type | string | Type of routing instance |

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceInterfaceName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| name__1 | string | Interface name |

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceRouting__OptionsInstance__Import

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| instance__import | string | Import policy for instance RIBs |

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Hop

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| name__1 | string | Route |
| next__hop | string | Next hop to destination |

---

__resource_name:__ junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Table

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Routing instance name |
| name__1 | string | Route |
| next__table | string | Next hop to another table |

---

__resource_name:__ junos-vsrx_SecurityAddress__BookAddress__SetAddressName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address book |
| name__1 | string | Address set |
| name__2 | string | Security address name |

---

__resource_name:__ junos-vsrx_SecurityAddress__BookAddressIp__Prefix

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address book |
| name__1 | string | Address |
| ip__prefix | string | Numeric IPv4 or IPv6 address with prefix|

---

__resource_name:__ junos-vsrx_SecurityAddress__BookAddressRange__AddressToRange__High

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address book |
| name__1 | string | Address |
| name__2 | string | Range address |
| range__high | string | Upper limit of address range |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationPoolAddressIpaddr

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address pool |
| ipaddr | string | IPv4 or IPv6 address or address range |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationPoolAddressPort

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address pool |
| port | string | Specify the port value |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationPoolRouting__InstanceRi__Name

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address pool |
| ri__name | string | Routing-instance name |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationRule__SetFromInterface

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| interface | string | Source interface list |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__AddressDst__Addr

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule name |
| dst__addr | string | IPv4 or IPv6 destination address |

---

__resource_name:__ junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule name |
| name__2 | string | Port or lower limit of port range |


---

__resource_name:__ junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule name |
| name__2 | string | Port or lower limit of port range |


---

__resource_name:__ junos-vsrx_SecurityNatDestinationRule__SetRuleThenDestination__NatPoolPool__Name

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule name |
| pool__name | string | Name of Destination NAT pool |


---

__resource_name:__ junos-vsrx_SecurityNatProxy__ArpInterfaceAddressToIpaddr

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Interface |
| name__1 | string | Address |
| ipaddr | string | Upper limit of address range |


---

__resource_name:__ junos-vsrx_SecurityNatSourcePoolAddressToIpaddr

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Address pool |
| name__1 | string | Address |
| ipaddr | string | IPv4 or IPv6 upper limit of address range |


---

__resource_name:__ junos-vsrx_SecurityNatSourceRule__SetFromZone

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| zone | string | Source zone list |


---

__resource_name:__ junos-vsrx_SecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule |
| source__address | string | Source address |


---

__resource_name:__ junos-vsrx_SecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| name__1 | string | Rule |
| pool__name | string | Name of Source NAT pool |


---

__resource_name:__ junos-vsrx_SecurityNatSourceRule__SetToZone

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Rule set |
| zone | string | Destination zone list |


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyMatchDestination__Address

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| destination__address | string | Match destination address |


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyMatchSource__Address

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| source__address | string | Match source address |


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyThenCountApply__Groups

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyThenDeny

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| deny | string | Deny packets |


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyThenDeny

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| deny | string | Deny packets |

*Deny is a presence key, use it like this:* `deny = ""`


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyThenLogSession__Init

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| session__init | string | Log at session init time |

*Session__init is a presence key, use it like this:* `session__init = ""`


---

__resource_name:__ junos-vsrx_SecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Policy Name |
| from__zone__name | string | Zone |
| to__zone__name | string | Zone |
| pair__policy | string | Policy in the reverse direction, to form a pair |


---

__resource_name:__ junos-vsrx_SecurityZonesSecurity__ZoneHost__Inbound__TrafficSystem__ServicesName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Security zone |
| name__1 | string | System services name |


---

__resource_name:__ junos-vsrx_SecurityZonesSecurity__ZoneInterfacesName

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | Security zone |
| name__1 | string | Logical interface |


---

__resource_name:__ junos-vsrx_SystemRoot__AuthenticationEncrypted__Password

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| encrypted__password | string | Encrypted password string |


---

__resource_name:__ junos-vsrx_SystemServicesSshPort

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| port | string | Port number to accept incoming connections |


---

__resource_name:__ junos-vsrx_SystemServicesWeb__ManagementHttpInterface

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| interface | string | Interfaces that accept HTTP access |

*One interface per resource*

---

__resource_name:__ junos-vsrx_SystemServicesWeb__ManagementHttpsInterface

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| interface | string | Interfaces that accept HTTP access |

*One interface per resource*

---

__resource_name:__ junos-vsrx_SystemServicesWeb__ManagementHttpsSystem__Generated__Certificate

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| system__generated__certificate | string | X.509 certificate generated automatically by system |

*Session__init is a presence key, use it like this:* `system__generated__certificate = ""`

---

__resource_name:__ junos-vsrx_SystemSyslogFileContentsAny

| Key Name | Value Type | Description |  
| :---: | :---: | :---: |
| resource_name | string | standard resource_name key |
| name | string | File name |
| name__1 | string | Contents |
| any | string | All levels |

*any is a presence key, use it like this:* `any = ""`