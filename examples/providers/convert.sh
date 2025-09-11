#!/bin/bash

set -e

jtaf-xml2tf -x ../evpn-vxlan-dc/dc1/*{spine,leaf}*.xml -j terraform-provider-junos-qfx-evpn-vxlan/trimmed_schema.json -t qfx -d ../terraform_files -u jcluser -p 'Juniper!1'
