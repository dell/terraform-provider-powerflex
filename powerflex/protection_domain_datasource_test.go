package powerflex

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type pdDataPoints struct {
	name  string
	name2 string
	name3 string
	id    string
}

var protectiondomainTestData pdDataPoints = pdDataPoints{
	id:    "4eeb304600000000",
	name:  "domain1",
	name2: "domain2",
	name3: "domain3",
}

// TestAccProtectionDomainDataSource tests the protectiondomain data source
// where it fetches the protectiondomains based on protectiondomain id/name or storage pool id/name
// and if nothing is mentioned , then return all protectiondomains
func TestAccProtectionDomainDataSource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//retrieving protection domain based on id
			{
				Config: ProtectionDomainDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first protection domain to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_protectiondomain.all", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.0.name", protectiondomainTestData.name),
				),
			},
			//retrieving protection domain based on name
			{
				Config: ProtectionDomainDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first protection domain to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.0.name", protectiondomainTestData.name),
				),
			},
			//retrieving all the protection domains
			{
				Config: ProtectionDomainDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the volume to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.0.name", protectiondomainTestData.name),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.1.name", protectiondomainTestData.name2),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.all", "protection_domains.2.name", protectiondomainTestData.name3),
				),
			},
		},
	})
}

var ProtectionDomainDataSourceConfig1 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}
data "powerflex_protection_domain" "all" {						
	id = "4eeb304600000000"
`

var ProtectionDomainDataSourceConfig2 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}
data "powerflex_protection_domain" "all" {						
	name = "domain1"
}
`

var ProtectionDomainDataSourceConfig3 = `
provider "powerflex" {
	username = "` + username + `"
	password = "` + password + `"
	endpoint = "` + endpoint + `"
	insecure = true
}
data "powerflex_protection_domain" "all" {						
}
`
