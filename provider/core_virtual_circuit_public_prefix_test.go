// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	VirtualCircuitPublicPrefixResourceConfig = VirtualCircuitPublicPrefixResourceDependencies + `

`
	VirtualCircuitPublicPrefixPropertyVariables = `
variable "virtual_circuit_public_prefix_verification_state" { default = "COMPLETED" }

`
	VirtualCircuitPublicPrefixResourceDependencies = VirtualCircuitPropertyVariables + VirtualCircuitPublicPropertyVariables + VirtualCircuitPublicRequiredOnlyResource
)

func TestCoreVirtualCircuitPublicPrefixResource_basic(t *testing.T) {
	region := getRequiredEnvSetting("region")
	if !strings.EqualFold("r1", region) {
		t.Skip("FastConnect tests are not yet supported in production regions")
	}

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_virtual_circuit_public_prefixes.test_virtual_circuit_public_prefixes"

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify datasource
			{
				Config: config + `
variable "virtual_circuit_public_prefix_verification_state" { default = "COMPLETED" }

data "oci_core_virtual_circuit_public_prefixes" "test_virtual_circuit_public_prefixes" {
	#Required
	virtual_circuit_id = "${oci_core_virtual_circuit.test_virtual_circuit.id}"

	#Optional
	verification_state = "${var.virtual_circuit_public_prefix_verification_state}"
	verification_state = "${var.virtual_circuit_public_prefix_verification_state}"
}
                ` + compartmentIdVariableStr + VirtualCircuitPublicPrefixResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "verification_state", "COMPLETED"),
					//resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_id"),

					resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_public_prefixes.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_public_prefixes.0.cidr_block"),
					resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_public_prefixes.0.verification_state"),
					//resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_public_prefixes.0.virtual_circuit_id"),
				),
			},
		},
	})
}
