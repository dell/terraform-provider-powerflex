/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	//"regexp"
	"encoding/json"
	"os"

	//"log"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var FRDataSourceConfig1 = `
data "powerflex_firmware_repository" "test" {
	firmware_repository_ids = ["8aaa3fda8f5c2609018f854266e12865", "8aaa3fda8f5c2609018f857b6c0d2ede"]
	}
`

func TestAccFRDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfigForTesting + FRDataSourceConfig1,
				//Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDummyFRDataSource(t *testing.T) {
	data, err := os.ReadFile("/root/newworkspace2/terraform-provider-powerflex/powerflex/provider/dummy.json")
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the JSON data into a struct
	var fr scaleiotypes.FirmwareRepositoryDetails
	err = json.Unmarshal(data, &fr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(helper.GetAllFirmwareRepositoryState(&fr))
}
