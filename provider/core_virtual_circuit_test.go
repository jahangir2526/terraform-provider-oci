// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	VirtualCircuitPublicRequiredOnlyResource = VirtualCircuitResourceDependencies + `
resource "oci_core_virtual_circuit" "test_virtual_circuit" {
	#Required
	compartment_id = "${var.compartment_id}"
	type = "${var.virtual_circuit_type}"

 	#Required for PUBLIC Virtual Circuit
	cross_connect_mappings {
		cross_connect_or_cross_connect_group_id = "${oci_core_cross_connect.test_cross_connect.cross_connect_group_id}"
		vlan = "${var.virtual_circuit_cross_connect_mappings_vlan}"
	}
	customer_bgp_asn = "${var.virtual_circuit_customer_bgp_asn}"
	public_prefixes {
		cidr_block = "${var.virtual_circuit_public_prefixes_cidr_block}"
	}
}
`

	VirtualCircuitRequiredOnlyResource = VirtualCircuitResourceDependencies + `
resource "oci_core_virtual_circuit" "test_virtual_circuit" {
	#Required
	compartment_id = "${var.compartment_id}"
	type = "${var.virtual_circuit_type}"

 	#Required for PRIVATE Virtual Circuit
	cross_connect_mappings {
		cross_connect_or_cross_connect_group_id = "${oci_core_cross_connect.test_cross_connect.cross_connect_group_id}"
		customer_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip}"
		oracle_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip}"
		vlan = "${var.virtual_circuit_cross_connect_mappings_vlan}"
	}
	customer_bgp_asn = "${var.virtual_circuit_customer_bgp_asn}"
	gateway_id = "${oci_core_drg.test_drg.id}"
}
`

	VirtualCircuitProviderResourceConfig = VirtualCircuitResourceDependencies + `
data "oci_core_fast_connect_provider_services" "test_fast_connect_provider_services" {
	#Required
	compartment_id = "${var.compartment_id}"

	filter {
		name = "type"
		values = [ "LAYER2" ]
	}

	filter {
		name = "private_peering_bgp_management"
		values = [ "CUSTOMER_MANAGED" ]
	}

	filter {
		name = "supported_virtual_circuit_types"
		values = [ "${var.virtual_circuit_type}" ]
	}

	filter {
		name = "public_peering_bgp_management"
		values = [ "ORACLE_MANAGED" ]
	}
}

resource "oci_core_virtual_circuit" "test_virtual_circuit" {
	#Required
	compartment_id = "${var.compartment_id}"
	type = "${var.virtual_circuit_type}"

	#Required for PRIVATE VirtualCircuit with Provider
	bandwidth_shape_name = "${var.virtual_circuit_bandwidth_shape_name}"
	cross_connect_mappings {
		customer_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip}"
		oracle_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip}"
	}
	customer_bgp_asn = "${var.virtual_circuit_customer_bgp_asn}"
	display_name = "${var.virtual_circuit_display_name}"
	gateway_id = "${oci_core_drg.test_drg.id}"
	provider_service_id = "${data.oci_core_fast_connect_provider_services.test_fast_connect_provider_services.fast_connect_provider_services.0.id}"
	region = "${var.virtual_circuit_region}"
}
`

	VirtualCircuitResourceConfig = VirtualCircuitResourceDependencies + `
resource "oci_core_virtual_circuit" "test_virtual_circuit" {
	#Required
	compartment_id = "${var.compartment_id}"
	type = "${var.virtual_circuit_type}"

	#Optional
	bandwidth_shape_name = "${var.virtual_circuit_bandwidth_shape_name}"
	cross_connect_mappings {

		#Optional
		#bgp_md5auth_key = "${var.virtual_circuit_cross_connect_mappings_bgp_md5auth_key}"
		cross_connect_or_cross_connect_group_id = "${oci_core_cross_connect.test_cross_connect.cross_connect_group_id}"
		customer_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip}"
		oracle_bgp_peering_ip = "${var.virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip}"
		vlan = "${var.virtual_circuit_cross_connect_mappings_vlan}"
	}
	customer_bgp_asn = "${var.virtual_circuit_customer_bgp_asn}"
	display_name = "${var.virtual_circuit_display_name}"
	gateway_id = "${oci_core_drg.test_drg.id}"
	region = "${var.virtual_circuit_region}"
}
`

	VirtualCircuitPropertyVariables = `
