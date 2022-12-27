package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

var endpoint = os.Getenv("POWERFLEX_ENDPOINT")
var username = os.Getenv("POWERFLEX_USERNAME")
var password = os.Getenv("POWERFLEX_PASSWORD")

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
