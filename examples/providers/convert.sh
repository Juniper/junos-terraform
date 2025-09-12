#!/bin/bash

set -e

jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*{spine,leaf}*.xml ../evpn-vxlan-dc/dc2/*spine*.xml -j terraform-provider-junos-qfx-evpn-vxlan/trimmed_schema.json -t qfx-evpn-vxlan  -d ../terraform_files/junos-qfx-evpn-vlan -u jcluser -p 'Juniper!1'

jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*firewall*.xml ../evpn-vxlan-dc/dc2/*firewall*.xml -j terraform-provider-junos-srx-evpn-vxlan/trimmed_schema.json -t srx-evpn-vxlan  -d ../terraform_files/junos-srx-evpn-vxlan -u jcluser -p 'Juniper!1'

