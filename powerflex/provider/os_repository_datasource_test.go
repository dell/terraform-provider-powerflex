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
	"fmt"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var OSRepoDataSourceConfig1 = `
data "powerflex_os_repository" "test" {
}
`

var OSRepoDataSourceConfig2 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repository ids, names, source path, etc.
	filter {
		id = ["` + OSRepoID1 + `"]
	}
  }
`

var OSRepoDataSourceConfig3 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repository ids, names, source path, etc.
	filter {
		id = ["` + OSRepoID1 + `"]
		name = ["` + OSRepoName1 + `"]
		source_path = ["` + OSRepoSourcePath + `"]
		repo_type = ["` + OSRepoType + `"]
		state = ["` + OSRepoState + `"]
		created_by = ["` + OSRepoCreatedBy + `"]
	}
  }
`
var OSRepoDataSourceConfig4 = `
data "powerflex_os_repository" "test" {
	# this datasource supports filters like os repsoitory ids, names, source path, etc.
	filter {
		id = ["^tfacc_.*$"]
	}
  }
`

func TestAccDatasourceOSRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.id", OSRepoID1),
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.name", OSRepoName1),
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.source_path", OSRepoSourcePath),
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.repo_type", OSRepoType),
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.state", OSRepoState),
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.created_by", OSRepoCreatedBy),
				),
			},
			{
				Config: ProviderConfigForTesting + OSRepoDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerflex_os_repository.test", "os_repositories.0.id", OSRepoID1),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAllOsRepositories).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + OSRepoDataSourceConfig1,
				ExpectError: regexp.MustCompile(`.*Error in getting OS repository details*.`),
			},
			// Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + OSRepoDataSourceConfig1,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// Filter Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetDataSourceByValue).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + OSRepoDataSourceConfig2,
				ExpectError: regexp.MustCompile(`.*Error in getting OS repository details*.`),
			},
		},
	})
}
