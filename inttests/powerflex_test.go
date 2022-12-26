package powerflextesting

import (
	"terraform-provider-powerflex/powerflex"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var NAME = ""

const ProviderConfigForTesting = `
	provider "powerflex" {
		insecure = ""
		usecerts = ""
		powerflex_version = ""
	}
`

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerflex": providerserver.NewProtocol6WithError(powerflex.New()),
}
