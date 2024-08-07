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
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var osResp *scaleiotypes.OSRepository = &scaleiotypes.OSRepository{
	Name: "TestTFOS",
}
var localMocker1 *Mocker
var localMocker2 *Mocker

func TestAccOsRepositoryResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error during create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMocker1 != nil {
						localMocker1.UnPatch()
					}
					if localMocker2 != nil {
						localMocker2.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateOSRepository).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + osRepoResource,
				ExpectError: regexp.MustCompile(`.*Could not create OS Repository*.`),
			},
			// Error getting OS Repository ID
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMocker1 != nil {
						localMocker1.UnPatch()
					}
					if localMocker2 != nil {
						localMocker2.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateOSRepository).Return(osResp, nil).Build()
					localMocker2 = Mock(helper.GetOsRepositoryID).Return("", fmt.Errorf("Get error")).Build()
				},
				Config:      ProviderConfigForTesting + osRepoResource,
				ExpectError: regexp.MustCompile(`.*Could not get the OS Repository id*.`),
			},
			// Error getting OS Repository Details
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMocker1 != nil {
						localMocker1.UnPatch()
					}
					if localMocker2 != nil {
						localMocker2.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateOSRepository).Return(osResp, nil).Build()
					localMocker2 = Mock(helper.GetOsRepositoryID).Return("", nil).Build()
					localMocker1 = Mock(helper.GetOSRepositoryByID).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + osRepoResource,
				ExpectError: regexp.MustCompile(`.*Could not get the OS Repository Details*.`),
			},
			// Successful Create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMocker1 != nil {
						localMocker1.UnPatch()
					}
					if localMocker2 != nil {
						localMocker2.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + osRepoResource,
			},
			// Read error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					if localMocker1 != nil {
						localMocker1.UnPatch()
					}
					if localMocker2 != nil {
						localMocker2.UnPatch()
					}
					localMocker1 = Mock(helper.GetOSRepositoryByID).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + osRepoResource,
				ExpectError: regexp.MustCompile(`.*Could not get the OS Repository Details*.`),
			},
		}})
}

var osRepoResource = `
resource "powerflex_os_repository" "test" {
	name = "TestTFOS"
	repo_type = "ISO"
	source_path = "` + OSRepoSourcePath + `"
	image_type = "vmware_esxi"
 } 
`
