/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type dataPoints struct {
	storagePoolID string
	volumeType    string
	dataLayout    string
}

var volumeTestData dataPoints

func init() {
	volumeTestData.storagePoolID = "c98e26e500000000"
	volumeTestData.volumeType = "ThinProvisioned"
	volumeTestData.dataLayout = "MediumGranularity"
}

// TestAccVolumeDataSource tests the volume data source
// where it fetches the volumes based on volume id/name or storage pool id/name
// and if nothing is mentioned , then return all volumes
func TestAccDatasourceVolume(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving volume based on id
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first volume to ensure attributes are correctly set
					resource.TestCheckResourceAttrPair("data.powerflex_volume.all", "volumes.0.id", "powerflex_volume.ref-vol", "id"),
					resource.TestCheckResourceAttrPair("data.powerflex_volume.all", "volumes.0.name", "powerflex_volume.ref-vol", "name"),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.storage_pool_id", volumeTestData.storagePoolID),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.volume_type", volumeTestData.volumeType),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.data_layout", volumeTestData.dataLayout),
				),
			},
			//retrieving volume based on name
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first volume to ensure attributes are correctly set
					resource.TestCheckResourceAttrPair("data.powerflex_volume.all", "volumes.0.id", "powerflex_volume.ref-vol", "id"),
					resource.TestCheckResourceAttrPair("data.powerflex_volume.all", "volumes.0.name", "powerflex_volume.ref-vol", "name"),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.storage_pool_id", volumeTestData.storagePoolID),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.volume_type", volumeTestData.volumeType),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.data_layout", volumeTestData.dataLayout),
				),
			},
			//retrieving volume based on storage pool id
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the volume to ensure storage pool id attributes is correctly set
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.storage_pool_id", volumeTestData.storagePoolID),
				),
			},
			//retrieving volume based on storage pool name
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig4,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the volume to ensure storage pool id attributes is correctly set
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.storage_pool_id", volumeTestData.storagePoolID),
				),
			},
			//retrieving all the volumes
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig5,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the volume to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.storage_pool_id", volumeTestData.storagePoolID),
				),
			},
		},
	})
}

var VolumeDataSourceConfig1 = create8gbVol + `
data "powerflex_volume" "all" {						
	id = resource.powerflex_volume.ref-vol.id
}
`

var VolumeDataSourceConfig2 = create8gbVol + `
data "powerflex_volume" "all" {						
	name = resource.powerflex_volume.ref-vol.name
}
`

var VolumeDataSourceConfig3 = create8gbVol + `
data "powerflex_volume" "all" {						
	storage_pool_id = "c98e26e500000000"
}
`

var VolumeDataSourceConfig4 = create8gbVol + `
data "powerflex_volume" "all" {						
	storage_pool_name = "pool1"
}
`

var VolumeDataSourceConfig5 = create8gbVol + `
data "powerflex_volume" "all" {						
}
`
