package powerflex

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

// TestAccSnapshotPolicyDataSource tests the snapshot policy data source
// where it fetches the snapshot policies based on snapshot policy id/name
// and if nothing is mentioned , then return all snapshot policies
func TestAccSnapshotPolicyDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving snapshot policy based on id
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first snapshot policy to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp1", "snapshotpolicies.0.id", "15ad99b900000001"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp1", "snapshotpolicies.0.name", "sample_snap_policy_1"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp1", "snapshotpolicies.#", "1"),
				),
			},
			//retrieving snapshot policy based on name
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first snapshot policy to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp2", "snapshotpolicies.0.id", "15ad99b900000001"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp2", "snapshotpolicies.0.name", "sample_snap_policy_1"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp2", "snapshotpolicies.#", "1"),
				),
			},
			//retrieving all snapshot policies
			{
				Config: ProviderConfigForTesting + SnapshotPolicyDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first snapshot policy to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp3", "snapshotpolicies.0.id", "15ad99b800000000"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp3", "snapshotpolicies.#", "4"),
					resource.TestCheckResourceAttr("data.powerflex_snapshotpolicy.sp3", "snapshotpolicies.0.name", "sample_snap_policy"),
				),
			},
			//retrieving snapshot policy with empty snapshot policy id
			{
				Config:      ProviderConfigForTesting + SnapshotPolicyDataSourceConfig4,
				ExpectError: regexp.MustCompile(".*Invalid Attribute Value Length.*"),
			},
			//retrieving snapshot policy with incorrect snapshot policy id
			{
				Config:      ProviderConfigForTesting + SnapshotPolicyDataSourceConfig5,
				ExpectError: regexp.MustCompile(".*Unable to Read Powerflex Snapshot Policy.*"),
			},
		},
	})
}

var SnapshotPolicyDataSourceConfig1 = `
data "powerflex_snapshotpolicy" "sp1" {						
	id = "15ad99b900000001"
}
`

var SnapshotPolicyDataSourceConfig2 = `
data "powerflex_snapshotpolicy" "sp2" {						
	name = "sample_snap_policy_1"
}
`

var SnapshotPolicyDataSourceConfig3 = `
data "powerflex_snapshotpolicy" "sp3" {						
}
`

var SnapshotPolicyDataSourceConfig4 = `
data "powerflex_snapshotpolicy" "sp4" {						
	id = ""
}
`

var SnapshotPolicyDataSourceConfig5 = `
data "powerflex_snapshotpolicy" "sp5" {	
	id = "15ad99b9000"					
}
`
