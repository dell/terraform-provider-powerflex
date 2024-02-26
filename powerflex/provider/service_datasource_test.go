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
	"testing"
	"regexp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var ServiceDataSourceConfig = `
data "powerflex_service" "example" {				
}
`

var ServiceDataSourceConfig2 = `
data "powerflex_service" "example" {
	service_ids = ["` + ServiceDataPoints.ServiceID + `"]					
}
`

var ServiceDataSourceConfig3 = `
data "powerflex_service" "example" {
	service_names = ["` + ServiceDataPoints.ServiceName + `"]						
}
`

var ServiceDataSourceConfig4 = `
data "powerflex_service" "example" {
	service_ids = ["invalid"]					
}
`

func TestAccServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + ServiceDataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: ProviderConfigForTesting + ServiceDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_service.example", "service_details.0.id", ServiceDataPoints.ServiceID),
				),
			},
			{
				Config: ProviderConfigForTesting + ServiceDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_service.example", "service_details.0.deployment_name", ServiceDataPoints.ServiceName),
				),
			},
			{
				Config:      ProviderConfigForTesting + ServiceDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting service details using id*.`),
			},
		},
	})
}
