package powerflex

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

var ProviderConfigForTesting = ``

type sdsDataPoints struct {
	SdsIp1   string
	SdsIp2   string
	SdsIp3   string
	SdsIp4   string
	SdsIp5   string
	SdsIp6   string
	SdsIp7   string
	SdsIp8   string
	SdsIp9   string
	SdsIp10  string
	SdsIp11  string
	SdcIp    string
	volName  string
	volName2 string
	volName3 string
}

var SdsResourceTestData sdsDataPoints

func init() {
	godotenv.Load("POWERFLEX_TERRAFORM_TEST.env")

	username := os.Getenv("POWERFLEX_USERNAME")
	password := os.Getenv("POWERFLEX_PASSWORD")
	endpoint := os.Getenv("POWERFLEX_ENDPOINT")
	SdsResourceTestData.SdsIp1 = os.Getenv("POWERFLEX_SDS_IP_1")
	SdsResourceTestData.SdsIp2 = os.Getenv("POWERFLEX_SDS_IP_2")
	SdsResourceTestData.SdsIp3 = os.Getenv("POWERFLEX_SDS_IP_3")
	SdsResourceTestData.SdsIp4 = os.Getenv("POWERFLEX_SDS_IP_4")
	SdsResourceTestData.SdsIp5 = os.Getenv("POWERFLEX_SDS_IP_5")
	SdsResourceTestData.SdsIp6 = os.Getenv("POWERFLEX_SDS_IP_6")
	SdsResourceTestData.SdsIp7 = os.Getenv("POWERFLEX_SDS_IP_7")
	SdsResourceTestData.SdsIp8 = os.Getenv("POWERFLEX_SDS_IP_8")
	SdsResourceTestData.SdsIp9 = os.Getenv("POWERFLEX_SDS_IP_9")
	SdsResourceTestData.SdsIp10 = os.Getenv("POWERFLEX_SDS_IP_10")
	SdsResourceTestData.SdsIp11 = os.Getenv("POWERFLEX_SDS_IP_11")
	SdsResourceTestData.SdcIp = os.Getenv("POWERFLEX_SDC_IP")
	SdsResourceTestData.volName = os.Getenv("POWERFLEX_VOLUME_NAME")
	SdsResourceTestData.volName2 = os.Getenv("POWERFLEX_VOLUME_NAME_2")
	SdsResourceTestData.volName3 = os.Getenv("POWERFLEX_VOLUME_NAME_3")

	ProviderConfigForTesting = fmt.Sprintf(`
		provider "powerflex" {
			username = "%s"
			password = "%s"
			endpoint = "%s"
		}
	`, username, password, endpoint)
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerflex": providerserver.NewProtocol6WithError(New()),
	}
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("POWERFLEX_USERNAME"); v == "" {
		t.Fatal("POWERFLEX_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("POWERFLEX_PASSWORD"); v == "" {
		t.Fatal("POWERFLEX_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("POWERFLEX_ENDPOINT"); v == "" {
		t.Fatal("POWERFLEX_ENDPOINT must be set for acceptance tests")
	}
}
