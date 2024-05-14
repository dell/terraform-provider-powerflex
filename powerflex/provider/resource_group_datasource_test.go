/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var ResourceGroupDataSourceConfig = `
data "powerflex_resource_group" "example" {				
}
`

var ResourceGroupDataSourceConfig2 = `
data "powerflex_resource_group" "example" {
	service_ids = ["` + ServiceDataPoints.ServiceID + `"]					
}
`

var ResourceGroupDataSourceConfig3 = `
data "powerflex_resource_group" "example" {
	service_names = ["` + ServiceDataPoints.ServiceName + `"]						
}
`

var ResourceGroupDataSourceConfig4 = `
data "powerflex_resource_group" "example" {
	service_ids = ["invalid"]					
}
`

func TestAccResourceGroupDataSource(t *testing.T) {
	t.Skip("Skipping this test case")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "service_details.0.id", ServiceDataPoints.ServiceID),
				),
			},
			{
				Config: ProviderConfigForTesting + ResourceGroupDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_resource_group.example", "service_details.0.deployment_name", ServiceDataPoints.ServiceName),
				),
			},
			{
				Config:      ProviderConfigForTesting + ResourceGroupDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting service details using id*.`),
			},
		},
	})
}