variable "virtual_circuit_bandwidth_shape_name" { default = "10 Gbps" }
variable "virtual_circuit_cross_connect_mappings_bgp_md5auth_key" { default = "bgpMd5AuthKey" }
variable "virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip" { default = "10.0.0.18/31" }
variable "virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip" { default = "10.0.0.19/31" }
variable "virtual_circuit_cross_connect_mappings_vlan" { default = 200 }
variable "virtual_circuit_customer_bgp_asn" { default = 10 }
variable "virtual_circuit_display_name" { default = "displayName" }
variable "virtual_circuit_public_prefixes_cidr_block" { default = "0.0.0.0/5" }
variable "virtual_circuit_region" { default = "r1" }
variable "virtual_circuit_state" { default = "AVAILABLE" }
`

	VirtualCircuitPrivatePropertyVariables = `
variable "virtual_circuit_type" { default = "PRIVATE" }

`

	VirtualCircuitPublicPropertyVariables = `
variable "virtual_circuit_type" { default = "PUBLIC" }

`

	//VirtualCircuitResourceDependencies = GatewayPropertyVariables + GatewayResourceConfig + ProviderServicePropertyVariables + ProviderServiceResourceConfig
	VirtualCircuitResourceDependencies = DrgPropertyVariables + DrgRequiredOnlyResource + CrossConnectPropertyVariables + CrossConnectResourceConfig
)

func TestCoreVirtualCircuitResource_basic(t *testing.T) {
	region := getRequiredEnvSetting("region")
	if !strings.EqualFold("r1", region) {
		t.Skip("FastConnect tests are not yet supported in production regions")
	}

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_virtual_circuit.test_virtual_circuit"
	datasourceName := "data.oci_core_virtual_circuits.test_virtual_circuits"
	singularDatasourceName := "data.oci_core_virtual_circuit.test_virtual_circuit"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create - PUBLIC Virtual Circuit
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + VirtualCircuitPropertyVariables + VirtualCircuitPublicPropertyVariables + compartmentIdVariableStr + VirtualCircuitPublicRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_mappings.0.cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.vlan", "200"),
					resource.TestCheckResourceAttr(resourceName, "customer_bgp_asn", "10"),
					resource.TestCheckResourceAttr(resourceName, "public_prefixes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_prefixes.0.cidr_block", "0.0.0.0/5"),
					resource.TestCheckResourceAttr(resourceName, "type", "PUBLIC"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},
			// delete before next create
			{
				Config: config + compartmentIdVariableStr + VirtualCircuitResourceDependencies,
			},
			// verify create - PRIVATE Virtual Circuit with Provider
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + VirtualCircuitPropertyVariables + VirtualCircuitPrivatePropertyVariables + compartmentIdVariableStr + VirtualCircuitProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.18/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.19/31"),
					resource.TestCheckResourceAttr(resourceName, "customer_bgp_asn", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_service_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_state", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "PRIVATE"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},
			// delete before next create
			{
				Config: config + compartmentIdVariableStr + VirtualCircuitResourceDependencies,
			},

			// verify create - PRIVATE Virtual Circuit
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + VirtualCircuitPropertyVariables + VirtualCircuitPrivatePropertyVariables + compartmentIdVariableStr + VirtualCircuitRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_mappings.0.cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.18/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.19/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.vlan", "200"),
					resource.TestCheckResourceAttr(resourceName, "customer_bgp_asn", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "type", "PRIVATE"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},
			// delete before next create
			{
				Config: config + compartmentIdVariableStr + VirtualCircuitResourceDependencies,
			},

			// verify create with optionals
			{
				Config: config + VirtualCircuitPropertyVariables + VirtualCircuitPrivatePropertyVariables + compartmentIdVariableStr + VirtualCircuitResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bandwidth_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.#", "1"),
					//resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.bgp_md5auth_key", "bgpMd5AuthKey"),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_mappings.0.cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.18/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.19/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.vlan", "200"),
					resource.TestCheckResourceAttr(resourceName, "customer_bgp_asn", "10"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					//resource.TestCheckResourceAttrSet(resourceName, "provider_service_id"),
					//resource.TestCheckResourceAttr(resourceName, "public_prefixes.#", "1"),
					//resource.TestCheckResourceAttr(resourceName, "public_prefixes.0.cidr_block", "0.0.0.0/5"),
					resource.TestCheckResourceAttr(resourceName, "region", "r1"),
					resource.TestCheckResourceAttr(resourceName, "type", "PRIVATE"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + `
