bgp = {
  EVPN_iBGP = {
    neighbor_type     = "internal",
    local-address     = "10.30.100.3",
    af-evpn-signaling = 1,
    cluster           = "10.30.100.3",
    local-as          = "65200",
    multipath         = "1",
    neighbors = [
      { name = "10.30.100.4" }
    ]
  }
  IPCLOS_eBGP = {
    neighbor_type         = "external",
    mtu-discovery         = "1",
    local-as              = "65501",
    multipath-multiple-as = "1",
    bfd-multiplier        = 3,
    bfd-minimum-interval  = 1000,
    neighbors = [
      { name = "10.30.135.2", description = "EBGP peering to 10.30.135.2", peer-as = 65503 },
      { name = "10.30.136.2", description = "EBGP peering to 10.30.136.2", peer-as = 65504 },
      { name = "10.30.137.2", description = "EBGP peering to 10.30.137.2", peer-as = 65505 },
      { name = "10.30.131.2", description = "EBGP peering to 10.30.131.2", peer-as = 65506 },
      { name = "10.30.132.2", description = "EBGP peering to 10.30.132.2", peer-as = 65507 },
    ],
    vpn-apply-export = "1"

  }
}
