package main

import (
	"context"
	"terraform-provider-powerflex/powerflex"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), powerflex.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/dell/powerflex",
	})
}
