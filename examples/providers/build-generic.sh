#!/bin/bash
# Build generic providers using --generic flag.
# Mirrors examples/providers/build.sh but uses schema-driven generic provider.
#
# Usage: cd examples/providers && bash build-generic.sh

set -e

echo "Building generic QFX provider..."
jtaf-yang2go --generic \
  -p ../yang/18.2/18.2R3/common \
  ../yang/18.2/18.2R3/junos-qfx/conf/*.yang \
  -x ../evpn-vxlan-dc/dc1/dc1-*leaf* \
     ../evpn-vxlan-dc/dc1/dc1-*spine* \
     ../evpn-vxlan-dc/dc2/dc2-*spine* \
  -t vqfx-evpn-vxlan

echo ""
echo "Building generic SRX provider..."
jtaf-yang2go --generic \
  -p ../yang/18.2/18.2R3/common \
  ../yang/18.2/18.2R3/junos-es/conf/*.yang \
  -x ../evpn-vxlan-dc/dc1/dc1-*firewall* \
     ../evpn-vxlan-dc/dc2/dc2-*firewall* \
  -t vsrx-evpn-vxlan

echo ""
echo "Compiling QFX provider..."
cd terraform-provider-junos-vqfx-evpn-vxlan
go build .
echo "  Binary: $(du -sh terraform-provider-junos-vqfx-evpn-vxlan | cut -f1)"
cd ..

echo ""
echo "Compiling SRX provider..."
cd terraform-provider-junos-vsrx-evpn-vxlan
go build .
echo "  Binary: $(du -sh terraform-provider-junos-vsrx-evpn-vxlan | cut -f1)"
cd ..

echo ""
echo "Installing providers to \$GOPATH/bin..."
cd terraform-provider-junos-vqfx-evpn-vxlan && go install . && cd ..
cd terraform-provider-junos-vsrx-evpn-vxlan && go install . && cd ..

echo ""
echo "Done. Providers installed to $(go env GOPATH)/bin/"
echo ""
echo "Next steps:"
echo "  1. Run: bash convert.sh           (generate .tf test files)"
echo "  2. Set up ~/.terraformrc          (see examples/example-terraformrc)"
echo "  3. cd ../terraform_files"
echo "  4. terraform validate"
echo "  5. terraform plan"
echo "  6. terraform apply -auto-approve"
