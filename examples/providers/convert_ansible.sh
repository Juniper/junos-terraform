#!/bin/bash

set -e

CDIR=`pwd`
cd ../../
pip install .
cd "$CDIR"


jtaf-xml2yaml -x ../ansible/*interfaces.xml -j ansible-provider-junos-qfx-ansible-test/trimmed_schema.json -d ansible_files
#jtaf-xml2yaml -x ../evpn-vxlan-dc/dc1/*firewall*.xml ../evpn-vxlan-dc/dc2/*firewall*.xml -j terraform-provider-junos-vsrx-evpn-vxlan/trimmed_schema.json -t vsrx-evpn-vxlan -d ../terraform_files  -u jcluser -p 'Juniper!1'

