package main

import (
	"context"
	"log"
	"terraform-provider-powerflex/powerflex/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name powerflex

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/dell/powerflex. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "registry.terraform.io/dell/powerflex",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
