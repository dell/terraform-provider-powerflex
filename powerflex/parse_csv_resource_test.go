package powerflex

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVolumeDataSource tests the volume data source
// where it fetches the volumes based on volume id/name or storage pool id/name
// and if nothing is mentioned , then return all volumes
func TestAccParseCSVResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving volume based on id
			{
				Config: ProviderConfigForGatewayTesting + ParseCSVConfig1,
				Check:  resource.ComposeAggregateTestCheckFunc(
				// Verify the first volume to ensure attributes are correctly set
				),
			},
		},
	})
}

var ParseCSVConfig1 = `
resource "parse_csv" "all" {						
	csv_detail = [
				  {
					ip = "ip1"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
			   },
			   {
				ip = "ip2"
				password = "dangerous"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "Yes"
		   },
		   {
			ip = "ip3"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	   },
	   {
		ip = "ip4"
		password = "dangerous"
		operating_system = "linux"
		is_mdm_or_tb = "Standby"
		is_sdc = "Yes"
   },
		]
}
`
