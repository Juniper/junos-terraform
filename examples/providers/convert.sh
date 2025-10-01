#!/bin/bash

set -e

jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*{spine,leaf}*.xml ../evpn-vxlan-dc/dc2/*spine*.xml -j terraform-provider-junos-vqfx-evpn-vxlan/trimmed_schema.json -t vqfx-evpn-vxlan -d ../terraform_files -u jcluser -p 'Juniper!1'
jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*firewall*.xml ../evpn-vxlan-dc/dc2/*firewall*.xml -j terraform-provider-junos-vsrx-evpn-vxlan/trimmed_schema.json -t vsrx-evpn-vxlan -d ../terraform_files  -u jcluser -p 'Juniper!1'