variable "virtual_circuit_bandwidth_shape_name" { default = "20 Gbps" }
variable "virtual_circuit_cross_connect_mappings_bgp_md5auth_key" { default = "bgpMd5AuthKey2" }
variable "virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip" { default = "10.0.0.20/31" }
variable "virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip" { default = "10.0.0.21/31" }
variable "virtual_circuit_cross_connect_mappings_vlan" { default = 300 }
variable "virtual_circuit_customer_bgp_asn" { default = 11 }
variable "virtual_circuit_display_name" { default = "displayName2" }
variable "virtual_circuit_public_prefixes_cidr_block" { default = "0.0.0.0/5" }
variable "virtual_circuit_region" { default = "r1" }
variable "virtual_circuit_state" { default = "AVAILABLE" }
variable "virtual_circuit_type" { default = "PRIVATE" }

				                ` + compartmentIdVariableStr + VirtualCircuitResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bandwidth_shape_name", "20 Gbps"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.#", "1"),
					//resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.bgp_md5auth_key", "bgpMd5AuthKey2"),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_mappings.0.cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.20/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.21/31"),
					resource.TestCheckResourceAttr(resourceName, "cross_connect_mappings.0.vlan", "300"),
					resource.TestCheckResourceAttr(resourceName, "customer_bgp_asn", "11"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					//resource.TestCheckResourceAttrSet(resourceName, "provider_service_id"),
					//resource.TestCheckResourceAttr(resourceName, "public_prefixes.#", "1"),
					//resource.TestCheckResourceAttr(resourceName, "public_prefixes.0.cidr_block", "0.0.0.0/5"),
					resource.TestCheckResourceAttr(resourceName, "region", "r1"),
					resource.TestCheckResourceAttr(resourceName, "type", "PRIVATE"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "virtual_circuit_bandwidth_shape_name" { default = "20 Gbps" }
variable "virtual_circuit_cross_connect_mappings_bgp_md5auth_key" { default = "bgpMd5AuthKey2" }
variable "virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip" { default = "10.0.0.20/31" }
variable "virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip" { default = "10.0.0.21/31" }
variable "virtual_circuit_cross_connect_mappings_vlan" { default = 300 }
variable "virtual_circuit_customer_bgp_asn" { default = 11 }
variable "virtual_circuit_display_name" { default = "displayName2" }
variable "virtual_circuit_public_prefixes_cidr_block" { default = "0.0.0.0/5" }
variable "virtual_circuit_region" { default = "r1" }
variable "virtual_circuit_state" { default = "AVAILABLE" }
variable "virtual_circuit_type" { default = "PRIVATE" }

data "oci_core_virtual_circuits" "test_virtual_circuits" {
	#Required
	compartment_id = "${var.compartment_id}"

	#Optional
	display_name = "${var.virtual_circuit_display_name}"
	#state = "${var.virtual_circuit_state}"

	filter {
		name = "id"
		values = ["${oci_core_virtual_circuit.test_virtual_circuit.id}"]
	}
}
				                ` + compartmentIdVariableStr + VirtualCircuitResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					//resource.TestCheckResourceAttrSet(datasourceName, "gateway_id"),
					//resource.TestCheckResourceAttrSet(datasourceName, "provider_service_id"),
					//resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),

					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.bandwidth_shape_name", "20 Gbps"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.cross_connect_mappings.#", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.cross_connect_mappings.0.bgp_md5auth_key", "bgpMd5AuthKey2"),
					resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuits.0.cross_connect_mappings.0.cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.20/31"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.21/31"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.cross_connect_mappings.0.vlan", "300"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.customer_bgp_asn", "11"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuits.0.gateway_id"),
					//resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuits.0.provider_service_id"),
					//resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.public_prefixes.#", "1"),
					//resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.region", "r1"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_circuits.0.type", "PRIVATE"),
				),
			},
			// verify singular datasource
			{
				Config: config + `
variable "virtual_circuit_bandwidth_shape_name" { default = "20 Gbps" }
variable "virtual_circuit_cross_connect_mappings_bgp_md5auth_key" { default = "bgpMd5AuthKey2" }
variable "virtual_circuit_cross_connect_mappings_customer_bgp_peering_ip" { default = "10.0.0.20/31" }
variable "virtual_circuit_cross_connect_mappings_oracle_bgp_peering_ip" { default = "10.0.0.21/31" }
variable "virtual_circuit_cross_connect_mappings_vlan" { default = 300 }
variable "virtual_circuit_customer_bgp_asn" { default = 11 }
variable "virtual_circuit_display_name" { default = "displayName2" }
variable "virtual_circuit_public_prefixes_cidr_block" { default = "0.0.0.0/5" }
variable "virtual_circuit_region" { default = "r1" }
variable "virtual_circuit_state" { default = "AVAILABLE" }
variable "virtual_circuit_type" { default = "PRIVATE" }

data "oci_core_virtual_circuit" "test_virtual_circuit" {
	#Required
	virtual_circuit_id = "${oci_core_virtual_circuit.test_virtual_circuit.id}"
}
				                ` + compartmentIdVariableStr + VirtualCircuitResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "gateway_id"),
					//resource.TestCheckResourceAttrSet(singularDatasourceName, "provider_service_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "virtual_circuit_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "bandwidth_shape_name", "20 Gbps"),
					resource.TestCheckResourceAttr(singularDatasourceName, "bgp_management", "CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(singularDatasourceName, "bgp_session_state", "DOWN"),
					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "cross_connect_mappings.#", "1"),
					//resource.TestCheckResourceAttr(singularDatasourceName, "cross_connect_mappings.0.bgp_md5auth_key", "bgpMd5AuthKey2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "cross_connect_mappings.0.customer_bgp_peering_ip", "10.0.0.20/31"),
					resource.TestCheckResourceAttr(singularDatasourceName, "cross_connect_mappings.0.oracle_bgp_peering_ip", "10.0.0.21/31"),
					resource.TestCheckResourceAttr(singularDatasourceName, "cross_connect_mappings.0.vlan", "300"),
					resource.TestCheckResourceAttr(singularDatasourceName, "customer_bgp_asn", "11"),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "oracle_bgp_asn", "31898"),
					//resource.TestCheckResourceAttr(singularDatasourceName, "provider_state", "providerState2"),
					//resource.TestCheckResourceAttr(singularDatasourceName, "public_prefixes.#", "1"),
					//resource.TestCheckResourceAttr(singularDatasourceName, "reference_comment", "referenceComment2"),
					//resource.TestCheckResourceAttr(singularDatasourceName, "region", "region"),
					resource.TestCheckResourceAttr(singularDatasourceName, "service_type", "COLOCATED"),
					resource.TestCheckResourceAttr(singularDatasourceName, "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttr(singularDatasourceName, "type", "PRIVATE"),
				),
			},
		},
	})
}
