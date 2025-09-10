#!/bin/bash

set -e

jtaf-xml2tf -x examples/evpn-vxlan-dc/dc1/*{spine,leaf}*.xml -j examples/providers/terraform-provider-junos-qfx-evpn-vxlan/trimmed_schema.json -t qfx -d examples/terraform_files -u jcluser -p 'Junipergit -help'
