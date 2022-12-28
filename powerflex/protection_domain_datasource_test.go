package powerflex

import (
	"fmt"
	"os"
	"regexp"
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
	name3: "domain_1",
}

var (
	ProtectionDomainDataSourceConfig1 string
	ProtectionDomainDataSourceConfig2 string
	ProtectionDomainDataSourceConfig3 string
	ProtectionDomainDataSourceConfig4 string
)

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
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd1", "protection_domains.0.name", protectiondomainTestData.name),
				),
			},
			//retrieving protection domain based on name
			{
				Config: ProtectionDomainDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first protection domain to ensure attributes are correctly set
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.#", "1"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd2", "protection_domains.0.name", protectiondomainTestData.name),
					// resource.TestCheckOutput("pdResult", "domain1"),
				),
			},
			//retrieving all the protection domains
			{
				Config: ProtectionDomainDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the volume to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.#", "3"),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.0.id", protectiondomainTestData.id),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.0.name", protectiondomainTestData.name),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.1.name", protectiondomainTestData.name2),
					resource.TestCheckResourceAttr("data.powerflex_protection_domain.pd3", "protection_domains.2.name", protectiondomainTestData.name3),
					// resource.TestCheckOutput("pdResult3", "domain1"),
				),
			},
			//retrieving all the protection domains
			{
				Config:      ProtectionDomainDataSourceConfig4,
				ExpectError: regexp.MustCompile(""),
			},
		},
	})
}

func init() {
	username = os.Getenv("POWERFLEX_USERNAME")
	password = os.Getenv("POWERFLEX_PASSWORD")
	endpoint = os.Getenv("POWERFLEX_ENDPOINT")

	ProtectionDomainDataSourceConfig1 = fmt.Sprintf(`
provider "powerflex" {
	username = "%s"
	password = "%s"
	endpoint = "%s"
	insecure = true
}
data "powerflex_protection_domain" "pd1" {						
	id = "4eeb304600000000"
}
output "pdResult1" {
	value = data.powerflex_protection_domain.pd1.protection_domains[0].name
}
`, username, password, endpoint)

	ProtectionDomainDataSourceConfig2 = fmt.Sprintf(`
provider "powerflex" {
	username = "%s"
	password = "%s"
	endpoint = "%s"
	insecure = true
}
data "powerflex_protection_domain" "pd2" {			
	name = "domain1"
}
output "pdResult2" {
	value = data.powerflex_protection_domain.pd2.protection_domains[0].name
}
`, username, password, endpoint)

	ProtectionDomainDataSourceConfig3 = fmt.Sprintf(`
provider "powerflex" {
	username = "%s"
	password = "%s"
	endpoint = "%s"
	insecure = true
}
data "powerflex_protection_domain" "pd3" {						
}
output "pdResult3" {
	value = data.powerflex_protection_domain.pd3.protection_domains
}
`, username, password, endpoint)

	ProtectionDomainDataSourceConfig4 = fmt.Sprintf(`
provider "powerflex" {
	username = "%s"
	password = "%s"
	endpoint = "%s"
	insecure = true
}
data "powerflex_protection_domain" "pd4" {
	id = "blahblahid"
	name = "blahblahname"					
}
output "pdResult3" {
	value = data.powerflex_protection_domain.pd3.protection_domains
}
`, username, password, endpoint)
}
