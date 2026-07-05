/*
 * Copyright 2026 Vitruvian Software
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/accesscontextmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumiverse/pulumi-time/sdk/go/time"

	networking "github.com/VitruvianSoftware/pulumi-library/go/pkg/network"
	netpeering "github.com/VitruvianSoftware/pulumi-library/go/pkg/network_peering"
	vpc_sc "github.com/VitruvianSoftware/pulumi-library/go/pkg/vpc_service_controls"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := loadNetConfig(ctx)

		// Prepend the bootstrap deploy service accounts to the VPC-SC perimeter
		// members (read live from the bootstrap stack, mirroring upstream).
		cfg.DeploySAMembers = readBootstrapDeploySAs(ctx, cfg.BootstrapStackName)

		// Resolve the org stack once for the whole stage. It publishes the net-hub
		// project id (net_hub_project_id), the per-env spoke network project id
		// (<env>_network_project_id) and the access-context-manager policy id/number
		// consumed below. Reading the project ids here mirrors TF's
		// data.terraform_remote_state.org wiring, so they no longer have to be
		// duplicated in this stack's own config. GetOutputDetails resolves
		// synchronously, so the values are usable as plain strings; each falls back
		// to its config value when the org stack has not been deployed yet
		// (e.g. preview / unit tests).
		var orgStack *pulumi.StackReference
		if cfg.OrgStackName != "" {
			var err error
			orgStack, err = pulumi.NewStackReference(ctx, "org", &pulumi.StackReferenceArgs{
				Name: pulumi.String(cfg.OrgStackName),
			})
			if err != nil {
				return err
			}
			// net_hub_project_id is needed by the hub stack and by each spoke
			// (for the spoke<->hub peering reference).
			if id := stackOutputString(orgStack, "net_hub_project_id"); id != "" {
				cfg.HubProjectID = id
			}
			// <env>_network_project_id backs the spoke shared-VPC host (spoke branch only).
			if id := stackOutputString(orgStack, fmt.Sprintf("%s_network_project_id", cfg.Env)); id != "" {
				cfg.SpokeProjectID = id
			}
		}

		// Compute environment-specific advertised IP ranges
		advertisedRanges := []networking.AdvertisedIPRange{
			{Range: cfg.PscIP + "/32", Description: "PSC Endpoint"},
		}
		// Hub (env_code "c") also advertises the DNS forwarding source range
		if cfg.EnvCode == "c" {
			advertisedRanges = append([]networking.AdvertisedIPRange{
				{Range: "35.199.192.0/19", Description: "Google DNS Forwarding Source"},
			}, advertisedRanges...)
		}

		// ====================================================================
		// HUB ENVIRONMENT (Deployed in "shared" / env_code "c")
		// ====================================================================
		if cfg.Env == "shared" {
			// 1. Hierarchical Firewall Policy (org/folder level)
			_, err := networking.NewHierarchicalFirewallPolicy(ctx, "hierarchical-fw", &networking.HierarchicalFirewallPolicyArgs{
				ParentID:      pulumi.String(cfg.ParentID),
				ShortName:     fmt.Sprintf("fw-%s-svpc-hierarchical", cfg.Env),
				Description:   "Hierarchical firewall rules",
				Associations:  cfg.FirewallAssociations,
				EnableLogging: cfg.FirewallPoliciesEnableLogging,
			})
			if err != nil {
				return err
			}

			// 2. Hub Shared VPC Host
			if _, err := compute.NewSharedVPCHostProject(ctx, "hub-svpc-host", &compute.SharedVPCHostProjectArgs{
				Project: pulumi.String(cfg.HubProjectID),
			}); err != nil {
				return err
			}

			// 3. Hub VPC & Subnets
			hubVpcName := fmt.Sprintf("vpc-%s-svpc-hub", cfg.EnvCode)
			hubNetOpts := &networking.NetworkingArgs{
				ProjectID: pulumi.String(cfg.HubProjectID),
				VPCName:   pulumi.String(hubVpcName),
				EnablePSA: true,
				Subnets: []networking.SubnetArgs{
					{
						// Upstream hub secondary_ranges = {} (no gke-pod/gke-svc on the hub).
						Name:             fmt.Sprintf("sb-%s-svpc-hub-%s", cfg.EnvCode, cfg.Region1),
						Region:           cfg.Region1,
						CIDR:             cfg.HubSubnet1Cidr,
						FlowLogs:         true,
						FlowLogsInterval: cfg.VpcFlowLogs.AggregationInterval,
						FlowLogsSampling: cfg.VpcFlowLogs.FlowSampling,
						FlowLogsMetadata: cfg.VpcFlowLogs.Metadata,
					},
					{
						Name:             fmt.Sprintf("sb-%s-svpc-hub-%s", cfg.EnvCode, cfg.Region2),
						Region:           cfg.Region2,
						CIDR:             cfg.HubSubnet2Cidr,
						FlowLogs:         true,
						FlowLogsInterval: cfg.VpcFlowLogs.AggregationInterval,
						FlowLogsSampling: cfg.VpcFlowLogs.FlowSampling,
						FlowLogsMetadata: cfg.VpcFlowLogs.Metadata,
					},
					{
						// Hub regional managed-proxy subnets (upstream net-hubs.tf) so the
						// hub VPC can host L7 internal load balancers.
						Name:    fmt.Sprintf("sb-%s-svpc-hub-%s-proxy", cfg.EnvCode, cfg.Region1),
						Region:  cfg.Region1,
						CIDR:    cfg.HubProxy1Cidr,
						Role:    "ACTIVE",
						Purpose: "REGIONAL_MANAGED_PROXY",
					},
					{
						Name:    fmt.Sprintf("sb-%s-svpc-hub-%s-proxy", cfg.EnvCode, cfg.Region2),
						Region:  cfg.Region2,
						CIDR:    cfg.HubProxy2Cidr,
						Role:    "ACTIVE",
						Purpose: "REGIONAL_MANAGED_PROXY",
					},
				},
			}

			hubVpc, err := networking.NewNetworking(ctx, "hub", hubNetOpts)
			if err != nil {
				return err
			}

			// Hub Egress internet route (tag-based, only when NAT is enabled)
			if cfg.NatEnabled {
				_, err = compute.NewRoute(ctx, "hub-egress-internet", &compute.RouteArgs{
					Project:        pulumi.String(cfg.HubProjectID),
					Name:           pulumi.String(fmt.Sprintf("rt-%s-hub-1000-egress-internet-default", cfg.EnvCode)),
					Network:        hubVpc.VPC.ID(),
					DestRange:      pulumi.String("0.0.0.0/0"),
					NextHopGateway: pulumi.String("default-internet-gateway"),
					Priority:       pulumi.Int(1000),
					Tags:           pulumi.StringArray{pulumi.String("egress-internet")},
				}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
				if err != nil {
					return err
				}
			}

			// Windows KMS activation route (upstream windows_activation_enabled).
			if cfg.WindowsActivationEnabled {
				_, err = compute.NewRoute(ctx, "hub-windows-kms", &compute.RouteArgs{
					Project:        pulumi.String(cfg.HubProjectID),
					Name:           pulumi.String(fmt.Sprintf("rt-%s-svpc-hub-1000-all-default-windows-kms", cfg.EnvCode)),
					Network:        hubVpc.VPC.ID(),
					DestRange:      pulumi.String("35.190.247.13/32"),
					NextHopGateway: pulumi.String("default-internet-gateway"),
					Priority:       pulumi.Int(1000),
				}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
				if err != nil {
					return err
				}
			}

			// 4. Hub VPC-Level Firewall — data-driven rules
			hubFw, err := networking.NewNetworkFirewallPolicy(ctx, "hub-vpc-fw", &networking.NetworkFirewallPolicyArgs{
				ProjectID:  pulumi.String(cfg.HubProjectID),
				PolicyName: fmt.Sprintf("fp-%s-hub-firewalls", cfg.EnvCode),
				TargetVPCs: []pulumi.StringInput{
					pulumi.Sprintf("projects/%s/global/networks/%s", cfg.HubProjectID, hubVpc.VPC.Name),
				},
				Rules: networking.BuildFoundationRules(cfg.EnvCode, cfg.FirewallPoliciesEnableLogging, cfg.PscIP+"/32", []string{cfg.HubSubnet1Cidr, cfg.HubSubnet2Cidr}, false),
			}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
			if err != nil {
				return err
			}

			// 5. PSC on hub
			_, err = networking.NewPrivateServiceConnect(ctx, "hub-psc", &networking.PrivateServiceConnectArgs{
				ProjectID:            pulumi.String(cfg.HubProjectID),
				NetworkSelfLink:      hubVpc.VPC.SelfLink,
				DnsCode:              fmt.Sprintf("dz-%s-hub", cfg.EnvCode),
				IPAddress:            cfg.PscIP,
				ForwardingRuleTarget: "vpc-sc",
			}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
			if err != nil {
				return err
			}

			// 6. DNS Policy on hub
			hubDnsPolicy, err := dns.NewPolicy(ctx, "hub-dns-policy", &dns.PolicyArgs{
				Project:                 pulumi.String(cfg.HubProjectID),
				Name:                    pulumi.String(fmt.Sprintf("dp-%s-hub-default-policy", cfg.EnvCode)),
				EnableInboundForwarding: pulumi.Bool(true),
				EnableLogging:           pulumi.Bool(cfg.DnsEnableLogging),
				Networks: dns.PolicyNetworkArray{
					&dns.PolicyNetworkArgs{
						NetworkUrl: hubVpc.VPC.SelfLink,
					},
				},
			}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
			if err != nil {
				return err
			}

			// 7. DNS forwarding zone on hub
			_, err = networking.NewDnsZone(ctx, "dns-forwarding", &networking.DnsZoneArgs{
				ProjectID:                 pulumi.String(cfg.HubProjectID),
				Name:                      "fz-dns-hub",
				Domain:                    cfg.Domain,
				Type:                      "forwarding",
				NetworkSelfLink:           hubVpc.VPC.SelfLink,
				TargetNameServerAddresses: cfg.TargetNameServers,
			})
			if err != nil {
				return err
			}

			// 8. Transitivity Appliance — MIG+ILB per region.
			// Gated on enable_hub_and_spoke_transitivity (default true preserves the
			// prior always-on behavior; set false to disable).
			if cfg.EnableTransitivity {
				_, err = networking.NewTransitivityAppliance(ctx, "transitivity", &networking.TransitivityApplianceArgs{
					ProjectID:   pulumi.String(cfg.HubProjectID),
					Regions:     []string{cfg.Region1, cfg.Region2},
					Network:     hubVpc.VPC.SelfLink,
					NetworkName: hubVpcName,
					Subnetworks: map[string]pulumi.StringInput{
						cfg.Region1: hubVpc.Subnets[fmt.Sprintf("sb-%s-svpc-hub-%s", cfg.EnvCode, cfg.Region1)].SelfLink,
						cfg.Region2: hubVpc.Subnets[fmt.Sprintf("sb-%s-svpc-hub-%s", cfg.EnvCode, cfg.Region2)].SelfLink,
					},
					// Upstream regional_aggregates (net-hubs-transitivity.tf).
					RegionalAggregates: map[string][]string{
						cfg.Region1: {"10.8.0.0/16", "100.72.0.0/16"},
						cfg.Region2: {"10.9.0.0/16", "100.73.0.0/16"},
					},
					FirewallPolicy: hubFw.Policy.Name,
				}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
				if err != nil {
					return err
				}

				// 8.1. Hub Firewall — Allow Health Checks to Transitivity ILBs
				_, err = compute.NewFirewall(ctx, "fw-hub-allow-health-checks", &compute.FirewallArgs{
					Project: pulumi.String(cfg.HubProjectID),
					Name:    pulumi.String(fmt.Sprintf("fw-%s-hub-allow-health-checks", cfg.EnvCode)),
					Network: hubVpc.VPC.SelfLink,
					Allows: compute.FirewallAllowArray{
						&compute.FirewallAllowArgs{
							Protocol: pulumi.String("tcp"),
							Ports:    pulumi.StringArray{pulumi.String("22")},
						},
					},
					SourceRanges: pulumi.StringArray{
						pulumi.String("130.211.0.0/22"),
						pulumi.String("35.191.0.0/16"),
					},
					TargetTags: pulumi.StringArray{
						pulumi.String("allow-transitivity"),
					},
				}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
				if err != nil {
					return err
				}
			} // end enable_hub_and_spoke_transitivity

			// 9. Hub BGP Routers — 4 total (2 per region), hub only (not on spokes).
			// Upstream names region1 routers cr5/cr6 and region2 routers cr7/cr8.
			for _, rc := range []struct {
				region string
				crIdxs []string
			}{
				{cfg.Region1, []string{"5", "6"}},
				{cfg.Region2, []string{"7", "8"}},
			} {
				for _, crIdx := range rc.crIdxs {
					_, err = networking.NewCloudRouter(ctx, fmt.Sprintf("hub-cr-%s-cr%s", rc.region, crIdx), &networking.RouterArgs{
						ProjectID:          pulumi.String(cfg.HubProjectID),
						Region:             rc.region,
						Network:            hubVpc.VPC.SelfLink,
						BgpAsn:             cfg.BgpAsn,
						AdvertisedGroups:   []string{"ALL_SUBNETS"},
						AdvertisedIpRanges: advertisedRanges,
						EnableNat:          false,
					}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
					if err != nil {
						return err
					}
				}
			}

			// 10. Separate NAT routers on hub (gated on nat_enabled).
			if cfg.NatEnabled {
				for _, reg := range []string{cfg.Region1, cfg.Region2} {
					_, err = networking.NewCloudRouter(ctx, fmt.Sprintf("hub-nat-%s", reg), &networking.RouterArgs{
						ProjectID:       pulumi.String(cfg.HubProjectID),
						Region:          reg,
						Network:         hubVpc.VPC.SelfLink,
						BgpAsn:          cfg.NatBgpAsn,
						EnableNat:       true,
						NatNumAddresses: cfg.NatNumAddresses,
					}, pulumi.DependsOn([]pulumi.Resource{hubVpc.VPC}))
					if err != nil {
						return err
					}
				}
			}

			// VPC-SC perimeter on the hub project — created once here in the shared/hub
			// stack that owns net-hub (TF applies it unconditionally in envs/shared).
			if orgStack != nil {
				var hubPolicyID pulumi.StringInput = orgStack.GetStringOutput(pulumi.String("access_context_manager_policy_id"))
				if cfg.PolicyID != "" {
					hubPolicyID = pulumi.String(cfg.PolicyID)
				}
				// Prepend the bootstrap deploy service accounts (organization/networks/
				// projects step SAs) to the perimeter members so the deploy pipeline is
				// not locked out (upstream net-hubs.tf members/members_dry_run).
				hubPerimeterMembers := mergeMembers(cfg.DeploySAMembers, cfg.VpcScMembers)
				_, err = vpc_sc.NewVpcServiceControls(ctx, "hub-vpc-sc-perimeter", &vpc_sc.VpcServiceControlsArgs{
					PolicyID:           hubPolicyID,
					Prefix:             "c_hub",
					Members:            hubPerimeterMembers,
					MembersDryRun:      hubPerimeterMembers,
					ProjectNumbers:     pulumi.StringArray{orgStack.GetStringOutput(pulumi.String("net_hub_project_number"))},
					RestrictedServices: cfg.VpcScRestrictedServices,
					Enforce:            cfg.EnforceVpcSc,
				})
				if err != nil {
					return err
				}
			}

			// Exports — matching TF 3-networks-hub-and-spoke/envs/shared/outputs.tf
			ctx.Export("shared_vpc_host_project_id", pulumi.String(cfg.HubProjectID))
			ctx.Export("network_name", hubVpc.VPC.Name)
			ctx.Export("dns_policy", hubDnsPolicy.ID()) // DNS policy ID
			return nil
		}

		// ====================================================================
		// SPOKE ENVIRONMENT (dev, nonprod, prod)
		// ====================================================================

		if _, err := compute.NewSharedVPCHostProject(ctx, "spoke-svpc-host", &compute.SharedVPCHostProjectArgs{
			Project: pulumi.String(cfg.SpokeProjectID),
		}); err != nil {
			return err
		}

		spokeVpcName := fmt.Sprintf("vpc-%s-svpc-spoke", cfg.EnvCode)
		spokeNetOpts := &networking.NetworkingArgs{
			ProjectID: pulumi.String(cfg.SpokeProjectID),
			VPCName:   pulumi.String(spokeVpcName),
			EnablePSA: true,
			Subnets: []networking.SubnetArgs{
				{
					Name:   fmt.Sprintf("sb-%s-svpc-spoke-%s", cfg.EnvCode, cfg.Region1),
					Region: cfg.Region1,
					CIDR:   cfg.SpokeSubnet1Cidr,
					SecondaryRanges: []networking.SecondaryRangeArgs{
						{RangeName: fmt.Sprintf("rn-%s-spoke-%s-gke-pod", cfg.EnvCode, cfg.Region1), CIDR: cfg.SpokeGkePod1Cidr},
						{RangeName: fmt.Sprintf("rn-%s-spoke-%s-gke-svc", cfg.EnvCode, cfg.Region1), CIDR: cfg.SpokeGkeSvc1Cidr},
					},
					FlowLogs:         true,
					FlowLogsInterval: cfg.VpcFlowLogs.AggregationInterval,
					FlowLogsSampling: cfg.VpcFlowLogs.FlowSampling,
					FlowLogsMetadata: cfg.VpcFlowLogs.Metadata,
				},
				{
					Name:   fmt.Sprintf("sb-%s-svpc-spoke-%s", cfg.EnvCode, cfg.Region2),
					Region: cfg.Region2,
					CIDR:   cfg.SpokeSubnet2Cidr,
					SecondaryRanges: []networking.SecondaryRangeArgs{
						{RangeName: fmt.Sprintf("rn-%s-spoke-%s-gke-pod", cfg.EnvCode, cfg.Region2), CIDR: cfg.SpokeGkePod2Cidr},
						{RangeName: fmt.Sprintf("rn-%s-spoke-%s-gke-svc", cfg.EnvCode, cfg.Region2), CIDR: cfg.SpokeGkeSvc2Cidr},
					},
					FlowLogs:         true,
					FlowLogsInterval: cfg.VpcFlowLogs.AggregationInterval,
					FlowLogsSampling: cfg.VpcFlowLogs.FlowSampling,
					FlowLogsMetadata: cfg.VpcFlowLogs.Metadata,
				},
				{
					Name:    fmt.Sprintf("sb-%s-svpc-spoke-%s-proxy", cfg.EnvCode, cfg.Region1),
					Region:  cfg.Region1,
					CIDR:    cfg.SpokeProxy1Cidr,
					Role:    "ACTIVE",
					Purpose: "REGIONAL_MANAGED_PROXY",
				},
				{
					Name:    fmt.Sprintf("sb-%s-svpc-spoke-%s-proxy", cfg.EnvCode, cfg.Region2),
					Region:  cfg.Region2,
					CIDR:    cfg.SpokeProxy2Cidr,
					Role:    "ACTIVE",
					Purpose: "REGIONAL_MANAGED_PROXY",
				},
			},
		}

		spokeVpc, err := networking.NewNetworking(ctx, "spoke", spokeNetOpts)
		if err != nil {
			return err
		}

		// Bi-Directional VPC Peering (Spoke <-> Hub). The component creates both
		// reciprocal peerings — the spoke side (imports the hub's custom routes)
		// and the inverted hub side (exports them), ordered so the hub side
		// depends on the spoke side. Explicit names reproduce the upstream
		// "np-<local>-<peer>" convention.
		hubVpcRef := fmt.Sprintf("projects/%s/global/networks/vpc-c-svpc-hub", cfg.HubProjectID)

		_, err = netpeering.NewNetworkPeering(ctx, "spoke-hub", &netpeering.NetworkPeeringArgs{
			LocalNetwork:       spokeVpc.VPC.SelfLink,
			PeerNetwork:        pulumi.String(hubVpcRef),
			ImportCustomRoutes: true, // spoke imports the hub's custom routes
			LocalName:          fmt.Sprintf("np-%s-svpc-spoke-vpc-c-svpc-hub", cfg.EnvCode),
			PeerName:           fmt.Sprintf("np-vpc-c-svpc-hub-%s-svpc-spoke", cfg.EnvCode),
		})
		if err != nil {
			return err
		}

		// Spoke VPC-Level Firewall — data-driven rules
		_, err = networking.NewNetworkFirewallPolicy(ctx, "spoke-vpc-fw", &networking.NetworkFirewallPolicyArgs{
			ProjectID:  pulumi.String(cfg.SpokeProjectID),
			PolicyName: fmt.Sprintf("fp-%s-spoke-firewalls", cfg.EnvCode),
			TargetVPCs: []pulumi.StringInput{
				pulumi.Sprintf("projects/%s/global/networks/%s", cfg.SpokeProjectID, spokeVpc.VPC.Name),
			},
			Rules: networking.BuildFoundationRules(cfg.EnvCode, cfg.FirewallPoliciesEnableLogging, cfg.PscIP+"/32", []string{cfg.SpokeSubnet1Cidr, cfg.SpokeSubnet2Cidr}, false),
		}, pulumi.DependsOn([]pulumi.Resource{spokeVpc.VPC}))
		if err != nil {
			return err
		}

		// PSC on spoke
		_, err = networking.NewPrivateServiceConnect(ctx, "spoke-psc", &networking.PrivateServiceConnectArgs{
			ProjectID:            pulumi.String(cfg.SpokeProjectID),
			NetworkSelfLink:      spokeVpc.VPC.SelfLink,
			DnsCode:              fmt.Sprintf("dz-%s-spoke", cfg.EnvCode),
			IPAddress:            cfg.PscIP,
			ForwardingRuleTarget: "vpc-sc",
		}, pulumi.DependsOn([]pulumi.Resource{spokeVpc.VPC}))
		if err != nil {
			return err
		}

		// DNS Policy on spoke
		_, err = dns.NewPolicy(ctx, "spoke-dns-policy", &dns.PolicyArgs{
			Project:                 pulumi.String(cfg.SpokeProjectID),
			Name:                    pulumi.String(fmt.Sprintf("dp-%s-spoke-default-policy", cfg.EnvCode)),
			EnableInboundForwarding: pulumi.Bool(true),
			EnableLogging:           pulumi.Bool(cfg.DnsEnableLogging),
			Networks: dns.PolicyNetworkArray{
				&dns.PolicyNetworkArgs{
					NetworkUrl: spokeVpc.VPC.SelfLink,
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{spokeVpc.VPC}))
		if err != nil {
			return err
		}

		// DNS peering from spoke to hub
		_, err = networking.NewDnsZone(ctx, "dns-peering", &networking.DnsZoneArgs{
			ProjectID:             pulumi.String(cfg.SpokeProjectID),
			Name:                  fmt.Sprintf("dz-%s-svpc-to-dns-hub", cfg.EnvCode),
			Domain:                cfg.Domain,
			Type:                  "peering",
			NetworkSelfLink:       spokeVpc.VPC.SelfLink,
			TargetNetworkSelfLink: pulumi.String(hubVpcRef),
		})
		if err != nil {
			return err
		}

		// NAT routers on spoke (spokes don't get BGP routers in hub-and-spoke).
		// Gated on nat_enabled.
		if cfg.NatEnabled {
			for _, reg := range []string{cfg.Region1, cfg.Region2} {
				_, err = networking.NewCloudRouter(ctx, fmt.Sprintf("spoke-nat-%s", reg), &networking.RouterArgs{
					ProjectID:       pulumi.String(cfg.SpokeProjectID),
					Region:          reg,
					Network:         spokeVpc.VPC.SelfLink,
					BgpAsn:          cfg.NatBgpAsn,
					EnableNat:       true,
					NatNumAddresses: cfg.NatNumAddresses,
				}, pulumi.DependsOn([]pulumi.Resource{spokeVpc.VPC}))
				if err != nil {
					return err
				}
			}
		}

		// Windows KMS activation route on spoke (upstream windows_activation_enabled).
		if cfg.WindowsActivationEnabled {
			_, err = compute.NewRoute(ctx, "spoke-windows-kms", &compute.RouteArgs{
				Project:        pulumi.String(cfg.SpokeProjectID),
				Name:           pulumi.String(fmt.Sprintf("rt-%s-svpc-spoke-1000-all-default-windows-kms", cfg.EnvCode)),
				Network:        spokeVpc.VPC.ID(),
				DestRange:      pulumi.String("35.190.247.13/32"),
				NextHopGateway: pulumi.String("default-internet-gateway"),
				Priority:       pulumi.Int(1000),
			}, pulumi.DependsOn([]pulumi.Resource{spokeVpc.VPC}))
			if err != nil {
				return err
			}
		}

		var acmPolicyID pulumi.StringOutput
		if orgStack != nil {
			acmPolicyID = orgStack.GetStringOutput(pulumi.String("access_context_manager_policy_id"))
		} else {
			acmPolicyID = pulumi.String("").ToStringOutput()
		}

		var finalPolicyID pulumi.StringInput
		if cfg.PolicyID != "" {
			finalPolicyID = pulumi.String(cfg.PolicyID)
		} else {
			finalPolicyID = acmPolicyID
		}

		// VPC-SC on spoke
		var perimeterName pulumi.StringOutput
		var accessLevelName pulumi.StringOutput
		var accessLevelDryRunName pulumi.StringOutput
		{
			// Prepend the bootstrap deploy service accounts to the perimeter members
			// (upstream base_env members/members_dry_run).
			spokePerimeterMembers := mergeMembers(cfg.DeploySAMembers, cfg.VpcScMembers)
			perimeter, err := vpc_sc.NewVpcServiceControls(ctx, "vpc-sc-perimeter", &vpc_sc.VpcServiceControlsArgs{
				PolicyID:              finalPolicyID,
				Prefix:                fmt.Sprintf("%s_spoke", cfg.EnvCode),
				Members:               spokePerimeterMembers,
				MembersDryRun:         spokePerimeterMembers,
				ProjectNumbers:        pulumi.ToStringArray(cfg.VpcScProjects),
				RestrictedServices:    cfg.VpcScRestrictedServices,
				Enforce:               cfg.EnforceVpcSc,
				IngressPolicies:       cfg.VpcScIngressPolicies,
				EgressPolicies:        cfg.VpcScEgressPolicies,
				IngressPoliciesDryRun: cfg.VpcScIngressPoliciesDryRun,
				EgressPoliciesDryRun:  cfg.VpcScEgressPoliciesDryRun,
			})
			if err != nil {
				return err
			}

			vpcScSleep, err := time.NewSleep(ctx, "vpc-sc-propagation-wait", &time.SleepArgs{
				CreateDuration: pulumi.String("60s"),
			}, pulumi.DependsOn([]pulumi.Resource{perimeter.Perimeter}))
			if err != nil {
				return err
			}

			perimeterName = pulumi.All(vpcScSleep.ID(), perimeter.Perimeter.Name).ApplyT(func(args []interface{}) string {
				return args[1].(string)
			}).(pulumi.StringOutput)

			accessLevelName = perimeter.AccessLevel.Name
			accessLevelDryRunName = perimeter.AccessLevelDryRun.Name
		}

		// Exports — matching TF 3-networks-hub-and-spoke/envs/{env}/outputs.tf
		ctx.Export("access_context_manager_policy_id", acmPolicyID)
		ctx.Export("shared_vpc_host_project_id", pulumi.String(cfg.SpokeProjectID))
		ctx.Export("network_name", spokeVpc.VPC.Name)
		ctx.Export("network_self_link", spokeVpc.VPC.SelfLink)

		// Subnet exports as arrays (matching TF subnets_names/ips/self_links/secondary_ranges)
		var subnetNames, subnetIPs, subnetSelfLinks pulumi.StringArray
		for _, subnet := range spokeVpc.Subnets {
			subnetNames = append(subnetNames, subnet.Name)
			subnetIPs = append(subnetIPs, subnet.IpCidrRange)
			subnetSelfLinks = append(subnetSelfLinks, subnet.SelfLink)
		}
		ctx.Export("subnets_names", subnetNames)
		ctx.Export("subnets_ips", subnetIPs)
		ctx.Export("subnets_self_links", subnetSelfLinks)
		// Secondary ranges: build a list from each subnet's secondary_ip_ranges.
		// TF outputs this as a list of objects with range_name and ip_cidr_range.
		var secondaryRangesList pulumi.ArrayOutput
		for _, subnet := range spokeVpc.Subnets {
			secondaryRangesList = pulumi.All(secondaryRangesList, subnet.SecondaryIpRanges).ApplyT(func(args []interface{}) []interface{} {
				existing, _ := args[0].([]interface{})
				ranges, _ := args[1].([]interface{})
				return append(existing, ranges...)
			}).(pulumi.ArrayOutput)
		}
		if secondaryRangesList == (pulumi.ArrayOutput{}) {
			ctx.Export("subnets_secondary_ranges", pulumi.ToStringArray([]string{}))
		} else {
			ctx.Export("subnets_secondary_ranges", secondaryRangesList)
		}
		ctx.Export("enforce_vpcsc", pulumi.Bool(cfg.EnforceVpcSc))
		ctx.Export("service_perimeter_name", perimeterName)
		ctx.Export("access_level_name", accessLevelName)
		ctx.Export("access_level_name_dry_run", accessLevelDryRunName)

		return nil
	})
}

type VpcFlowLogsConfig struct {
	AggregationInterval string  `json:"aggregation_interval"`
	FlowSampling        float64 `json:"flow_sampling"`
	Metadata            string  `json:"metadata"`
}

type NetConfig struct {
	Env                           string
	EnvCode                       string
	HubProjectID                  string
	SpokeProjectID                string
	Region1                       string
	Region2                       string
	ParentID                      string
	Domain                        string
	PolicyID                      string
	OrgStackName                  string
	PscIP                         string
	BgpAsn                        int
	NatBgpAsn                     int
	NatNumAddresses               int
	TargetNameServers             []string
	VpcScMembers                  []string
	VpcScProjects                 []string
	VpcScRestrictedServices       []string
	VpcScIngressPolicies          accesscontextmanager.ServicePerimeterStatusIngressPolicyArray
	VpcScEgressPolicies           accesscontextmanager.ServicePerimeterStatusEgressPolicyArray
	VpcScIngressPoliciesDryRun    accesscontextmanager.ServicePerimeterSpecIngressPolicyArray
	VpcScEgressPoliciesDryRun     accesscontextmanager.ServicePerimeterSpecEgressPolicyArray
	SpokeSubnet1Cidr              string
	SpokeSubnet2Cidr              string
	SpokeProxy1Cidr               string
	SpokeProxy2Cidr               string
	SpokeGkePod1Cidr              string
	SpokeGkeSvc1Cidr              string
	SpokeGkePod2Cidr              string
	SpokeGkeSvc2Cidr              string
	HubSubnet1Cidr                string
	HubSubnet2Cidr                string
	HubProxy1Cidr                 string
	HubProxy2Cidr                 string
	FirewallAssociations          []string
	FirewallPoliciesEnableLogging bool
	DnsEnableLogging              bool
	EnforceVpcSc                  bool
	EnableTransitivity            bool
	NatEnabled                    bool
	WindowsActivationEnabled      bool
	BootstrapStackName            string
	DeploySAMembers               []string
	VpcFlowLogs                   *VpcFlowLogsConfig
}

// readBootstrapDeploySAs reads the organization/networks/projects step deploy
// service-account emails from the bootstrap stack (mirrors upstream
// modules/base_env/remote.tf) and returns them as `serviceAccount:<email>`
// members. Degrades gracefully to nil when the stack/outputs are unavailable
// (e.g. unit tests with mocks, or a not-yet-deployed bootstrap stack).
func readBootstrapDeploySAs(ctx *pulumi.Context, stackName string) []string {
	if stackName == "" {
		return nil
	}
	stack, err := pulumi.NewStackReference(ctx, "bootstrap", &pulumi.StackReferenceArgs{
		Name: pulumi.String(stackName),
	})
	if err != nil {
		return nil
	}
	var out []string
	for _, outName := range []string{
		"organization_step_terraform_service_account_email",
		"networks_step_terraform_service_account_email",
		"projects_step_terraform_service_account_email",
	} {
		if email := stackOutputString(stack, outName); email != "" {
			out = append(out, "serviceAccount:"+email)
		}
	}
	return out
}

// stackOutputString reads a single string output from a StackReference using
// GetOutputDetails, which resolves synchronously so the value is usable as a
// plain Go string. It returns "" when the output is absent or not yet available
// (e.g. the referenced stack has not been deployed, or during preview / unit
// tests with mocks), letting callers fall back to config.
func stackOutputString(stack *pulumi.StackReference, name string) string {
	details, err := stack.GetOutputDetails(name)
	if err != nil || details.Value == nil {
		return ""
	}
	s, ok := details.Value.(string)
	if !ok {
		return ""
	}
	return s
}

// mergeMembers prepends the bootstrap deploy service-account members (organization/
// networks/projects step SAs) ahead of the operator-supplied VPC-SC members,
// de-duplicating. Mirrors upstream distinct(concat([...SAs...], additional_members)).
func mergeMembers(deploySAs, members []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(deploySAs)+len(members))
	for _, m := range deploySAs {
		if m == "" || seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	for _, m := range members {
		if m == "" || seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	return out
}

func loadNetConfig(ctx *pulumi.Context) *NetConfig {
	conf := config.New(ctx, "")

	c := &NetConfig{
		Env:            conf.Require("env"),
		EnvCode:        conf.Require("env_code"),
		HubProjectID:   conf.Get("hub_project_id"),
		SpokeProjectID: conf.Get("spoke_project_id"),
		Region1:        conf.Get("region1"),
		Region2:        conf.Get("region2"),
		ParentID:       conf.Require("parent_id"),
		Domain:         conf.Get("domain"),
		PolicyID:       conf.Get("policy_id"),
		OrgStackName:   conf.Get("org_stack_name"),
		PscIP:          conf.Get("psc_ip"),
	}
	conf.GetObject("vpc_sc_members", &c.VpcScMembers)
	conf.GetObject("vpc_sc_projects", &c.VpcScProjects)
	conf.GetObject("vpc_sc_restricted_services", &c.VpcScRestrictedServices)
	conf.GetObject("vpc_sc_ingress_policies", &c.VpcScIngressPolicies)
	conf.GetObject("vpc_sc_egress_policies", &c.VpcScEgressPolicies)
	conf.GetObject("vpc_sc_ingress_policies_dry_run", &c.VpcScIngressPoliciesDryRun)
	conf.GetObject("vpc_sc_egress_policies_dry_run", &c.VpcScEgressPoliciesDryRun)
	conf.GetObject("target_name_servers", &c.TargetNameServers)
	conf.GetObject("firewall_associations", &c.FirewallAssociations)

	var flowLogs VpcFlowLogsConfig
	if err := conf.GetObject("vpc_flow_logs", &flowLogs); err == nil {
		c.VpcFlowLogs = &flowLogs
	} else {
		// Default matches TF upstream default
		c.VpcFlowLogs = &VpcFlowLogsConfig{
			AggregationInterval: "INTERVAL_5_SEC",
			FlowSampling:        0.5,
			Metadata:            "INCLUDE_ALL_METADATA",
		}
	}

	if val, err := conf.TryBool("firewall_policies_enable_logging"); err == nil {
		c.FirewallPoliciesEnableLogging = val
	} else {
		c.FirewallPoliciesEnableLogging = true // Default to true matching TF
	}

	if val, err := conf.TryBool("dns_enable_logging"); err == nil {
		c.DnsEnableLogging = val
	} else {
		c.DnsEnableLogging = true
	}

	if val, err := conf.TryBool("enforce_vpcsc"); err == nil {
		c.EnforceVpcSc = val
	} else {
		c.EnforceVpcSc = false // TF defaults enforce_vpcsc=false (dry-run first)
	}

	if val, err := conf.TryBool("enable_hub_and_spoke_transitivity"); err == nil {
		c.EnableTransitivity = val
	} else {
		c.EnableTransitivity = true // preserve prior always-on behavior
	}

	if val, err := conf.TryBool("nat_enabled"); err == nil {
		c.NatEnabled = val
	} else {
		c.NatEnabled = true // preserve current behavior (upstream default is false)
	}

	if val, err := conf.TryBool("windows_activation_enabled"); err == nil {
		c.WindowsActivationEnabled = val
	} else {
		c.WindowsActivationEnabled = false
	}

	// Bootstrap deploy service-account emails are read from the bootstrap stack in
	// main() via readBootstrapDeploySAs (mirrors upstream modules/base_env/remote.tf);
	// bootstrap_stack_name guards the lookup and degrades gracefully when unset.
	c.BootstrapStackName = conf.Get("bootstrap_stack_name")
	if c.BootstrapStackName == "" {
		c.BootstrapStackName = "bootstrap"
	}

	if c.Region1 == "" {
		c.Region1 = "us-central1"
	}
	if c.Region2 == "" {
		c.Region2 = "us-west1"
	}
	if c.Domain == "" {
		c.Domain = "example.com."
	}
	if c.OrgStackName == "" {
		c.OrgStackName = "org"
	}
	if c.PscIP == "" {
		c.PscIP = "10.17.0.6"
	}
	if len(c.VpcScRestrictedServices) == 0 {
		c.VpcScRestrictedServices = vpc_sc.GetDefaultRestrictedServices()
	}
	if len(c.FirewallAssociations) == 0 {
		c.FirewallAssociations = []string{c.ParentID} // Fallback to parent
	}
	if len(c.TargetNameServers) == 0 {
		c.TargetNameServers = []string{"10.0.0.1"}
	}

	// Assign CIDRs based on EnvCode to avoid peering overlaps
	// Hub defaults match upstream envs/shared/net-hubs.tf.
	if c.HubSubnet1Cidr == "" {
		c.HubSubnet1Cidr = "10.8.0.0/18"
	}
	if c.HubSubnet2Cidr == "" {
		c.HubSubnet2Cidr = "10.9.0.0/18"
	}
	if c.HubProxy1Cidr == "" {
		c.HubProxy1Cidr = "10.26.0.0/23"
	}
	if c.HubProxy2Cidr == "" {
		c.HubProxy2Cidr = "10.27.0.0/23"
	}

	if c.EnvCode == "d" {
		c.SpokeSubnet1Cidr = "10.8.64.0/18"
		c.SpokeSubnet2Cidr = "10.9.64.0/18"
		c.SpokeProxy1Cidr = "10.26.2.0/23"
		c.SpokeProxy2Cidr = "10.27.2.0/23"
		c.SpokeGkePod1Cidr = "100.72.64.0/18"
		c.SpokeGkeSvc1Cidr = "100.73.64.0/18"
		c.SpokeGkePod2Cidr = "100.74.64.0/18"
		c.SpokeGkeSvc2Cidr = "100.75.64.0/18"
	} else if c.EnvCode == "n" {
		c.SpokeSubnet1Cidr = "10.8.128.0/18"
		c.SpokeSubnet2Cidr = "10.9.128.0/18"
		c.SpokeProxy1Cidr = "10.26.4.0/23"
		c.SpokeProxy2Cidr = "10.27.4.0/23"
		c.SpokeGkePod1Cidr = "100.72.128.0/18"
		c.SpokeGkeSvc1Cidr = "100.73.128.0/18"
		c.SpokeGkePod2Cidr = "100.74.128.0/18"
		c.SpokeGkeSvc2Cidr = "100.75.128.0/18"
	} else if c.EnvCode == "p" {
		c.SpokeSubnet1Cidr = "10.8.192.0/18"
		c.SpokeSubnet2Cidr = "10.9.192.0/18"
		c.SpokeProxy1Cidr = "10.26.6.0/23"
		c.SpokeProxy2Cidr = "10.27.6.0/23"
		c.SpokeGkePod1Cidr = "100.72.192.0/18"
		c.SpokeGkeSvc1Cidr = "100.73.192.0/18"
		c.SpokeGkePod2Cidr = "100.74.192.0/18"
		c.SpokeGkeSvc2Cidr = "100.75.192.0/18"
	} else {
		// Fallback
		c.SpokeSubnet1Cidr = "10.8.64.0/18"
		c.SpokeSubnet2Cidr = "10.9.64.0/18"
		c.SpokeProxy1Cidr = "10.26.2.0/23"
		c.SpokeProxy2Cidr = "10.27.2.0/23"
		c.SpokeGkePod1Cidr = "100.72.64.0/18"
		c.SpokeGkeSvc1Cidr = "100.73.64.0/18"
		c.SpokeGkePod2Cidr = "100.74.64.0/18"
		c.SpokeGkeSvc2Cidr = "100.75.64.0/18"
	}

	c.BgpAsn = 64514
	c.NatBgpAsn = 64514
	c.NatNumAddresses = 2

	return c
}
