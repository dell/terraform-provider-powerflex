package main

import (
	"context"
	"terraform-provider-powerflex/powerflex"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name powerflex

func main() {
	providerserver.Serve(context.Background(), powerflex.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as dell.com/dev/powerflex. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "dell.com/dev/powerflex",
	})
}
