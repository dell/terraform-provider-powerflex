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
		PreCheck:                 func() { testAccPreCheck(t) },
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
resource "powerflex_parse_csv" "test-csv2" {
	mdm_ip = "10.247.66.67"
	mdm_password = "Password123"
	lia_password="Password123"
csv_detail = [
				  {
					ip = "10.247.66.67"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
			   },
				{
					ip = "10.247.100.214"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Secondary"
					is_sdc = "Yes"
			},
		   {
			ip = "10.247.66.194"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	   },
	   {
		ip = "10.247.103.163"
		password = "dangerous"
		operating_system = "linux"
		is_mdm_or_tb = "Standby"
		is_sdc = "Yes"
   },
		]
    }
`
