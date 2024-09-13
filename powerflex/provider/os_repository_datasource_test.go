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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var OSRepoDataSourceConfig1 = `
data "powerflex_os_repository" "test" {
}
`

var OSRepoDataSourceConfig2 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repsoitory ids, names
	filter {
		os_repo_ids = ["` + OSRepoID1 + `"]
	}
  }
`

var OSRepoDataSourceConfig3 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repsoitory ids, names
	filter {
		os_repo_names = ["` + OSRepoName1 + `"]
	}
  }
`

var OSRepoDataSourceConfig4 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repsoitory ids, names
	filter {
		os_repo_ids = ["invalid_id"]
	}
  }
`
var OSRepoDataSourceConfig5 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repsoitory ids, names
	filter {
		os_repo_names = ["invalid_name"]
	}
  }
`

func TestAccDatasourceOSRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + OSRepoDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerflex_os_repository.test", "os_repositories.#"),
				),
			},
			{
				Config: ProviderConfigForTesting + OSRepoDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.id", OSRepoID1),
				),
			},
			{
				Config: ProviderConfigForTesting + OSRepoDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.name", OSRepoName1),
				),
			},
		},
	})
}

func TestAccDatasourceOSRepoNegative(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForTesting + OSRepoDataSourceConfig4,
				ExpectError: regexp.MustCompile(`.*Error in getting OS repository details using id*`),
			},
			{
				Config:      ProviderConfigForTesting + OSRepoDataSourceConfig5,
				ExpectError: regexp.MustCompile(`.*Error in getting OS repository details by names*`),
			},
		},
	})
}
