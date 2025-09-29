#!/bin/bash

set -e

jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*{spine,leaf}*.xml ../evpn-vxlan-dc/dc2/*spine*.xml -j terraform-provider-junos-vqfx/trimmed_schema.json -t vqfx -d ../terraform_files -u jcluser -p 'Juniper!1'
jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*firewall*.xml ../evpn-vxlan-dc/dc2/*firewall*.xml -j terraform-provider-junos-vsrx/trimmed_schema.json -t vsrx -d ../terraform_files  -u jcluser -p 'Juniper!1'

