package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"terraform-provider-powerflex/powerflex"
)

var (
	version string = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	err := providerserver.Serve(context.Background(), powerflex.New, providerserver.ServeOpts{
		Address: "dell.com/ses/powerflex",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
