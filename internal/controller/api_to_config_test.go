// SPDX-License-Identifier:Apache-2.0

package controller

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1beta1 "github.com/metallb/frrk8s/api/v1beta1"
	"github.com/metallb/frrk8s/internal/frr"
	"github.com/metallb/frrk8s/internal/ipfamily"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		name     string
		fromK8s  []v1beta1.FRRConfiguration
		expected *frr.Config
		err      error
	}{

		{
			name: "Single Router and Neighbor",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65001,
									ID:  "192.0.2.1",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65002,
											Address: "192.0.2.2",
											Port:    179,
										},
									},
									VRF:      "",
									Prefixes: []string{"192.0.2.0/24"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65001,
						RouterID: "192.0.2.1",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65002@192.0.2.2",
								ASN:      65002,
								Addr:     "192.0.2.2",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						VRF:          "",
						IPV4Prefixes: []string{"192.0.2.0/24"},
						IPV6Prefixes: []string{},
					},
				},
			},
			err: nil,
		},
		{
			name: "Multiple Routers and Neighbors",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65010,
									ID:  "192.0.2.5",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65011,
											Address: "192.0.2.6",
											Port:    179,
										},
										{
											ASN:     65012,
											Address: "192.0.2.7",
											Port:    179,
										},
									},
									VRF:      "",
									Prefixes: []string{"192.0.2.0/24"},
								},
								{
									ASN: 65013,
									ID:  "2001:db8::3",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65014,
											Address: "2001:db8::4",
											Port:    179,
										},
									},
									VRF:      "vrf2",
									Prefixes: []string{"2001:db8::/64"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65010,
						RouterID: "192.0.2.5",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65011@192.0.2.6",
								ASN:      65011,
								Addr:     "192.0.2.6",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65012@192.0.2.7",
								ASN:      65012,
								Addr:     "192.0.2.7",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						VRF:          "",
						IPV4Prefixes: []string{"192.0.2.0/24"},
						IPV6Prefixes: []string{},
					},
					{
						MyASN:    65013,
						RouterID: "2001:db8::3",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv6,
								Name:     "65014@2001:db8::4",
								ASN:      65014,
								Addr:     "2001:db8::4",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						VRF:          "vrf2",
						IPV4Prefixes: []string{},
						IPV6Prefixes: []string{"2001:db8::/64"},
					},
				},
			},
			err: nil,
		},
		{
			name: "IPv4 Neighbor with IPv4 and IPv6 Prefixes",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65020,
									ID:  "192.0.2.10",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65021,
											Address: "192.0.2.11",
											Port:    179,
										},
									},
									VRF:      "",
									Prefixes: []string{"192.0.2.0/24", "2001:db8::/64"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65020,
						RouterID: "192.0.2.10",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65021@192.0.2.11",
								ASN:      65021,
								Addr:     "192.0.2.11",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						IPV4Prefixes: []string{"192.0.2.0/24"},
						IPV6Prefixes: []string{"2001:db8::/64"},
					},
				},
			},
			err: nil,
		},
		{
			name: "Empty Configuration",
			fromK8s: []v1beta1.FRRConfiguration{
				{},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{},
			},
			err: nil,
		},
		{
			name: "Non default VRF",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65030,
									ID:  "192.0.2.15",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65031,
											Address: "192.0.2.16",
											Port:    179,
										},
									},
									VRF:      "vrf1",
									Prefixes: []string{"192.0.2.0/24"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65030,
						RouterID: "192.0.2.15",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65031@192.0.2.16",
								ASN:      65031,
								Addr:     "192.0.2.16",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						VRF:          "vrf1",
						IPV4Prefixes: []string{"192.0.2.0/24"},
						IPV6Prefixes: []string{},
					},
				},
			},
			err: nil,
		},
		{
			name: "Neighbor with ToAdvertise",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Prefixes: []string{"192.0.2.0/24"},
													Mode:     v1beta1.AllowRestricted,
												},
											},
										},
									},
									Prefixes: []string{"192.0.2.0/24"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65040,
						RouterID: "192.0.2.20",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.21",
								ASN:      65041,
								Addr:     "192.0.2.21",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.2.0/24",
										},
									},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						IPV4Prefixes: []string{"192.0.2.0/24"},
						IPV6Prefixes: []string{},
					},
				},
			},
			err: nil,
		},
		{
			name: "Two Neighbor with ToAdvertise, one advertise all",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Prefixes: []string{"192.0.2.0/24", "192.0.4.0/24"},
													Mode:     v1beta1.AllowRestricted,
												},
											},
										},
										{
											ASN:     65041,
											Address: "192.0.2.22",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Mode: v1beta1.AllowAll,
												},
											},
										},
									},
									Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24", "2001:db8::/64"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65040,
						RouterID: "192.0.2.20",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.21",
								ASN:      65041,
								Addr:     "192.0.2.21",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.2.0/24",
										},
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.4.0/24",
										},
									},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.22",
								ASN:      65041,
								Addr:     "192.0.2.22",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.2.0/24",
										},
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.3.0/24",
										},
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.4.0/24",
										},
									},
									PrefixesV6: []frr.OutgoingFilter{
										{
											IPFamily: ipfamily.IPv6,
											Prefix:   "2001:db8::/64",
										},
									},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						IPV4Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24"},
						IPV6Prefixes: []string{"2001:db8::/64"},
					},
				},
			},
			err: nil,
		},
		{
			name: "Two Neighbor with ToAdvertise, one advertise all, both with communities",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Prefixes: []string{"192.0.2.0/24", "192.0.4.0/24"},
													Mode:     v1beta1.AllowRestricted,
												},
												PrefixesWithCommunity: []v1beta1.CommunityPrefixes{
													{
														Prefixes:  []string{"192.0.2.0/24", "192.0.4.0/24"},
														Community: "10:100",
													},
													{
														Prefixes:  []string{"192.0.2.0/24"},
														Community: "10:102",
													},
													{
														Prefixes:  []string{"192.0.2.0/24", "192.0.4.0/24"},
														Community: "large:123:456:7890",
													},
													{
														Prefixes:  []string{"192.0.4.0/24"},
														Community: "large:123:456:7892",
													},
													{
														Prefixes:  []string{"192.0.4.0/24"},
														Community: "10:104",
													},
												},
											},
										},
										{
											ASN:     65041,
											Address: "192.0.2.22",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Mode: v1beta1.AllowAll,
												},
												PrefixesWithCommunity: []v1beta1.CommunityPrefixes{
													{
														Prefixes:  []string{"192.0.2.0/24", "192.0.4.0/24"},
														Community: "10:100",
													},
													{
														Prefixes:  []string{"192.0.2.0/24"},
														Community: "10:102",
													},
													{
														Prefixes:  []string{"192.0.2.0/24", "2001:db8::/64"},
														Community: "10:108",
													},
												},
											},
										},
									},
									Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24", "2001:db8::/64"},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:    65040,
						RouterID: "192.0.2.20",
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.21",
								ASN:      65041,
								Addr:     "192.0.2.21",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{
										{
											IPFamily:         ipfamily.IPv4,
											Prefix:           "192.0.2.0/24",
											Communities:      []string{"10:100", "10:102"},
											LargeCommunities: []string{"123:456:7890"},
										},
										{
											IPFamily:         ipfamily.IPv4,
											Prefix:           "192.0.4.0/24",
											Communities:      []string{"10:100", "10:104"},
											LargeCommunities: []string{"123:456:7890", "123:456:7892"},
										},
									},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.22",
								ASN:      65041,
								Addr:     "192.0.2.22",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{
										{
											IPFamily:    ipfamily.IPv4,
											Prefix:      "192.0.2.0/24",
											Communities: []string{"10:100", "10:102", "10:108"},
										},
										{
											IPFamily: ipfamily.IPv4,
											Prefix:   "192.0.3.0/24",
										},
										{
											IPFamily:    ipfamily.IPv4,
											Prefix:      "192.0.4.0/24",
											Communities: []string{"10:100"},
										},
									},
									PrefixesV6: []frr.OutgoingFilter{
										{
											IPFamily:    ipfamily.IPv6,
											Prefix:      "2001:db8::/64",
											Communities: []string{"10:108"},
										},
									},
								},
								Incoming: frr.AllowedIn{
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
						IPV4Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24"},
						IPV6Prefixes: []string{"2001:db8::/64"},
					},
				},
			},
			err: nil,
		},
		{
			name: "One neighbor, trying to set community on an unallowed prefix",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToAdvertise: v1beta1.Advertise{
												Allowed: v1beta1.AllowedPrefixes{
													Prefixes: []string{"192.0.2.0/24", "192.0.4.0/24"},
													Mode:     v1beta1.AllowRestricted,
												},
												PrefixesWithCommunity: []v1beta1.CommunityPrefixes{
													{
														Prefixes:  []string{"192.0.2.0/24", "192.0.4.0/24"},
														Community: "10:100",
													},
													{
														Prefixes:  []string{"192.0.10.10/32"}, // not allowed
														Community: "10:100",
													},
												},
											},
										},
									},
									Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24", "2001:db8::/64"},
								},
							},
						},
					},
				},
			},
			expected: nil,
			err:      fmt.Errorf("prefix %s with community %s not in allowed list for neighbor %s", "192.0.10.10/32", "10:100", "192.0.2.21"),
		},
		{
			name: "Neighbor with ToReceiveAll",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToReceive: v1beta1.Receive{
												Allowed: v1beta1.AllowedPrefixes{
													Mode: v1beta1.AllowAll,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:        65040,
						RouterID:     "192.0.2.20",
						IPV4Prefixes: []string{},
						IPV6Prefixes: []string{},
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.21",
								ASN:      65041,
								Addr:     "192.0.2.21",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									All:        true,
									PrefixesV4: []frr.IncomingFilter{},
									PrefixesV6: []frr.IncomingFilter{},
								},
							},
						},
					},
				},
			},
			err: nil,
		}, {
			name: "Neighbor with ToReceive some ips only",
			fromK8s: []v1beta1.FRRConfiguration{
				{
					Spec: v1beta1.FRRConfigurationSpec{
						BGP: v1beta1.BGPConfig{
							Routers: []v1beta1.Router{
								{
									ASN: 65040,
									ID:  "192.0.2.20",
									Neighbors: []v1beta1.Neighbor{
										{
											ASN:     65041,
											Address: "192.0.2.21",
											Port:    179,
											ToReceive: v1beta1.Receive{
												Allowed: v1beta1.AllowedPrefixes{
													Prefixes: []string{"192.0.2.0/24", "192.0.3.0/24", "192.0.4.0/24", "2001:db8::/64"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &frr.Config{
				Routers: []*frr.RouterConfig{
					{
						MyASN:        65040,
						RouterID:     "192.0.2.20",
						IPV4Prefixes: []string{},
						IPV6Prefixes: []string{},
						Neighbors: []*frr.NeighborConfig{
							{
								IPFamily: ipfamily.IPv4,
								Name:     "65041@192.0.2.21",
								ASN:      65041,
								Addr:     "192.0.2.21",
								Port:     179,
								Outgoing: frr.AllowedOut{
									PrefixesV4: []frr.OutgoingFilter{},
									PrefixesV6: []frr.OutgoingFilter{},
								},
								Incoming: frr.AllowedIn{
									All: false,
									PrefixesV4: []frr.IncomingFilter{
										{IPFamily: "ipv4", Prefix: "192.0.2.0/24"},
										{IPFamily: "ipv4", Prefix: "192.0.3.0/24"},
										{IPFamily: "ipv4", Prefix: "192.0.4.0/24"},
									},
									PrefixesV6: []frr.IncomingFilter{
										{IPFamily: "ipv6", Prefix: "2001:db8::/64"},
									},
								},
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frr, err := apiToFRR(test.fromK8s[0]) // TODO: pass the array when we start supporting merge
			if test.err != nil && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if test.err == nil && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if diff := cmp.Diff(frr, test.expected); diff != "" {
				t.Fatalf("config different from expected: %s", diff)
			}
		})
	}
}