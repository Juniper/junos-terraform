package patch

import (
  "testing"
)

var TrimmedSchemaJSON = `{
  "path": "",
  "root": {
    "children": [
      {
        "children": [
          {
            "children": [
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "union",
                        "name": "device-count",
                        "path": "chassis/aggregated-devices/ethernet",
                        "type": "leaf",
                        "types": [
                          {
                            "path": "chassis/aggregated-devices/ethernet/device-count",
                            "patterns": [
                              "\u003c.*\u003e|$.*"
                            ],
                            "type": "string"
                          },
                          {
                            "path": "chassis/aggregated-devices/ethernet/device-count",
                            "type": "uint32"
                          }
                        ]
                      }
                    ],
                    "name": "ethernet",
                    "path": "chassis/aggregated-devices",
                    "type": "container"
                  }
                ],
                "name": "aggregated-devices",
                "path": "chassis",
                "type": "container"
              }
            ],
            "name": "chassis",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "lengths": [
                      {
                        "max": 127,
                        "min": 1,
                        "path": "forwarding-options/storm-control-profiles/name"
                      }
                    ],
                    "name": "name",
                    "path": "forwarding-options/storm-control-profiles",
                    "type": "leaf"
                  },
                  {
                    "name": "all",
                    "path": "forwarding-options/storm-control-profiles",
                    "type": "container"
                  }
                ],
                "key": "name",
                "name": "storm-control-profiles",
                "path": "forwarding-options",
                "type": "list"
              }
            ],
            "name": "forwarding-options",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "name",
                    "path": "interfaces/interface",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "string",
                    "name": "description",
                    "path": "interfaces/interface",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "empty",
                    "name": "vlan-tagging",
                    "path": "interfaces/interface",
                    "type": "leaf"
                  },
                  {
                    "children": [
                      {
                        "base-type": "string",
                        "leaf-type": "jt:esi",
                        "name": "identifier",
                        "path": "interfaces/interface/esi",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "empty",
                        "name": "all-active",
                        "path": "interfaces/interface/esi",
                        "type": "leaf"
                      }
                    ],
                    "name": "esi",
                    "path": "interfaces/interface",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "leaf-type": "union",
                            "name": "bundle",
                            "path": "interfaces/interface/ether-options/ieee-802.3ad",
                            "type": "leaf",
                            "types": [
                              {
                                "path": "interfaces/interface/ether-options/ieee-802.3ad/bundle",
                                "type": "string"
                              },
                              {
                                "path": "interfaces/interface/ether-options/ieee-802.3ad/bundle",
                                "patterns": [
                                  "\u003c.*\u003e|$.*"
                                ],
                                "type": "string"
                              }
                            ]
                          }
                        ],
                        "name": "ieee-802.3ad",
                        "path": "interfaces/interface/ether-options",
                        "type": "container"
                      }
                    ],
                    "name": "ether-options",
                    "path": "interfaces/interface",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "leaf-type": "empty",
                            "name": "active",
                            "path": "interfaces/interface/aggregated-ether-options/lacp",
                            "type": "leaf"
                          },
                          {
                            "enums": [
                              {
                                "id": "fast",
                                "path": "interfaces/interface/aggregated-ether-options/lacp/periodic",
                                "value": 0
                              },
                              {
                                "id": "slow",
                                "path": "interfaces/interface/aggregated-ether-options/lacp/periodic",
                                "value": 1
                              }
                            ],
                            "leaf-type": "enumeration",
                            "name": "periodic",
                            "path": "interfaces/interface/aggregated-ether-options/lacp",
                            "type": "leaf"
                          },
                          {
                            "base-type": "string",
                            "leaf-type": "jt:mac-addr",
                            "name": "system-id",
                            "path": "interfaces/interface/aggregated-ether-options/lacp",
                            "type": "leaf"
                          }
                        ],
                        "name": "lacp",
                        "path": "interfaces/interface/aggregated-ether-options",
                        "type": "container"
                      }
                    ],
                    "name": "aggregated-ether-options",
                    "path": "interfaces/interface",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "interfaces/interface/unit",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "string",
                        "name": "description",
                        "path": "interfaces/interface/unit",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "string",
                        "name": "vlan-id",
                        "path": "interfaces/interface/unit",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "children": [
                                  {
                                    "base-type": "string",
                                    "leaf-type": "jt:ipv4prefix",
                                    "name": "name",
                                    "path": "interfaces/interface/unit/family/inet/address",
                                    "type": "leaf"
                                  }
                                ],
                                "key": "name",
                                "name": "address",
                                "ordered-by": "user",
                                "path": "interfaces/interface/unit/family/inet",
                                "type": "list"
                              }
                            ],
                            "name": "inet",
                            "path": "interfaces/interface/unit/family",
                            "type": "container"
                          },
                          {
                            "children": [
                              {
                                "children": [
                                  {
                                    "leaf-type": "string",
                                    "name": "members",
                                    "ordered-by": "user",
                                    "path": "interfaces/interface/unit/family/ethernet-switching/vlan",
                                    "type": "leaf-list"
                                  }
                                ],
                                "name": "vlan",
                                "path": "interfaces/interface/unit/family/ethernet-switching",
                                "type": "container"
                              }
                            ],
                            "name": "ethernet-switching",
                            "path": "interfaces/interface/unit/family",
                            "type": "container"
                          }
                        ],
                        "name": "family",
                        "path": "interfaces/interface/unit",
                        "type": "container"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:mac-unicast",
                        "name": "mac",
                        "path": "interfaces/interface/unit",
                        "type": "leaf"
                      }
                    ],
                    "key": "name",
                    "name": "unit",
                    "path": "interfaces/interface",
                    "type": "list"
                  }
                ],
                "key": "name",
                "name": "interface",
                "path": "interfaces",
                "type": "list"
              }
            ],
            "name": "interfaces",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "name",
                    "path": "policy-options/policy-statement",
                    "type": "leaf"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "policy-options/policy-statement/term",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "enums": [
                              {
                                "id": "aggregate",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 0
                              },
                              {
                                "id": "bgp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 1
                              },
                              {
                                "id": "direct",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 2
                              },
                              {
                                "id": "dvmrp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 3
                              },
                              {
                                "id": "isis",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 4
                              },
                              {
                                "id": "esis",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 5
                              },
                              {
                                "id": "l2circuit",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 6
                              },
                              {
                                "id": "l2vpn",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 7
                              },
                              {
                                "id": "local",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 8
                              },
                              {
                                "id": "ospf",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 9
                              },
                              {
                                "id": "ospf2",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 10
                              },
                              {
                                "id": "ospf3",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 11
                              },
                              {
                                "id": "pim",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 12
                              },
                              {
                                "id": "rip",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 13
                              },
                              {
                                "id": "ripng",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 14
                              },
                              {
                                "id": "static",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 15
                              },
                              {
                                "id": "arp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 16
                              },
                              {
                                "id": "frr",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 17
                              },
                              {
                                "id": "mpls",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 18
                              },
                              {
                                "id": "ldp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 19
                              },
                              {
                                "id": "rsvp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 20
                              },
                              {
                                "id": "msdp",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 21
                              },
                              {
                                "id": "route-target",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 22
                              },
                              {
                                "id": "access",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 23
                              },
                              {
                                "id": "access-internal",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 24
                              },
                              {
                                "id": "anchor",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 25
                              },
                              {
                                "id": "bgp-static",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 26
                              },
                              {
                                "id": "vpls",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 27
                              },
                              {
                                "id": "evpn",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 28
                              },
                              {
                                "id": "spring-te",
                                "path": "policy-options/policy-statement/term/from/protocol",
                                "value": 29
                              }
                            ],
                            "leaf-type": "enumeration",
                            "name": "protocol",
                            "ordered-by": "user",
                            "path": "policy-options/policy-statement/term/from",
                            "type": "leaf-list"
                          },
                          {
                            "children": [
                              {
                                "base-type": "string",
                                "leaf-type": "jt:ipprefix",
                                "name": "address",
                                "path": "policy-options/policy-statement/term/from/route-filter",
                                "type": "leaf"
                              },
                              {
                                "leaf-type": "string",
                                "name": "exact",
                                "path": "policy-options/policy-statement/term/from/route-filter",
                                "type": "leaf"
                              },
                              {
                                "leaf-type": "string",
                                "name": "orlonger",
                                "path": "policy-options/policy-statement/term/from/route-filter",
                                "type": "leaf"
                              },
                              {
                                "leaf-type": "string",
                                "name": "prefix-length-range",
                                "path": "policy-options/policy-statement/term/from/route-filter",
                                "type": "leaf"
                              }
                            ],
                            "key": "address choice-ident choice-value",
                            "name": "route-filter",
                            "ordered-by": "user",
                            "path": "policy-options/policy-statement/term/from",
                            "type": "list"
                          }
                        ],
                        "name": "from",
                        "path": "policy-options/policy-statement/term",
                        "type": "container"
                      },
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "leaf-type": "string",
                                "name": "add",
                                "path": "policy-options/policy-statement/term/then/community",
                                "type": "leaf"
                              },
                              {
                                "leaf-type": "string",
                                "name": "community-name",
                                "path": "policy-options/policy-statement/term/then/community",
                                "type": "leaf"
                              }
                            ],
                            "key": "choice-ident choice-value community-name",
                            "name": "community",
                            "ordered-by": "user",
                            "path": "policy-options/policy-statement/term/then",
                            "type": "list"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "accept",
                            "path": "policy-options/policy-statement/term/then",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "reject",
                            "path": "policy-options/policy-statement/term/then",
                            "type": "leaf"
                          }
                        ],
                        "name": "then",
                        "path": "policy-options/policy-statement/term",
                        "type": "container"
                      }
                    ],
                    "key": "name",
                    "name": "term",
                    "ordered-by": "user",
                    "path": "policy-options/policy-statement",
                    "type": "list"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "leaf-type": "empty",
                            "name": "per-packet",
                            "path": "policy-options/policy-statement/then/load-balance",
                            "type": "leaf"
                          }
                        ],
                        "name": "load-balance",
                        "path": "policy-options/policy-statement/then",
                        "type": "container"
                      }
                    ],
                    "name": "then",
                    "path": "policy-options/policy-statement",
                    "type": "container"
                  }
                ],
                "key": "name",
                "name": "policy-statement",
                "path": "policy-options",
                "type": "list"
              },
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "name",
                    "path": "policy-options/community",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "string",
                    "name": "members",
                    "ordered-by": "user",
                    "path": "policy-options/community",
                    "type": "leaf-list"
                  }
                ],
                "key": "name",
                "name": "community",
                "path": "policy-options",
                "type": "list"
              }
            ],
            "name": "policy-options",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "enums": [
                          {
                            "id": "internal",
                            "path": "protocols/bgp/group/type",
                            "value": 0
                          },
                          {
                            "id": "external",
                            "path": "protocols/bgp/group/type",
                            "value": 1
                          }
                        ],
                        "leaf-type": "enumeration",
                        "name": "type",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "empty",
                            "name": "no-nexthop-change",
                            "path": "protocols/bgp/group/multihop",
                            "type": "leaf"
                          }
                        ],
                        "name": "multihop",
                        "path": "protocols/bgp/group",
                        "type": "container"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:ipaddr",
                        "name": "local-address",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "empty",
                        "name": "mtu-discovery",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:policy-algebra",
                        "name": "import",
                        "ordered-by": "user",
                        "path": "protocols/bgp/group",
                        "type": "leaf-list"
                      },
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "children": [
                                  {
                                    "children": [
                                      {
                                        "children": [
                                          {
                                            "leaf-type": "union",
                                            "name": "routing-uptime",
                                            "path": "protocols/bgp/group/family/evpn/signaling/delay-route-advertisements/minimum-delay",
                                            "type": "leaf",
                                            "types": [
                                              {
                                                "path": "protocols/bgp/group/family/evpn/signaling/delay-route-advertisements/minimum-delay/routing-uptime",
                                                "patterns": [
                                                  "\u003c.*\u003e|$.*"
                                                ],
                                                "type": "string"
                                              },
                                              {
                                                "path": "protocols/bgp/group/family/evpn/signaling/delay-route-advertisements/minimum-delay/routing-uptime",
                                                "ranges": [
                                                  {
                                                    "max": 36000,
                                                    "min": 1,
                                                    "path": "protocols/bgp/group/family/evpn/signaling/delay-route-advertisements/minimum-delay/routing-uptime"
                                                  }
                                                ],
                                                "type": "uint32"
                                              }
                                            ]
                                          }
                                        ],
                                        "name": "minimum-delay",
                                        "path": "protocols/bgp/group/family/evpn/signaling/delay-route-advertisements",
                                        "type": "container"
                                      }
                                    ],
                                    "name": "delay-route-advertisements",
                                    "path": "protocols/bgp/group/family/evpn/signaling",
                                    "type": "container"
                                  }
                                ],
                                "name": "signaling",
                                "path": "protocols/bgp/group/family/evpn",
                                "type": "container"
                              }
                            ],
                            "name": "evpn",
                            "path": "protocols/bgp/group/family",
                            "type": "container"
                          }
                        ],
                        "name": "family",
                        "path": "protocols/bgp/group",
                        "type": "container"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:policy-algebra",
                        "name": "export",
                        "ordered-by": "user",
                        "path": "protocols/bgp/group",
                        "type": "leaf-list"
                      },
                      {
                        "leaf-type": "empty",
                        "name": "vpn-apply-export",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:areaid",
                        "name": "cluster",
                        "path": "protocols/bgp/group",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "string",
                            "name": "as-number",
                            "path": "protocols/bgp/group/local-as",
                            "type": "leaf"
                          }
                        ],
                        "name": "local-as",
                        "path": "protocols/bgp/group",
                        "type": "container"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "empty",
                            "name": "multiple-as",
                            "path": "protocols/bgp/group/multipath",
                            "type": "leaf"
                          }
                        ],
                        "name": "multipath",
                        "path": "protocols/bgp/group",
                        "type": "container"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "union",
                            "name": "minimum-interval",
                            "path": "protocols/bgp/group/bfd-liveness-detection",
                            "type": "leaf",
                            "types": [
                              {
                                "path": "protocols/bgp/group/bfd-liveness-detection/minimum-interval",
                                "patterns": [
                                  "\u003c.*\u003e|$.*"
                                ],
                                "type": "string"
                              },
                              {
                                "path": "protocols/bgp/group/bfd-liveness-detection/minimum-interval",
                                "ranges": [
                                  {
                                    "max": 255000,
                                    "min": 1,
                                    "path": "protocols/bgp/group/bfd-liveness-detection/minimum-interval"
                                  }
                                ],
                                "type": "uint32"
                              }
                            ],
                            "units": "milliseconds"
                          },
                          {
                            "default": "3",
                            "leaf-type": "union",
                            "name": "multiplier",
                            "path": "protocols/bgp/group/bfd-liveness-detection",
                            "type": "leaf",
                            "types": [
                              {
                                "path": "protocols/bgp/group/bfd-liveness-detection/multiplier",
                                "patterns": [
                                  "\u003c.*\u003e|$.*"
                                ],
                                "type": "string"
                              },
                              {
                                "path": "protocols/bgp/group/bfd-liveness-detection/multiplier",
                                "ranges": [
                                  {
                                    "max": 255,
                                    "min": 1,
                                    "path": "protocols/bgp/group/bfd-liveness-detection/multiplier"
                                  }
                                ],
                                "type": "uint32"
                              }
                            ]
                          }
                        ],
                        "name": "bfd-liveness-detection",
                        "path": "protocols/bgp/group",
                        "type": "container"
                      },
                      {
                        "base-type": "string",
                        "leaf-type": "jt:ipprefix",
                        "name": "allow",
                        "ordered-by": "user",
                        "path": "protocols/bgp/group",
                        "type": "leaf-list"
                      },
                      {
                        "children": [
                          {
                            "base-type": "string",
                            "leaf-type": "jt:ipaddr",
                            "name": "name",
                            "path": "protocols/bgp/group/neighbor",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "string",
                            "lengths": [
                              {
                                "max": 255,
                                "min": 1,
                                "path": "protocols/bgp/group/neighbor/description"
                              }
                            ],
                            "name": "description",
                            "path": "protocols/bgp/group/neighbor",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "string",
                            "name": "peer-as",
                            "path": "protocols/bgp/group/neighbor",
                            "type": "leaf"
                          }
                        ],
                        "key": "name",
                        "name": "neighbor",
                        "ordered-by": "user",
                        "path": "protocols/bgp/group",
                        "type": "list"
                      }
                    ],
                    "key": "name",
                    "name": "group",
                    "ordered-by": "user",
                    "path": "protocols/bgp",
                    "type": "list"
                  }
                ],
                "name": "bgp",
                "path": "protocols",
                "type": "container"
              },
              {
                "children": [
                  {
                    "default": "mpls",
                    "enums": [
                      {
                        "id": "mpls",
                        "path": "protocols/evpn/encapsulation",
                        "value": 0
                      },
                      {
                        "id": "vxlan",
                        "path": "protocols/evpn/encapsulation",
                        "value": 1
                      }
                    ],
                    "leaf-type": "enumeration",
                    "name": "encapsulation",
                    "path": "protocols/evpn",
                    "type": "leaf"
                  },
                  {
                    "default": "ingress-replication",
                    "enums": [
                      {
                        "id": "ingress-replication",
                        "path": "protocols/evpn/multicast-mode",
                        "value": 0
                      }
                    ],
                    "leaf-type": "enumeration",
                    "name": "multicast-mode",
                    "path": "protocols/evpn",
                    "type": "leaf"
                  },
                  {
                    "enums": [
                      {
                        "id": "advertise",
                        "path": "protocols/evpn/default-gateway",
                        "value": 0
                      },
                      {
                        "id": "no-gateway-community",
                        "path": "protocols/evpn/default-gateway",
                        "value": 1
                      },
                      {
                        "id": "do-not-advertise",
                        "path": "protocols/evpn/default-gateway",
                        "value": 2
                      }
                    ],
                    "leaf-type": "enumeration",
                    "name": "default-gateway",
                    "path": "protocols/evpn",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "string",
                    "name": "extended-vni-list",
                    "path": "protocols/evpn",
                    "type": "leaf-list"
                  },
                  {
                    "leaf-type": "empty",
                    "name": "no-core-isolation",
                    "path": "protocols/evpn",
                    "type": "leaf"
                  }
                ],
                "name": "evpn",
                "path": "protocols",
                "type": "container"
              },
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "protocols/lldp/interface",
                        "type": "leaf"
                      }
                    ],
                    "key": "name",
                    "name": "interface",
                    "ordered-by": "user",
                    "path": "protocols/lldp",
                    "type": "list"
                  }
                ],
                "name": "lldp",
                "path": "protocols",
                "type": "container"
              },
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "protocols/igmp-snooping/vlan",
                        "type": "leaf"
                      }
                    ],
                    "key": "name",
                    "name": "vlan",
                    "ordered-by": "user",
                    "path": "protocols/igmp-snooping",
                    "type": "list"
                  }
                ],
                "name": "igmp-snooping",
                "path": "protocols",
                "type": "container"
              }
            ],
            "name": "protocols",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "name",
                    "path": "routing-instances/instance",
                    "type": "leaf"
                  },
                  {
                    "enums": [
                      {
                        "id": "forwarding",
                        "path": "routing-instances/instance/instance-type",
                        "value": 0
                      },
                      {
                        "id": "vrf",
                        "path": "routing-instances/instance/instance-type",
                        "value": 1
                      },
                      {
                        "id": "no-forwarding",
                        "path": "routing-instances/instance/instance-type",
                        "value": 2
                      },
                      {
                        "id": "l2vpn",
                        "path": "routing-instances/instance/instance-type",
                        "value": 3
                      },
                      {
                        "id": "vpls",
                        "path": "routing-instances/instance/instance-type",
                        "value": 4
                      },
                      {
                        "id": "virtual-switch",
                        "path": "routing-instances/instance/instance-type",
                        "value": 5
                      },
                      {
                        "id": "l2backhaul-vpn",
                        "path": "routing-instances/instance/instance-type",
                        "value": 6
                      },
                      {
                        "id": "virtual-router",
                        "path": "routing-instances/instance/instance-type",
                        "value": 7
                      },
                      {
                        "id": "layer2-control",
                        "path": "routing-instances/instance/instance-type",
                        "value": 8
                      },
                      {
                        "id": "mpls-internet-multicast",
                        "path": "routing-instances/instance/instance-type",
                        "value": 9
                      },
                      {
                        "id": "evpn",
                        "path": "routing-instances/instance/instance-type",
                        "value": 10
                      },
                      {
                        "id": "mpls-forwarding",
                        "path": "routing-instances/instance/instance-type",
                        "value": 11
                      },
                      {
                        "id": "evpn-vpws",
                        "path": "routing-instances/instance/instance-type",
                        "value": 12
                      }
                    ],
                    "leaf-type": "enumeration",
                    "name": "instance-type",
                    "path": "routing-instances/instance",
                    "type": "leaf"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "routing-instances/instance/interface",
                        "type": "leaf"
                      }
                    ],
                    "key": "name",
                    "name": "interface",
                    "path": "routing-instances/instance",
                    "type": "list"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "rd-type",
                        "path": "routing-instances/instance/route-distinguisher",
                        "type": "leaf"
                      }
                    ],
                    "name": "route-distinguisher",
                    "path": "routing-instances/instance",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "community",
                        "path": "routing-instances/instance/vrf-target",
                        "type": "leaf"
                      }
                    ],
                    "name": "vrf-target",
                    "path": "routing-instances/instance",
                    "type": "container"
                  },
                  {
                    "name": "vrf-table-label",
                    "path": "routing-instances/instance",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "name": "auto-export",
                        "path": "routing-instances/instance/routing-options",
                        "type": "container"
                      }
                    ],
                    "name": "routing-options",
                    "path": "routing-instances/instance",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "base-type": "string",
                            "leaf-type": "jt:policy-algebra",
                            "name": "export",
                            "ordered-by": "user",
                            "path": "routing-instances/instance/protocols/ospf",
                            "type": "leaf-list"
                          },
                          {
                            "children": [
                              {
                                "base-type": "string",
                                "leaf-type": "jt:areaid",
                                "name": "name",
                                "path": "routing-instances/instance/protocols/ospf/area",
                                "type": "leaf"
                              },
                              {
                                "children": [
                                  {
                                    "leaf-type": "union",
                                    "name": "name",
                                    "path": "routing-instances/instance/protocols/ospf/area/interface",
                                    "type": "leaf",
                                    "types": [
                                      {
                                        "path": "routing-instances/instance/protocols/ospf/area/interface/name",
                                        "type": "string"
                                      },
                                      {
                                        "path": "routing-instances/instance/protocols/ospf/area/interface/name",
                                        "patterns": [
                                          "\u003c.*\u003e|$.*"
                                        ],
                                        "type": "string"
                                      }
                                    ]
                                  },
                                  {
                                    "leaf-type": "union",
                                    "name": "metric",
                                    "path": "routing-instances/instance/protocols/ospf/area/interface",
                                    "type": "leaf",
                                    "types": [
                                      {
                                        "path": "routing-instances/instance/protocols/ospf/area/interface/metric",
                                        "patterns": [
                                          "\u003c.*\u003e|$.*"
                                        ],
                                        "type": "string"
                                      },
                                      {
                                        "path": "routing-instances/instance/protocols/ospf/area/interface/metric",
                                        "ranges": [
                                          {
                                            "max": 65535,
                                            "min": 1,
                                            "path": "routing-instances/instance/protocols/ospf/area/interface/metric"
                                          }
                                        ],
                                        "type": "uint16"
                                      }
                                    ]
                                  }
                                ],
                                "key": "name",
                                "name": "interface",
                                "ordered-by": "user",
                                "path": "routing-instances/instance/protocols/ospf/area",
                                "type": "list"
                              }
                            ],
                            "key": "name",
                            "name": "area",
                            "ordered-by": "user",
                            "path": "routing-instances/instance/protocols/ospf",
                            "type": "list"
                          }
                        ],
                        "name": "ospf",
                        "path": "routing-instances/instance/protocols",
                        "type": "container"
                      },
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "enums": [
                                  {
                                    "id": "gateway-address",
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/advertise",
                                    "value": 0
                                  },
                                  {
                                    "id": "direct-nexthop",
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/advertise",
                                    "value": 1
                                  }
                                ],
                                "leaf-type": "enumeration",
                                "name": "advertise",
                                "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes",
                                "type": "leaf"
                              },
                              {
                                "enums": [
                                  {
                                    "id": "mpls",
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/encapsulation",
                                    "value": 0
                                  },
                                  {
                                    "id": "vxlan",
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/encapsulation",
                                    "value": 1
                                  }
                                ],
                                "leaf-type": "enumeration",
                                "name": "encapsulation",
                                "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes",
                                "type": "leaf"
                              },
                              {
                                "leaf-type": "union",
                                "name": "vni",
                                "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes",
                                "type": "leaf",
                                "types": [
                                  {
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/vni",
                                    "patterns": [
                                      "\u003c.*\u003e|$.*"
                                    ],
                                    "type": "string"
                                  },
                                  {
                                    "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/vni",
                                    "ranges": [
                                      {
                                        "max": 16777214,
                                        "min": 1,
                                        "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes/vni"
                                      }
                                    ],
                                    "type": "uint32"
                                  }
                                ]
                              },
                              {
                                "base-type": "string",
                                "leaf-type": "jt:policy-algebra",
                                "name": "export",
                                "ordered-by": "user",
                                "path": "routing-instances/instance/protocols/evpn/ip-prefix-routes",
                                "type": "leaf-list"
                              }
                            ],
                            "name": "ip-prefix-routes",
                            "path": "routing-instances/instance/protocols/evpn",
                            "type": "container"
                          }
                        ],
                        "name": "evpn",
                        "path": "routing-instances/instance/protocols",
                        "type": "container"
                      }
                    ],
                    "name": "protocols",
                    "path": "routing-instances/instance",
                    "type": "container"
                  }
                ],
                "key": "name",
                "name": "instance",
                "path": "routing-instances",
                "type": "list"
              }
            ],
            "name": "routing-instances",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "children": [
                      {
                        "base-type": "string",
                        "leaf-type": "jt:ipprefix",
                        "name": "name",
                        "path": "routing-options/static/route",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "union",
                        "name": "next-hop",
                        "ordered-by": "user",
                        "path": "routing-options/static/route",
                        "type": "leaf-list",
                        "types": [
                          {
                            "path": "routing-options/static/route/next-hop",
                            "type": "string"
                          },
                          {
                            "path": "routing-options/static/route/next-hop",
                            "patterns": [
                              "\u003c.*\u003e|$.*"
                            ],
                            "type": "string"
                          }
                        ]
                      }
                    ],
                    "key": "name",
                    "name": "route",
                    "ordered-by": "user",
                    "path": "routing-options/static",
                    "type": "list"
                  }
                ],
                "name": "static",
                "path": "routing-options",
                "type": "container"
              },
              {
                "base-type": "string",
                "leaf-type": "jt:ipv4addr",
                "name": "router-id",
                "path": "routing-options",
                "type": "leaf"
              },
              {
                "children": [
                  {
                    "base-type": "string",
                    "leaf-type": "jt:policy-algebra",
                    "name": "export",
                    "ordered-by": "user",
                    "path": "routing-options/forwarding-table",
                    "type": "leaf-list"
                  },
                  {
                    "leaf-type": "empty",
                    "name": "ecmp-fast-reroute",
                    "path": "routing-options/forwarding-table",
                    "type": "leaf"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "leaf-type": "empty",
                            "name": "evpn",
                            "path": "routing-options/forwarding-table/chained-composite-next-hop/ingress",
                            "type": "leaf"
                          }
                        ],
                        "name": "ingress",
                        "path": "routing-options/forwarding-table/chained-composite-next-hop",
                        "type": "container"
                      }
                    ],
                    "name": "chained-composite-next-hop",
                    "path": "routing-options/forwarding-table",
                    "type": "container"
                  }
                ],
                "name": "forwarding-table",
                "path": "routing-options",
                "type": "container"
              }
            ],
            "name": "routing-options",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "leaf-type": "string",
                "name": "location",
                "path": "snmp",
                "type": "leaf"
              },
              {
                "leaf-type": "string",
                "name": "contact",
                "path": "snmp",
                "type": "leaf"
              },
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "name",
                    "path": "snmp/community",
                    "type": "leaf"
                  },
                  {
                    "enums": [
                      {
                        "id": "read-only",
                        "path": "snmp/community/authorization",
                        "value": 0
                      },
                      {
                        "id": "read-write",
                        "path": "snmp/community/authorization",
                        "value": 1
                      }
                    ],
                    "leaf-type": "enumeration",
                    "name": "authorization",
                    "path": "snmp/community",
                    "type": "leaf"
                  }
                ],
                "key": "name",
                "name": "community",
                "ordered-by": "user",
                "path": "snmp",
                "type": "list"
              }
            ],
            "name": "snmp",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "union",
                    "name": "interface-name",
                    "path": "switch-options/vtep-source-interface",
                    "type": "leaf",
                    "types": [
                      {
                        "path": "switch-options/vtep-source-interface/interface-name",
                        "type": "string"
                      },
                      {
                        "path": "switch-options/vtep-source-interface/interface-name",
                        "patterns": [
                          "\u003c.*\u003e|$.*"
                        ],
                        "type": "string"
                      }
                    ]
                  }
                ],
                "name": "vtep-source-interface",
                "path": "switch-options",
                "type": "container"
              },
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "rd-type",
                    "path": "switch-options/route-distinguisher",
                    "type": "leaf"
                  }
                ],
                "name": "route-distinguisher",
                "path": "switch-options",
                "type": "container"
              },
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "name": "community",
                    "path": "switch-options/vrf-target",
                    "type": "leaf"
                  },
                  {
                    "name": "auto",
                    "path": "switch-options/vrf-target",
                    "type": "container"
                  }
                ],
                "name": "vrf-target",
                "path": "switch-options",
                "type": "container"
              }
            ],
            "name": "switch-options",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "system/login/user",
                        "type": "leaf"
                      },
                      {
                        "leaf-type": "union",
                        "name": "uid",
                        "path": "system/login/user",
                        "type": "leaf",
                        "types": [
                          {
                            "path": "system/login/user/uid",
                            "patterns": [
                              "\u003c.*\u003e|$.*"
                            ],
                            "type": "string"
                          },
                          {
                            "path": "system/login/user/uid",
                            "ranges": [
                              {
                                "max": 64000,
                                "min": 100,
                                "path": "system/login/user/uid"
                              }
                            ],
                            "type": "uint32"
                          }
                        ]
                      },
                      {
                        "leaf-type": "string",
                        "name": "class",
                        "path": "system/login/user",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "string",
                            "lengths": [
                              {
                                "max": 128,
                                "min": 1,
                                "path": "system/login/user/authentication/encrypted-password"
                              }
                            ],
                            "name": "encrypted-password",
                            "path": "system/login/user/authentication",
                            "type": "leaf"
                          }
                        ],
                        "name": "authentication",
                        "path": "system/login/user",
                        "type": "container"
                      }
                    ],
                    "key": "name",
                    "name": "user",
                    "path": "system/login",
                    "type": "list"
                  },
                  {
                    "leaf-type": "string",
                    "lengths": [
                      {
                        "max": 2048,
                        "min": 1,
                        "path": "system/login/message"
                      }
                    ],
                    "name": "message",
                    "path": "system/login",
                    "type": "leaf"
                  }
                ],
                "name": "login",
                "path": "system",
                "type": "container"
              },
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "lengths": [
                      {
                        "max": 128,
                        "min": 1,
                        "path": "system/root-authentication/encrypted-password"
                      }
                    ],
                    "name": "encrypted-password",
                    "path": "system/root-authentication",
                    "type": "leaf"
                  }
                ],
                "name": "root-authentication",
                "path": "system",
                "type": "container"
              },
              {
                "leaf-type": "string",
                "lengths": [
                  {
                    "max": 255,
                    "min": 1,
                    "path": "system/host-name"
                  }
                ],
                "name": "host-name",
                "path": "system",
                "type": "leaf"
              },
              {
                "children": [
                  {
                    "children": [
                      {
                        "enums": [
                          {
                            "id": "allow",
                            "path": "system/services/ssh/root-login",
                            "value": 0
                          },
                          {
                            "id": "deny",
                            "path": "system/services/ssh/root-login",
                            "value": 1
                          },
                          {
                            "id": "deny-password",
                            "path": "system/services/ssh/root-login",
                            "value": 2
                          }
                        ],
                        "leaf-type": "enumeration",
                        "name": "root-login",
                        "path": "system/services/ssh",
                        "type": "leaf"
                      }
                    ],
                    "name": "ssh",
                    "path": "system/services",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "default": "5",
                                "leaf-type": "union",
                                "name": "max-connections",
                                "path": "system/services/extension-service/request-response/grpc",
                                "type": "leaf",
                                "types": [
                                  {
                                    "path": "system/services/extension-service/request-response/grpc/max-connections",
                                    "patterns": [
                                      "\u003c.*\u003e|$.*"
                                    ],
                                    "type": "string"
                                  },
                                  {
                                    "path": "system/services/extension-service/request-response/grpc/max-connections",
                                    "ranges": [
                                      {
                                        "max": 30,
                                        "min": 1,
                                        "path": "system/services/extension-service/request-response/grpc/max-connections"
                                      }
                                    ],
                                    "type": "uint32"
                                  }
                                ]
                              }
                            ],
                            "name": "grpc",
                            "path": "system/services/extension-service/request-response",
                            "type": "container"
                          }
                        ],
                        "name": "request-response",
                        "path": "system/services/extension-service",
                        "type": "container"
                      },
                      {
                        "children": [
                          {
                            "children": [
                              {
                                "base-type": "string",
                                "leaf-type": "jt:ipprefix-optional",
                                "name": "address",
                                "ordered-by": "user",
                                "path": "system/services/extension-service/notification/allow-clients",
                                "type": "leaf-list"
                              }
                            ],
                            "name": "allow-clients",
                            "path": "system/services/extension-service/notification",
                            "type": "container"
                          }
                        ],
                        "name": "notification",
                        "path": "system/services/extension-service",
                        "type": "container"
                      }
                    ],
                    "name": "extension-service",
                    "path": "system/services",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "name": "ssh",
                        "path": "system/services/netconf",
                        "type": "container"
                      }
                    ],
                    "name": "netconf",
                    "path": "system/services",
                    "type": "container"
                  },
                  {
                    "children": [
                      {
                        "children": [
                          {
                            "default": "3000",
                            "leaf-type": "union",
                            "name": "port",
                            "path": "system/services/rest/http",
                            "type": "leaf",
                            "types": [
                              {
                                "path": "system/services/rest/http/port",
                                "patterns": [
                                  "\u003c.*\u003e|$.*"
                                ],
                                "type": "string"
                              },
                              {
                                "path": "system/services/rest/http/port",
                                "ranges": [
                                  {
                                    "max": 65535,
                                    "min": 1024,
                                    "path": "system/services/rest/http/port"
                                  }
                                ],
                                "type": "uint32"
                              }
                            ]
                          }
                        ],
                        "name": "http",
                        "path": "system/services/rest",
                        "type": "container"
                      },
                      {
                        "leaf-type": "empty",
                        "name": "enable-explorer",
                        "path": "system/services/rest",
                        "type": "leaf"
                      }
                    ],
                    "name": "rest",
                    "path": "system/services",
                    "type": "container"
                  }
                ],
                "name": "services",
                "path": "system",
                "type": "container"
              },
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "system/syslog/user",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "enums": [
                              {
                                "id": "any",
                                "path": "system/syslog/user/contents/name",
                                "value": 0
                              },
                              {
                                "id": "authorization",
                                "path": "system/syslog/user/contents/name",
                                "value": 1
                              },
                              {
                                "id": "daemon",
                                "path": "system/syslog/user/contents/name",
                                "value": 2
                              },
                              {
                                "id": "ftp",
                                "path": "system/syslog/user/contents/name",
                                "value": 3
                              },
                              {
                                "id": "ntp",
                                "path": "system/syslog/user/contents/name",
                                "value": 4
                              },
                              {
                                "id": "security",
                                "path": "system/syslog/user/contents/name",
                                "value": 5
                              },
                              {
                                "id": "kernel",
                                "path": "system/syslog/user/contents/name",
                                "value": 6
                              },
                              {
                                "id": "user",
                                "path": "system/syslog/user/contents/name",
                                "value": 7
                              },
                              {
                                "id": "dfc",
                                "path": "system/syslog/user/contents/name",
                                "value": 8
                              },
                              {
                                "id": "external",
                                "path": "system/syslog/user/contents/name",
                                "value": 9
                              },
                              {
                                "id": "firewall",
                                "path": "system/syslog/user/contents/name",
                                "value": 10
                              },
                              {
                                "id": "pfe",
                                "path": "system/syslog/user/contents/name",
                                "value": 11
                              },
                              {
                                "id": "conflict-log",
                                "path": "system/syslog/user/contents/name",
                                "value": 12
                              },
                              {
                                "id": "change-log",
                                "path": "system/syslog/user/contents/name",
                                "value": 13
                              },
                              {
                                "id": "interactive-commands",
                                "path": "system/syslog/user/contents/name",
                                "value": 14
                              }
                            ],
                            "leaf-type": "enumeration",
                            "name": "name",
                            "path": "system/syslog/user/contents",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "emergency",
                            "path": "system/syslog/user/contents",
                            "type": "leaf"
                          }
                        ],
                        "key": "name",
                        "name": "contents",
                        "path": "system/syslog/user",
                        "type": "list"
                      }
                    ],
                    "key": "name",
                    "name": "user",
                    "ordered-by": "user",
                    "path": "system/syslog",
                    "type": "list"
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "lengths": [
                          {
                            "max": 1024,
                            "min": 1,
                            "path": "system/syslog/file/name"
                          }
                        ],
                        "name": "name",
                        "path": "system/syslog/file",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "enums": [
                              {
                                "id": "any",
                                "path": "system/syslog/file/contents/name",
                                "value": 0
                              },
                              {
                                "id": "authorization",
                                "path": "system/syslog/file/contents/name",
                                "value": 1
                              },
                              {
                                "id": "daemon",
                                "path": "system/syslog/file/contents/name",
                                "value": 2
                              },
                              {
                                "id": "ftp",
                                "path": "system/syslog/file/contents/name",
                                "value": 3
                              },
                              {
                                "id": "ntp",
                                "path": "system/syslog/file/contents/name",
                                "value": 4
                              },
                              {
                                "id": "security",
                                "path": "system/syslog/file/contents/name",
                                "value": 5
                              },
                              {
                                "id": "kernel",
                                "path": "system/syslog/file/contents/name",
                                "value": 6
                              },
                              {
                                "id": "user",
                                "path": "system/syslog/file/contents/name",
                                "value": 7
                              },
                              {
                                "id": "dfc",
                                "path": "system/syslog/file/contents/name",
                                "value": 8
                              },
                              {
                                "id": "external",
                                "path": "system/syslog/file/contents/name",
                                "value": 9
                              },
                              {
                                "id": "firewall",
                                "path": "system/syslog/file/contents/name",
                                "value": 10
                              },
                              {
                                "id": "pfe",
                                "path": "system/syslog/file/contents/name",
                                "value": 11
                              },
                              {
                                "id": "conflict-log",
                                "path": "system/syslog/file/contents/name",
                                "value": 12
                              },
                              {
                                "id": "change-log",
                                "path": "system/syslog/file/contents/name",
                                "value": 13
                              },
                              {
                                "id": "interactive-commands",
                                "path": "system/syslog/file/contents/name",
                                "value": 14
                              }
                            ],
                            "leaf-type": "enumeration",
                            "name": "name",
                            "path": "system/syslog/file/contents",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "any",
                            "path": "system/syslog/file/contents",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "notice",
                            "path": "system/syslog/file/contents",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "empty",
                            "name": "info",
                            "path": "system/syslog/file/contents",
                            "type": "leaf"
                          }
                        ],
                        "key": "name",
                        "name": "contents",
                        "path": "system/syslog/file",
                        "type": "list"
                      }
                    ],
                    "key": "name",
                    "name": "file",
                    "ordered-by": "user",
                    "path": "system/syslog",
                    "type": "list"
                  }
                ],
                "name": "syslog",
                "path": "system",
                "type": "container"
              },
              {
                "children": [
                  {
                    "children": [
                      {
                        "leaf-type": "string",
                        "name": "name",
                        "path": "system/extensions/providers",
                        "type": "leaf"
                      },
                      {
                        "children": [
                          {
                            "leaf-type": "string",
                            "name": "name",
                            "path": "system/extensions/providers/license-type",
                            "type": "leaf"
                          },
                          {
                            "leaf-type": "string",
                            "name": "deployment-scope",
                            "ordered-by": "user",
                            "path": "system/extensions/providers/license-type",
                            "type": "leaf-list"
                          }
                        ],
                        "key": "name",
                        "name": "license-type",
                        "ordered-by": "user",
                        "path": "system/extensions/providers",
                        "type": "list"
                      }
                    ],
                    "key": "name",
                    "name": "providers",
                    "ordered-by": "user",
                    "path": "system/extensions",
                    "type": "list"
                  }
                ],
                "name": "extensions",
                "path": "system",
                "type": "container"
              }
            ],
            "name": "system",
            "path": "",
            "type": "container"
          },
          {
            "children": [
              {
                "children": [
                  {
                    "leaf-type": "string",
                    "lengths": [
                      {
                        "max": 64,
                        "min": 2,
                        "path": "vlans/vlan/name"
                      }
                    ],
                    "name": "name",
                    "path": "vlans/vlan",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "string",
                    "name": "vlan-id",
                    "path": "vlans/vlan",
                    "type": "leaf"
                  },
                  {
                    "leaf-type": "union",
                    "name": "l3-interface",
                    "path": "vlans/vlan",
                    "type": "leaf",
                    "types": [
                      {
                        "path": "vlans/vlan/l3-interface",
                        "type": "string"
                      },
                      {
                        "path": "vlans/vlan/l3-interface",
                        "patterns": [
                          "\u003c.*\u003e|$.*"
                        ],
                        "type": "string"
                      }
                    ]
                  },
                  {
                    "children": [
                      {
                        "leaf-type": "union",
                        "name": "vni",
                        "path": "vlans/vlan/vxlan",
                        "type": "leaf",
                        "types": [
                          {
                            "path": "vlans/vlan/vxlan/vni",
                            "patterns": [
                              "\u003c.*\u003e|$.*"
                            ],
                            "type": "string"
                          },
                          {
                            "path": "vlans/vlan/vxlan/vni",
                            "ranges": [
                              {
                                "max": 16777214,
                                "min": 0,
                                "path": "vlans/vlan/vxlan/vni"
                              }
                            ],
                            "type": "uint32"
                          }
                        ]
                      }
                    ],
                    "name": "vxlan",
                    "path": "vlans/vlan",
                    "type": "container"
                  }
                ],
                "key": "name",
                "name": "vlan",
                "path": "vlans",
                "type": "list"
              }
            ],
            "name": "vlans",
            "path": "",
            "type": "container"
          }
        ],
        "config": "true",
        "name": "configuration",
        "path": "",
        "type": "container"
      }
    ],
    "name": "root",
    "path": "",
    "type": "container"
  }
}`

var idx map[string]*NodeInfo

func TestMain(m *testing.M) {
	// setup 
	var err error
	idx, err = UnmarshalTrimmedSchemaIndex(TrimmedSchemaJSON)
	if err != nil {
		panic(err)
	}
}

func TestCreateDiffPatch_ReplaceHostName_DeleteThenCreate(t *testing.T) {
	name := "base-config"

	editLeaf := map[string]Change{
		`configuration/groups[name="base-config"]/system/host-name`: {
			Op:   "replace",
			OldV: "dc1-leaf1",
			NewV: "dc1-leaf1-test",
		},
	}

	correctDiff := `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <groups>
    <name>base-config</name>
    <system>
      <host-name operation="delete"></host-name>
      <host-name operation="create">dc1-leaf1-test</host-name>
    </system>
  </groups>
</configuration>`

	diff, err := CreateDiffPatch(editLeaf, name)
	if err != nil {
		t.Fatalf("CreateDiffPatch returned error: %v", err)
	}

	if diff != correctDiff {
		t.Fatalf("diff mismatch\n--- got ---\n%s\n--- want ---\n%s\n", diff, correctDiff)
	}
}