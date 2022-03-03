// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// THIS SOFTWARE IS PROVIDED BY Juniper Networks, Inc. ''AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL Juniper Networks, Inc. BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
	"os"

	gonetconf "github.com/davedotdev/go-netconf/helpers/junos_helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProviderConfig is to hold client information
type ProviderConfig struct {
	*gonetconf.GoNCClient
	Host string
}

func check(err error) {
	if err != nil {
		// Some of these errors will be "normal".
		f, _ := os.OpenFile("jtaf_logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(err.Error() + "\n")
		f.Close()
		return
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:     d.Get("host").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		SSHKey:   d.Get("sshkey").(string),
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return &ProviderConfig{client, config.Host}, nil
}

// Provider returns a Terraform ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"junos-vsrx_SystemRoot__AuthenticationEncrypted__Password":                                          junosSystemRoot__AuthenticationEncrypted__Password(),
			"junos-vsrx_SystemServicesSshPort":                                                                  junosSystemServicesSshPort(),
			"junos-vsrx_SystemServicesWeb__ManagementHttpsSystem__Generated__Certificate":                       junosSystemServicesWeb__ManagementHttpsSystem__Generated__Certificate(),
			"junos-vsrx_SystemServicesWeb__ManagementHttpsInterface":                                            junosSystemServicesWeb__ManagementHttpsInterface(),
			"junos-vsrx_SystemServicesWeb__ManagementHttpInterface":                                             junosSystemServicesWeb__ManagementHttpInterface(),
			"junos-vsrx_SystemSyslogFileContentsAny":                                                            junosSystemSyslogFileContentsAny(),
			"junos-vsrx_Policy__OptionsPolicy__StatementTermFromInstance":                                       junosPolicy__OptionsPolicy__StatementTermFromInstance(),
			"junos-vsrx_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress":                           junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress(),
			"junos-vsrx_Policy__OptionsPolicy__StatementTermThenAccept":                                         junosPolicy__OptionsPolicy__StatementTermThenAccept(),
			"junos-vsrx_Policy__OptionsPolicy__StatementTermThenReject":                                         junosPolicy__OptionsPolicy__StatementTermThenReject(),
			"junos-vsrx_ApplicationsApplicationProtocol":                                                        junosApplicationsApplicationProtocol(),
			"junos-vsrx_ApplicationsApplicationDestination__Port":                                               junosApplicationsApplicationDestination__Port(),
			"junos-vsrx_InterfacesInterfaceDescription":                                                         junosInterfacesInterfaceDescription(),
			"junos-vsrx_InterfacesInterfaceUnitFamilyInetAddressName":                                           junosInterfacesInterfaceUnitFamilyInetAddressName(),
			"junos-vsrx_Forwarding__OptionsSamplingInputRate":                                                   junosForwarding__OptionsSamplingInputRate(),
			"junos-vsrx_Forwarding__OptionsSamplingFamilyInetOutputFileFilename":                                junosForwarding__OptionsSamplingFamilyInetOutputFileFilename(),
			"junos-vsrx_FirewallFilterTermFromProtocol":                                                         junosFirewallFilterTermFromProtocol(),
			"junos-vsrx_FirewallFilterTermThenSample":                                                           junosFirewallFilterTermThenSample(),
			"junos-vsrx_FirewallFilterTermThenAccept":                                                           junosFirewallFilterTermThenAccept(),
			"junos-vsrx_Routing__InstancesInstanceInstance__Type":                                               junosRouting__InstancesInstanceInstance__Type(),
			"junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Hop":                         junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Hop(),
			"junos-vsrx_Routing__InstancesInstanceRouting__OptionsStaticRouteNext__Table":                       junosRouting__InstancesInstanceRouting__OptionsStaticRouteNext__Table(),
			"junos-vsrx_Routing__InstancesInstanceInterfaceName":                                                junosRouting__InstancesInstanceInterfaceName(),
			"junos-vsrx_Routing__InstancesInstanceRouting__OptionsInstance__Import":                             junosRouting__InstancesInstanceRouting__OptionsInstance__Import(),
			"junos-vsrx_SecurityAddress__BookAddressIp__Prefix":                                                 junosSecurityAddress__BookAddressIp__Prefix(),
			"junos-vsrx_SecurityAddress__BookAddressRange__AddressToRange__High":                                junosSecurityAddress__BookAddressRange__AddressToRange__High(),
			"junos-vsrx_SecurityAddress__BookAddress__SetAddressName":                                           junosSecurityAddress__BookAddress__SetAddressName(),
			"junos-vsrx_SecurityNatSourcePoolAddressToIpaddr":                                                   junosSecurityNatSourcePoolAddressToIpaddr(),
			"junos-vsrx_SecurityNatSourceRule__SetFromZone":                                                     junosSecurityNatSourceRule__SetFromZone(),
			"junos-vsrx_SecurityNatSourceRule__SetToZone":                                                       junosSecurityNatSourceRule__SetToZone(),
			"junos-vsrx_SecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address":                     junosSecurityNatSourceRule__SetRuleSrc__Nat__Rule__MatchSource__Address(),
			"junos-vsrx_SecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name":                            junosSecurityNatSourceRule__SetRuleThenSource__NatPoolPool__Name(),
			"junos-vsrx_SecurityNatDestinationPoolRouting__InstanceRi__Name":                                    junosSecurityNatDestinationPoolRouting__InstanceRi__Name(),
			"junos-vsrx_SecurityNatDestinationPoolAddressIpaddr":                                                junosSecurityNatDestinationPoolAddressIpaddr(),
			"junos-vsrx_SecurityNatDestinationPoolAddressPort":                                                  junosSecurityNatDestinationPoolAddressPort(),
			"junos-vsrx_SecurityNatDestinationRule__SetFromInterface":                                           junosSecurityNatDestinationRule__SetFromInterface(),
			"junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__AddressDst__Addr": junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__AddressDst__Addr(),
			"junos-vsrx_SecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName":         junosSecurityNatDestinationRule__SetRuleDest__Nat__Rule__MatchDestination__PortName(),
			"junos-vsrx_SecurityNatDestinationRule__SetRuleThenDestination__NatPoolPool__Name":                  junosSecurityNatDestinationRule__SetRuleThenDestination__NatPoolPool__Name(),
			"junos-vsrx_SecurityNatProxy__ArpInterfaceAddressToIpaddr":                                          junosSecurityNatProxy__ArpInterfaceAddressToIpaddr(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyMatchSource__Address":                                       junosSecurityPoliciesPolicyPolicyMatchSource__Address(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyMatchDestination__Address":                                  junosSecurityPoliciesPolicyPolicyMatchDestination__Address(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyMatchApplication":                                           junosSecurityPoliciesPolicyPolicyMatchApplication(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy":                               junosSecurityPoliciesPolicyPolicyThenPermitTunnelPair__Policy(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyThenLogSession__Init":                                       junosSecurityPoliciesPolicyPolicyThenLogSession__Init(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyThenCountApply__Groups":                                     junosSecurityPoliciesPolicyPolicyThenCountApply__Groups(),
			"junos-vsrx_SecurityPoliciesPolicyPolicyThenDeny":                                                   junosSecurityPoliciesPolicyPolicyThenDeny(),
			"junos-vsrx_SecurityZonesSecurity__ZoneHost__Inbound__TrafficSystem__ServicesName":                  junosSecurityZonesSecurity__ZoneHost__Inbound__TrafficSystem__ServicesName(),
			"junos-vsrx_SecurityZonesSecurity__ZoneInterfacesName":                                              junosSecurityZonesSecurity__ZoneInterfacesName(),
			"junos-vsrx_commit":        junosCommit(),
			"junos-vsrx_destroycommit": junosDestroyCommit(),
		},
		ConfigureFunc: providerConfigure,
	}
}
