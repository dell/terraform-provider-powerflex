package powerflex

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

type dataPoints struct {
	name          string
	id            string
	storagePoolID string
	volumeType    string
	dataLayout    string
}

var volumeTestData dataPoints

func init() {
	volumeTestData.name = "cicd-dbc5a5909d"
	volumeTestData.id = "457752ff000000c7"
	volumeTestData.storagePoolID = "7630a24600000000"
	volumeTestData.volumeType = "ThinProvisioned"
	volumeTestData.dataLayout = "MediumGranularity"
}

// TestAccVolumeDataSource tests the volume data source
// where it fetches the volumes based on volume id/name or storage pool id/name
// and if nothing is mentioned , then return all volumes
func TestAccVolumeDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving volume based on id
			{
				Config: ProviderConfigForTesting + VolumeDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first volume to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.id", volumeTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.name", volumeTestData.name),
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
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.id", volumeTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_volume.all", "volumes.0.name", volumeTestData.name),
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

var VolumeDataSourceConfig1 = `
data "powerflex_volume" "all" {						
	id = "457752ff000000c7"
}
`

var VolumeDataSourceConfig2 = `
data "powerflex_volume" "all" {						
	name = "cicd-dbc5a5909d"
}
`

var VolumeDataSourceConfig3 = `
data "powerflex_volume" "all" {						
	storage_pool_id = "7630a24600000000"
}
`

var VolumeDataSourceConfig4 = `
data "powerflex_volume" "all" {						
	storage_pool_name = "pool1"
}
`

var VolumeDataSourceConfig5 = `
data "powerflex_volume" "all" {						
}
`
